package service

import (
	"context"

	"github.com/ONSdigital/dp-identity-api/api"
	cognitoClient "github.com/ONSdigital/dp-identity-api/cognito"
	"github.com/ONSdigital/dp-identity-api/config"
	health "github.com/ONSdigital/dp-identity-api/service/healthcheck"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// Service contains all the configs, server and clients to run the dp-identity-api API
type Service struct {
	Config      *config.Config
	Server      HTTPServer
	Router      *mux.Router
	Api         *api.API
	ServiceList *ExternalServiceList
	HealthCheck HealthChecker
}

// Run the service
func Run(ctx context.Context, cfg *config.Config, serviceList *ExternalServiceList, buildTime, gitCommit, version string, svcErrors chan error) (*Service, error) {

	log.Event(ctx, "running service", log.INFO)

	log.Event(ctx, "using service configuration", log.Data{"config": cfg}, log.INFO)
	
	r := mux.NewRouter()

	s := serviceList.GetHTTPServer(cfg.BindAddr, r)

	cognitoclient := serviceList.GetCognitoClient(cfg.AWSRegion)

	a := api.Setup(ctx, r, cognitoclient, cfg)
	
	hc, err := serviceList.GetHealthCheck(cfg, buildTime, gitCommit, version)

	if err != nil {
		log.Event(ctx, "could not instantiate healthcheck", log.FATAL, log.Error(err))
		return nil, err
	}

	if err := registerCheckers(ctx, hc, cognitoclient, &cfg.AWSCognitoUserPoolID); err != nil {
		return nil, errors.Wrap(err, "unable to register checkers")
	}

	r.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)
	hc.Start(ctx)

	// Run the http server in a new go-routine
	go func() {
		if err := s.ListenAndServe(); err != nil {
			svcErrors <- errors.Wrap(err, "failure in http listen and serve")
		}
	}()

	return &Service{
		Config:      cfg,
		Router:      r,
		Api:         a,
		HealthCheck: hc,
		ServiceList: serviceList,
		Server:      s,
	}, nil
}

// Close gracefully shuts the service down in the required order, with timeout
func (svc *Service) Close(ctx context.Context) error {
	timeout := svc.Config.GracefulShutdownTimeout
	log.Event(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout}, log.INFO)
	ctx, cancel := context.WithTimeout(ctx, timeout)

	// track shutown gracefully closes up
	var hasShutdownError bool

	go func() {
		defer cancel()

		// stop healthcheck, as it depends on everything else
		if svc.ServiceList.HealthCheck {
			svc.HealthCheck.Stop()
		}

		// stop any incoming requests before closing any outbound connections
		if err := svc.Server.Shutdown(ctx); err != nil {
			log.Event(ctx, "failed to shutdown http server", log.Error(err), log.ERROR)
			hasShutdownError = true
		}

		// ADD CODE HERE: Close other dependencies, in the expected order
	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	// timeout expired
	if ctx.Err() == context.DeadlineExceeded {
		log.Event(ctx, "shutdown timed out", log.ERROR, log.Error(ctx.Err()))
		return ctx.Err()
	}

	// other error
	if hasShutdownError {
		err := errors.New("failed to shutdown gracefully")
		log.Event(ctx, "failed to shutdown gracefully ", log.ERROR, log.Error(err))
		return err
	}

	log.Event(ctx, "graceful shutdown was successful", log.INFO)
	return nil
}

func registerCheckers(ctx context.Context, hc HealthChecker, client cognitoClient.Client, userPoolID *string) (err error) {
	hasErrors := false

	if err := hc.AddCheck("cognito healthchecker", health.CognitoHealthCheck(client, userPoolID)); err != nil {
		hasErrors = true
		log.Event(ctx, "error adding check for cognito client", log.ERROR, log.Error(err))
	}

	if hasErrors {
		return errors.New("Error(s) registering checkers for healthcheck")
	}

	return nil
}