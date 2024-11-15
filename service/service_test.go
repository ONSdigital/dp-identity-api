package service_test

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/ONSdigital/dp-authorisation/v2/authorisation"
	authorisationMock "github.com/ONSdigital/dp-authorisation/v2/authorisation/mock"

	"github.com/ONSdigital/dp-identity-api/v2/cognito"
	cognitoMock "github.com/ONSdigital/dp-identity-api/v2/cognito/mock"
	jwksMock "github.com/ONSdigital/dp-identity-api/v2/jwks/mock"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"

	"github.com/ONSdigital/dp-identity-api/v2/config"
	"github.com/ONSdigital/dp-identity-api/v2/service"

	serviceMock "github.com/ONSdigital/dp-identity-api/v2/service/mock"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	ctx           = context.Background()
	testBuildTime = "BuildTime"
	testGitCommit = "GitCommit"
	testVersion   = "Version"
	errServer     = errors.New("HTTP Server error")
)

var (
	errHealthcheck = errors.New("healthCheck error")
)

var funcDoGetHealthcheckErr = func(_ *config.Config, _, _, _ string) (service.HealthChecker, error) {
	return nil, errHealthcheck
}

var funcDoGetHTTPServerNil = func(_ string, _ http.Handler, _ *config.Config) service.HTTPServer {
	return nil
}

var jwksHandler = &jwksMock.ManagerMock{}

