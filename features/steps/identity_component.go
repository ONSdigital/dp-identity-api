package steps

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-authorisation/v2/authorisation"
	"github.com/ONSdigital/dp-authorisation/v2/authorisationtest"
	"github.com/ONSdigital/dp-permissions-api/sdk"

	"github.com/ONSdigital/dp-identity-api/cognito"
	cognitoMock "github.com/ONSdigital/dp-identity-api/cognito/mock"
	jwksMock "github.com/ONSdigital/dp-identity-api/jwks/mock"

	"github.com/ONSdigital/dp-identity-api/config"
	"github.com/ONSdigital/dp-identity-api/service"
	"github.com/ONSdigital/dp-identity-api/service/mock"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

type IdentityComponent struct {
	ErrorFeature            componenttest.ErrorFeature
	svcList                 *service.ExternalServiceList
	svc                     *service.Service
	errorChan               chan error
	Config                  *config.Config
	HTTPServer              *http.Server
	ServiceRunning          bool
	apiFeature              *componenttest.APIFeature
	CognitoClient           *cognitoMock.CognitoIdentityProviderClientStub
	AuthorisationMiddleware authorisation.Middleware
	JWKSHandler             *jwksMock.JWKSIntMock
}

func NewIdentityComponent() (*IdentityComponent, error) {
	svcErrors := make(chan error, 1)

	c := &IdentityComponent{
		HTTPServer:     &http.Server{},
		errorChan:      svcErrors,
		ServiceRunning: false,
	}

	var err error

	c.Config, err = config.Get()
	if err != nil {
		return nil, err
	}

	// set dummy user pool id
	c.Config.AWSCognitoUserPoolID = "eu-west-18_73289nds8w932"
	c.Config.AWSCognitoClientId = "client-aaa-bbb"
	c.Config.AWSCognitoClientSecret = "secret-ccc-ddd"
	c.Config.AWSAuthFlow = "USER_PASSWORD_AUTH"

	fakePermissionsAPI := setupFakePermissionsAPI()
	c.Config.AuthorisationConfig.PermissionsAPIURL = fakePermissionsAPI.URL()

	initMock := &mock.InitialiserMock{
		DoGetHealthCheckFunc:             c.DoGetHealthcheckOk,
		DoGetHTTPServerFunc:              c.DoGetHTTPServer,
		DoGetCognitoClientFunc:           c.DoGetCognitoClient,
		DoGetAuthorisationMiddlewareFunc: c.DoGetAuthorisationMiddleware,
	}

	c.svcList = service.NewServiceList(initMock)

	c.JWKSHandler = jwksMock.JWKSStubbed
	c.svc, err = service.Run(context.Background(), c.Config, c.svcList, c.JWKSHandler, "1", "", "", c.errorChan)
	if err != nil {
		return nil, err
	}

	c.ServiceRunning = true
	c.apiFeature = componenttest.NewAPIFeature(c.InitialiseService)

	return c, nil
}

func setupFakePermissionsAPI() *authorisationtest.FakePermissionsAPI {
	fakePermissionsAPI := authorisationtest.NewFakePermissionsAPI()
	bundle := getPermissionsBundle()
	fakePermissionsAPI.Reset()
	fakePermissionsAPI.UpdatePermissionsBundleResponse(bundle)
	return fakePermissionsAPI
}

func getPermissionsBundle() *sdk.Bundle {
	return &sdk.Bundle{
		"users:create": { // role
			"groups/role-admin": { // group
				{
					ID: "1", // policy
				},
			},
		},
		"users:read": { // role
			"groups/role-admin": { // group
				{
					ID: "2", // policy
				},
			},
		},
		"users:update": { // role
			"groups/role-admin": { // group
				{
					ID: "2", // policy
				},
			},
		},
		"groups:create": { // role
			"groups/role-admin": { // group
				{
					ID: "1", // policy
				},
			},
		},
		"groups:read": { // role
			"groups/role-admin": { // group
				{
					ID: "2", // policy
				},
			},
		},
		"groups:update": { // role
			"groups/role-admin": { // group
				{
					ID: "2", // policy
				},
			},
		},
		"groups:delete": { // role
			"groups/role-admin": { // group
				{
					ID: "2", // policy
				},
			},
		},
	}
}

func (c *IdentityComponent) Reset() *IdentityComponent {
	c.apiFeature.Reset()
	return c
}

func (c *IdentityComponent) Close() error {
	if c.svc != nil && c.ServiceRunning {
		c.svc.Close(context.Background())
		c.ServiceRunning = false
	}
	return nil
}

func (c *IdentityComponent) InitialiseService() (http.Handler, error) {
	return c.HTTPServer.Handler, nil
}

func (c *IdentityComponent) DoGetHealthcheckOk(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
	return &mock.HealthCheckerMock{
		AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
		StartFunc:    func(ctx context.Context) {},
		StopFunc:     func() {},
	}, nil
}

func (c *IdentityComponent) DoGetHTTPServer(bindAddr string, router http.Handler, cfg *config.Config) service.HTTPServer {
	c.HTTPServer.Addr = bindAddr
	c.HTTPServer.Handler = router
	return c.HTTPServer
}

func (c *IdentityComponent) DoGetCognitoClient(AWSRegion string) cognito.Client {
	c.CognitoClient = &cognitoMock.CognitoIdentityProviderClientStub{}
	return c.CognitoClient
}

func (c *IdentityComponent) DoGetAuthorisationMiddleware(ctx context.Context, cfg *authorisation.Config) (authorisation.Middleware, error) {
	middleware, err := authorisation.NewMiddlewareFromConfig(ctx, cfg, cfg.JWTVerificationPublicKeys)
	if err != nil {
		return nil, err
	}

	c.AuthorisationMiddleware = middleware
	return c.AuthorisationMiddleware, nil
}
