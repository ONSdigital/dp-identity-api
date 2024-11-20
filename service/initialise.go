package service

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-authorisation/v2/authorisation"
	cognitoclient "github.com/ONSdigital/dp-identity-api/v2/cognito"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/ONSdigital/dp-identity-api/v2/config"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
)

// ExternalServiceList holds the initialiser and initialisation state of external services.
type ExternalServiceList struct {
	AuthMiddleware bool
	HealthCheck    bool
	Init           Initialiser
}

// NewServiceList creates a new service list with the provided initialiser
func NewServiceList(initialiser Initialiser) *ExternalServiceList {
	return &ExternalServiceList{
		AuthMiddleware: false,
		HealthCheck:    false,
		Init:           initialiser,
	}
}

// Init implements the Initialiser interface to initialise dependencies
type Init struct{}

// GetHTTPServer creates an http server
func (e *ExternalServiceList) GetHTTPServer(bindAddr string, router http.Handler, cfg *config.Config) HTTPServer {
	s := e.Init.DoGetHTTPServer(bindAddr, router, cfg)
	return s
}

// GetHealthCheck creates a healthcheck with versionInfo and sets the HealthCheck flag to true
func (e *ExternalServiceList) GetHealthCheck(cfg *config.Config, buildTime, gitCommit, version string) (HealthChecker, error) {
	hc, err := e.Init.DoGetHealthCheck(cfg, buildTime, gitCommit, version)
	if err != nil {
		return nil, err
	}
	e.HealthCheck = true
	return hc, nil
}

// GetCognitoClient creates a cognito client
func (e *ExternalServiceList) GetCognitoClient(region string) cognitoclient.Client {
	client := e.Init.DoGetCognitoClient(region)
	return client
}

// GetAuthorisationMiddleware creates a new instance of authorisation.Middlware
func (e *ExternalServiceList) GetAuthorisationMiddleware(ctx context.Context, authorisationConfig *authorisation.Config) (authorisation.Middleware, error) {
	am, err := e.Init.DoGetAuthorisationMiddleware(ctx, authorisationConfig)
	if err != nil {
		return nil, err
	}
	e.AuthMiddleware = true
	return am, nil
}

// DoGetHTTPServer creates an HTTP Server with the provided bind address and router
func (e *Init) DoGetHTTPServer(bindAddr string, router http.Handler, cfg *config.Config) HTTPServer {
	s := dphttp.NewServer(bindAddr, router)
	s.HandleOSSignals = false
	if cfg.HTTPWriteTimeout != nil {
		s.WriteTimeout = *cfg.HTTPWriteTimeout
	}
	return s
}

// DoGetHealthCheck creates a healthcheck with versionInfo
func (e *Init) DoGetHealthCheck(cfg *config.Config, buildTime, gitCommit, version string) (HealthChecker, error) {
	versionInfo, err := healthcheck.NewVersionInfo(buildTime, gitCommit, version)
	if err != nil {
		return nil, err
	}
	hc := healthcheck.New(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)
	return &hc, nil
}

// DoGetCognitoClient creates a CognitoClient with the provided region
func (e *Init) DoGetCognitoClient(awsRegion string) cognitoclient.Client {
	client := cognito.New(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})), &aws.Config{Region: &awsRegion})
	return client
}

// DoGetAuthorisationMiddleware creates authorisation middleware for the given config
func (e *Init) DoGetAuthorisationMiddleware(ctx context.Context, authorisationConfig *authorisation.Config) (authorisation.Middleware, error) {
	return authorisation.NewFeatureFlaggedMiddleware(ctx, authorisationConfig, authorisationConfig.JWTVerificationPublicKeys)
}
