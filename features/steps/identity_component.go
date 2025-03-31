package steps

import (
	"context"
	"net/http"
	"time"

	"github.com/ONSdigital/dp-authorisation/v2/authorisation"
	"github.com/ONSdigital/dp-authorisation/v2/authorisationtest"
	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-identity-api/v2/cognito"
	cognitoMock "github.com/ONSdigital/dp-identity-api/v2/cognito/mock"
	"github.com/ONSdigital/dp-identity-api/v2/config"
	jwksMock "github.com/ONSdigital/dp-identity-api/v2/jwks/mock"
	"github.com/ONSdigital/dp-identity-api/v2/service"
	"github.com/ONSdigital/dp-identity-api/v2/service/mock"
	"github.com/ONSdigital/dp-permissions-api/sdk"
	"github.com/ONSdigital/log.go/v2/log"
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
	JWKSManager             *jwksMock.ManagerMock
}

func NewIdentityComponent() (*IdentityComponent, error) {
	svcErrors := make(chan error, 1)

	c := &IdentityComponent{
		HTTPServer: &http.Server{
			ReadHeaderTimeout: 5 * time.Second,
		},
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
	c.Config.AWSCognitoClientID = "client-aaa-bbb"
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

	c.JWKSManager = jwksMock.JWKSStubbed
	c.svc, err = service.Run(context.Background(), c.Config, c.svcList, c.JWKSManager, "1", "", "", c.errorChan)
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
	if err := fakePermissionsAPI.UpdatePermissionsBundleResponse(bundle); err != nil {
		log.Error(context.Background(), "failed to update permissions bundle response", err)
	}
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

func (c *IdentityComponent) DoGetHealthcheckOk(_ *config.Config, _, _, _ string) (service.HealthChecker, error) {
	return &mock.HealthCheckerMock{
		AddCheckFunc: func(_ string, _ healthcheck.Checker) error { return nil },
		StartFunc:    func(_ context.Context) {},
		StopFunc:     func() {},
	}, nil
}

func (c *IdentityComponent) DoGetHTTPServer(bindAddr string, router http.Handler, _ *config.Config) service.HTTPServer {
	c.HTTPServer.Addr = bindAddr
	c.HTTPServer.Handler = router
	return c.HTTPServer
}

func (c *IdentityComponent) DoGetCognitoClient(ctx context.Context, _ string) cognito.Client {
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
