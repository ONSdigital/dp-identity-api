package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ONSdigital/dp-identity-api/v2/config"
	"github.com/ONSdigital/dp-identity-api/v2/jwks"
	"github.com/ONSdigital/dp-identity-api/v2/service"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/pkg/errors"
)

const serviceName = "dp-identity-api"

var (
	// BuildTime represents the time in which the service was built
	BuildTime string
	// GitCommit represents the commit (SHA-1) hash of the service that is running
	GitCommit string
	// Version represents the version of the service that is running
	Version string

	/* NOTE: replace the above with the below to run code with for example vscode debugger.
	   BuildTime string = "1601119818"
	   GitCommit string = "6584b786caac36b6214ffe04bf62f058d4021538"
	   Version   string = "v0.1.0"
	*/
)

func main() {
	log.Namespace = serviceName
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatal(ctx, "fatal runtime error", err)
	}
}

func run(ctx context.Context) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	// Run the service, providing an error channel for fatal errors
	svcErrors := make(chan error, 1)
	svcList := service.NewServiceList(&service.Init{})

	log.Info(ctx, "dp-identity-api version", log.Data{"version": Version})

	// jks handler
	jwksHandler := &jwks.JWKS{}

	// Read config
	cfg, err := config.Get()
	if err != nil {
		return errors.Wrap(err, "error getting configuration")
	}

	// Retrieve the JWKS RSA Public Keys from Cognito on startup
	jwksRSAKeys, err := jwksHandler.GetJWKSRSAKeys(cfg.AWSRegion, cfg.AWSCognitoUserPoolID)
	if err != nil {
		log.Fatal(ctx, "could not retrieve the JWKS RSA public keys", err)
		return err
	}

	// Set the JWKS RSA public keys authorisation config
	cfg.AuthorisationConfig.JWTVerificationPublicKeys = jwksRSAKeys

	// sensitive fields are omitted from config.String().
	log.Info(ctx, "loaded config", log.Data{
		"config": cfg,
	})

	// Start service
	svc, err := service.Run(ctx, cfg, svcList, jwksHandler.DoGetJWKS(ctx), BuildTime, GitCommit, Version, svcErrors)
	if err != nil {
		return errors.Wrap(err, "running service failed")
	}

	// blocks until an os interrupt or a fatal error occurs
	select {
	case err := <-svcErrors:
		// ADD CODE HERE : call svc.Close(ctx) (or something specific)
		//  if there are any service connections like Kafka that you need to shut down
		return errors.Wrap(err, "service error received")
	case sig := <-signals:
		log.Info(ctx, "os signal received", log.Data{"signal": sig})
	}
	return svc.Close(ctx)
}
