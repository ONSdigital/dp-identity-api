package service

import (
	"context"

	"github.com/ONSdigital/dp-authorisation/v2/authorisation"
	"github.com/ONSdigital/dp-identity-api/api"
	cognitoClient "github.com/ONSdigital/dp-identity-api/cognito"
	"github.com/ONSdigital/dp-identity-api/config"
	"github.com/ONSdigital/dp-identity-api/jwks"
	health "github.com/ONSdigital/dp-identity-api/service/healthcheck"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// Service contains all the configs, server and clients to run the dp-identity-api API
type Service struct {
	Config                  *config.Config
	Server                  HTTPServer
	Router                  *mux.Router
	API                     *api.API
	ServiceList             *ExternalServiceList
	HealthCheck             HealthChecker
	authorisationMiddleware authorisation.Middleware
}

// Run the service
func Run(ctx context.Context, cfg *config.Config, serviceList *ExternalServiceList, jwksHandler jwks.JWKSInt, buildTime, gitCommit, version string, svcErrors chan error) (*Service, error) {
	log.Info(ctx, "running service")
	log.Info(ctx, "using service configuration", log.Data{"config": cfg})

	r := mux.NewRouter()

	s := serviceList.GetHTTPServer(cfg.BindAddr, r, cfg)

	cognitoclient := serviceList.GetCognitoClient(cfg.AWSRegion)

	authorisationMiddleware, err := serviceList.GetAuthorisationMiddleware(ctx, cfg.AuthorisationConfig)
	if err != nil {
		log.Fatal(ctx, "could not instantiate authorisation middleware", err)
		return nil, err
	}

	a, err := api.Setup(ctx, r, cognitoclient, cfg.AWSCognitoUserPoolID, cfg.AWSCognitoClientID, cfg.AWSCognitoClientSecret, cfg.AWSRegion, cfg.AWSAuthFlow, cfg.AllowedEmailDomains, authorisationMiddleware, jwksHandler)
	if err != nil {
		log.Fatal(ctx, "error returned from api setup", err)
		return nil, err
	}

	hc, err := serviceList.GetHealthCheck(cfg, buildTime, gitCommit, version)
	if err != nil {
		log.Fatal(ctx, "could not instantiate healthcheck", err)
		return nil, err
	}

	if err := registerCheckers(ctx, hc, cognitoclient, &cfg.AWSCognitoUserPoolID, authorisationMiddleware); err != nil {
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
		Config:                  cfg,
		Router:                  r,
		API:                     a,
		HealthCheck:             hc,
		ServiceList:             serviceList,
		Server:                  s,
		authorisationMiddleware: authorisationMiddleware,
	}, nil
}

// Close gracefully shuts the service down in the required order, with timeout
func (svc *Service) Close(ctx context.Context) error {
	timeout := svc.Config.GracefulShutdownTimeout
	log.Info(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout})
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
			log.Error(ctx, "failed to shutdown http server", err)
			hasShutdownError = true
		}

		if svc.ServiceList.AuthMiddleware {
			if err := svc.authorisationMiddleware.Close(ctx); err != nil {
				log.Error(ctx, "failed to close authorisation middleware", err)
				hasShutdownError = true
			}
		}
	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	// timeout expired
	if ctx.Err() == context.DeadlineExceeded {
		log.Error(ctx, "shutdown timed out", ctx.Err())
		return ctx.Err()
	}

	// other error
	if hasShutdownError {
		err := errors.New("failed to shutdown gracefully")
		log.Error(ctx, "failed to shutdown gracefully ", err)
		return err
	}

	log.Info(ctx, "graceful shutdown was successful")
	return nil
}

func registerCheckers(ctx context.Context, hc HealthChecker, client cognitoClient.Client, userPoolID *string, authorisationMiddleware authorisation.Middleware) (err error) {
	hasErrors := false

	if err := hc.AddCheck("Cognito", health.CognitoHealthCheck(ctx, client, userPoolID)); err != nil {
		hasErrors = true
		log.Error(ctx, "error adding health checker for Cognito", err)
	}

	if err := hc.AddCheck("Permissions API", authorisationMiddleware.HealthCheck); err != nil {
		hasErrors = true
		log.Error(ctx, "error adding health checker for Permissions API", err)
	}

	if hasErrors {
		return errors.New("Error(s) registering checkers for healthcheck")
	}

	return nil
}