func TestRun(t *testing.T) {
	Convey("Having a set of mocked dependencies", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)

		// set dummy config data
		cfg.AWSCognitoUserPoolID = "eu-west-18_73289nds8w932"
		cfg.AWSCognitoClientID = "client-aaa-bbb"
		cfg.AWSCognitoClientSecret = "secret-ccc-ddd"
		cfg.AWSAuthFlow = "authflow"

		hcMock := &serviceMock.HealthCheckerMock{
			AddCheckFunc: func(_ string, _ healthcheck.Checker) error { return nil },
			StartFunc:    func(_ context.Context) {},
		}

		serverWg := &sync.WaitGroup{}
		serverMock := &serviceMock.HTTPServerMock{
			ListenAndServeFunc: func() error {
				serverWg.Done()
				return nil
			},
		}

		failingServerMock := &serviceMock.HTTPServerMock{
			ListenAndServeFunc: func() error {
				serverWg.Done()
				return errServer
			},
		}

		funcDoGetHealthcheckOk := func(_ *config.Config, _, _, _ string) (service.HealthChecker, error) {
			return hcMock, nil
		}

		funcDoGetHTTPServer := func(_ string, _ http.Handler, _ *config.Config) service.HTTPServer {
			return serverMock
		}

		funcDoGetFailingHTTPSerer := func(_ string, _ http.Handler, _ *config.Config) service.HTTPServer {
			return failingServerMock
		}

		Convey("Given that initialising healthcheck returns an error", func() {
			initMock := &serviceMock.InitialiserMock{
				DoGetHTTPServerFunc:              funcDoGetHTTPServerNil,
				DoGetHealthCheckFunc:             funcDoGetHealthcheckErr,
				DoGetCognitoClientFunc:           DoGetCognitoClient,
				DoGetAuthorisationMiddlewareFunc: DoGetAuthorisationMiddleware,
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			_, err := service.Run(ctx, cfg, svcList, jwksHandler, testBuildTime, testGitCommit, testVersion, svcErrors)

			Convey("Then service Run fails with the same error and the flag is not set", func() {
				So(err, ShouldResemble, errHealthcheck)
				So(svcList.HealthCheck, ShouldBeFalse)
			})
		})

		Convey("Given that initialisation of the authorisation middleware fails", func() {
			expectedError := errors.New("failed to init authorisation middleware")
			initMock := &serviceMock.InitialiserMock{
				DoGetHealthCheckFunc:   funcDoGetHealthcheckOk,
				DoGetHTTPServerFunc:    funcDoGetFailingHTTPSerer,
				DoGetCognitoClientFunc: DoGetCognitoClient,
				DoGetAuthorisationMiddlewareFunc: func(_ context.Context, _ *authorisation.Config) (authorisation.Middleware, error) {
					return nil, expectedError
				},
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			_, err := service.Run(ctx, cfg, svcList, jwksHandler, testBuildTime, testGitCommit, testVersion, svcErrors)

			Convey("Then service Run fails with the expected error", func() {
				So(err, ShouldEqual, expectedError)
				So(svcList.AuthMiddleware, ShouldBeFalse)
			})
		})

		Convey("Given that all dependencies are successfully initialised but the http server fails", func() {
			initMock := &serviceMock.InitialiserMock{
				DoGetHealthCheckFunc:             funcDoGetHealthcheckOk,
				DoGetHTTPServerFunc:              funcDoGetFailingHTTPSerer,
				DoGetCognitoClientFunc:           DoGetCognitoClient,
				DoGetAuthorisationMiddlewareFunc: DoGetAuthorisationMiddleware,
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			serverWg.Add(1)
			_, err := service.Run(ctx, cfg, svcList, jwksHandler, testBuildTime, testGitCommit, testVersion, svcErrors)

			Convey("Then dependencies are initialised with all the flags being set", func() {
				So(err, ShouldBeNil)
				So(svcList.HealthCheck, ShouldBeTrue)
				So(svcList.AuthMiddleware, ShouldBeTrue)
			})

			Convey("But an error is returned in the error channel", func() {
				sErr := <-svcErrors
				So(sErr.Error(), ShouldResemble, fmt.Sprintf("failure in http listen and serve: %s", errServer.Error()))
				So(len(failingServerMock.ListenAndServeCalls()), ShouldEqual, 1)
			})
		})

		Convey("Given that all dependencies are successfully initialised", func() {
			initMock := &serviceMock.InitialiserMock{
				DoGetHTTPServerFunc:              funcDoGetHTTPServer,
				DoGetHealthCheckFunc:             funcDoGetHealthcheckOk,
				DoGetCognitoClientFunc:           DoGetCognitoClient,
				DoGetAuthorisationMiddlewareFunc: DoGetAuthorisationMiddleware,
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			serverWg.Add(1)
			_, err := service.Run(ctx, cfg, svcList, jwksHandler, testBuildTime, testGitCommit, testVersion, svcErrors)

			Convey("Then service Run succeeds and all the flags are set", func() {
				So(err, ShouldBeNil)
				So(svcList.HealthCheck, ShouldBeTrue)
				So(svcList.AuthMiddleware, ShouldBeTrue)
			})

			Convey("And the checkers are registered and the healthcheck and http server have started", func() {
				So(len(hcMock.AddCheckCalls()), ShouldEqual, 2)
				So(hcMock.AddCheckCalls()[0].Name, ShouldResemble, "Cognito")
				So(hcMock.AddCheckCalls()[1].Name, ShouldResemble, "Permissions API")

				So(len(initMock.DoGetHTTPServerCalls()), ShouldEqual, 1)
				So(initMock.DoGetHealthCheckCalls(), ShouldHaveLength, 1)
				So(initMock.DoGetHTTPServerCalls()[0].BindAddr, ShouldEqual, "localhost:25600")
				So(initMock.DoGetCognitoClientCalls(), ShouldHaveLength, 1)
				So(initMock.DoGetAuthorisationMiddlewareCalls(), ShouldHaveLength, 1)
				So(len(hcMock.StartCalls()), ShouldEqual, 1)
				//!!! a call needed to stop the server, maybe ?
				serverWg.Wait() // Wait for HTTP server go-routine to finish
				So(len(serverMock.ListenAndServeCalls()), ShouldEqual, 1)
			})
		})
	})
}

func TestClose(t *testing.T) {
	Convey("Having a correctly initialised service", t, func() {
		cfg, err := config.Get()

		// set dummy config data
		cfg.AWSCognitoUserPoolID = "eu-west-18_73289nds8w932"
		cfg.AWSCognitoClientID = "client-aaa-bbb"
		cfg.AWSCognitoClientSecret = "secret-ccc-ddd"
		cfg.AWSAuthFlow = "authflow"

		So(err, ShouldBeNil)

		hcStopped := false

		// healthcheck Stop does not depend on any other service being closed/stopped
		hcMock := &serviceMock.HealthCheckerMock{
			AddCheckFunc: func(_ string, _ healthcheck.Checker) error { return nil },
			StartFunc:    func(_ context.Context) {},
			StopFunc:     func() { hcStopped = true },
		}

		// server Shutdown will fail if healthcheck is not stopped
		serverMock := &serviceMock.HTTPServerMock{
			ListenAndServeFunc: func() error { return nil },
			ShutdownFunc: func(_ context.Context) error {
				if !hcStopped {
					return errors.New("Server stopped before healthcheck")
				}
				return nil
			},
		}
		authorisationMiddleware := &authorisationMock.MiddlewareMock{
			RequireFunc: func(_ string, handlerFunc http.HandlerFunc) http.HandlerFunc {
				return handlerFunc
			},
			CloseFunc: func(_ context.Context) error {
				return nil
			},
		}

		Convey("Closing the service results in all the dependencies being closed in the expected order", func() {
			initMock := &serviceMock.InitialiserMock{
				DoGetHTTPServerFunc: func(_ string, _ http.Handler, _ *config.Config) service.HTTPServer { return serverMock },
				DoGetHealthCheckFunc: func(_ *config.Config, _, _, _ string) (service.HealthChecker, error) {
					return hcMock, nil
				},
				DoGetCognitoClientFunc: DoGetCognitoClient,
				DoGetAuthorisationMiddlewareFunc: func(_ context.Context, _ *authorisation.Config) (authorisation.Middleware, error) {
					return authorisationMiddleware, nil
				},
			}

			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			svc, err := service.Run(ctx, cfg, svcList, jwksHandler, testBuildTime, testGitCommit, testVersion, svcErrors)
			So(err, ShouldBeNil)

			err = svc.Close(context.Background())
			So(err, ShouldBeNil)
			So(len(hcMock.StopCalls()), ShouldEqual, 1)
			So(len(serverMock.ShutdownCalls()), ShouldEqual, 1)
			So(len(authorisationMiddleware.CloseCalls()), ShouldEqual, 1)
		})

		Convey("If services fail to stop, the Close operation tries to close all dependencies and returns an error", func() {
			failingserverMock := &serviceMock.HTTPServerMock{
				ListenAndServeFunc: func() error { return nil },
				ShutdownFunc: func(_ context.Context) error {
					return errors.New("Failed to stop http server")
				},
			}
			authorisationMiddleware := &authorisationMock.MiddlewareMock{
				RequireFunc: func(_ string, handlerFunc http.HandlerFunc) http.HandlerFunc {
					return handlerFunc
				},
				CloseFunc: func(_ context.Context) error {
					return errors.New("failed to close authorisation middleware")
				},
			}

			initMock := &serviceMock.InitialiserMock{
				DoGetHTTPServerFunc: func(_ string, _ http.Handler, _ *config.Config) service.HTTPServer {
					return failingserverMock
				},
				DoGetHealthCheckFunc: func(_ *config.Config, _, _, _ string) (service.HealthChecker, error) {
					return hcMock, nil
				},
				DoGetCognitoClientFunc: DoGetCognitoClient,
				DoGetAuthorisationMiddlewareFunc: func(_ context.Context, _ *authorisation.Config) (authorisation.Middleware, error) {
					return authorisationMiddleware, nil
				},
			}

			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			svc, err := service.Run(ctx, cfg, svcList, jwksHandler, testBuildTime, testGitCommit, testVersion, svcErrors)
			So(err, ShouldBeNil)

			err = svc.Close(context.Background())
			So(err, ShouldNotBeNil)
			So(len(hcMock.StopCalls()), ShouldEqual, 1)
			So(len(failingserverMock.ShutdownCalls()), ShouldEqual, 1)
			So(len(authorisationMiddleware.CloseCalls()), ShouldEqual, 1)
		})

		Convey("If service times out while shutting down, the Close operation fails with the expected error", func() {
			cfg.GracefulShutdownTimeout = 1 * time.Millisecond
			timeoutServerMock := &serviceMock.HTTPServerMock{
				ListenAndServeFunc: func() error { return nil },
				ShutdownFunc: func(_ context.Context) error {
					time.Sleep(10 * time.Millisecond)
					return nil
				},
			}

			svcList := service.NewServiceList(nil)
			svcList.HealthCheck = true
			svc := service.Service{
				Config:      cfg,
				ServiceList: svcList,
				Server:      timeoutServerMock,
				HealthCheck: hcMock,
			}

			err = svc.Close(context.Background())
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldResemble, "context deadline exceeded")
			So(len(hcMock.StopCalls()), ShouldEqual, 1)
			So(len(timeoutServerMock.ShutdownCalls()), ShouldEqual, 1)
		})
	})
}

func DoGetCognitoClient(_ string) cognito.Client {
	return &cognitoMock.CognitoIdentityProviderClientStub{}
}

func DoGetAuthorisationMiddleware(_ context.Context, _ *authorisation.Config) (authorisation.Middleware, error) {
	return &authorisationMock.MiddlewareMock{
		RequireFunc: func(_ string, handlerFunc http.HandlerFunc) http.HandlerFunc {
			return handlerFunc
		},
	}, nil
}
