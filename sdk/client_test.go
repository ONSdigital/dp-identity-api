package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	healthcheck "github.com/ONSdigital/dp-api-clients-go/v2/health"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	dphttp "github.com/ONSdigital/dp-net/v2/http"

	"github.com/ONSdigital/dp-identity-api/models"
	. "github.com/smartystreets/goconvey/convey"
)

const testHost = "http://localhost:23900"

var (
	initialTestState = healthcheck.CreateCheckState(service)

	defaultCredentials = models.UserSignIn{
		Email:    "florence@magicroundabout.org",
		Password: "themagicroundaboutlives",
	}

	defaultExpirationTime        = "2023-09-27 17:30:00.000000000 +0000 UTC"
	defaultRefreshExpirationTime = "2023-09-28 16:30:00.000000000 +0000 UTC"
	defaultId                    = "testId1234"
	defaultRefreshToken          = "refreshtoken1234"
	defaultAuthorization         = "Bearer testToken1234"
)

func newMockHTTPClient(r *http.Response, err error) *dphttp.ClienterMock {
	return &dphttp.ClienterMock{
		SetPathsWithNoRetriesFunc: func(paths []string) {
			// This gets called by the mock, just don't do anything.
		},
		DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
			return r, err
		},
		GetPathsWithNoRetriesFunc: func() []string {
			return []string{"/healthcheck"}
		},
	}
}

func newIdentityAPIClient(_ *testing.T, httpClient *dphttp.ClienterMock) *Client {
	healthClient := healthcheck.NewClientWithClienter(service, testHost, httpClient)
	return NewWithHealthClient(healthClient)
}

func TestHealthCheckerClient(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	timePriorHealthCheck := time.Now().UTC()
	path := "/health"

	Convey("Given clienter.Do returns an error", t, func() {
		clientError := errors.New("unexpected error")
		httpClient := newMockHTTPClient(&http.Response{}, clientError)
		identityAPIClient := newIdentityAPIClient(t, httpClient)
		check := initialTestState

		Convey("When identity API client Checker is called", func() {
			err := identityAPIClient.Checker(ctx, &check)
			So(err, ShouldBeNil)

			Convey("Then the expected check is returned", func() {
				So(check.Name(), ShouldEqual, service)
				So(check.Status(), ShouldEqual, health.StatusCritical)
				So(check.StatusCode(), ShouldEqual, 0)
				So(check.Message(), ShouldEqual, clientError.Error())
				So(*check.LastChecked(), ShouldHappenAfter, timePriorHealthCheck)
				So(check.LastSuccess(), ShouldBeNil)
				So(*check.LastFailure(), ShouldHappenAfter, timePriorHealthCheck)
			})

			Convey("And client.Do should be called once with the expected parameters", func() {
				doCalls := httpClient.DoCalls()
				So(doCalls, ShouldHaveLength, 1)
				So(doCalls[0].Req.URL.Path, ShouldEqual, path)
			})
		})
	})

	Convey("Given a 500 response for health check", t, func() {
		httpClient := newMockHTTPClient(&http.Response{StatusCode: http.StatusInternalServerError}, nil)
		identityAPIClient := newIdentityAPIClient(t, httpClient)
		check := initialTestState

		Convey("When identity API client Checker is called", func() {
			err := identityAPIClient.Checker(ctx, &check)
			So(err, ShouldBeNil)

			Convey("Then the expected check is returned", func() {
				So(check.Name(), ShouldEqual, service)
				So(check.Status(), ShouldEqual, health.StatusCritical)
				So(check.StatusCode(), ShouldEqual, 500)
				So(check.Message(), ShouldEqual, service+healthcheck.StatusMessage[health.StatusCritical])
				So(*check.LastChecked(), ShouldHappenAfter, timePriorHealthCheck)
				So(check.LastSuccess(), ShouldBeNil)
				So(*check.LastFailure(), ShouldHappenAfter, timePriorHealthCheck)
			})

			Convey("And client.Do should be called once with the expected parameters", func() {
				doCalls := httpClient.DoCalls()
				So(doCalls, ShouldHaveLength, 1)
				So(doCalls[0].Req.URL.Path, ShouldEqual, path)
			})
		})
	})
}

func TestGetToken(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	tokensEndpoint := "/tokens"

	Convey("Given Get Token is returned successfully", t, func() {
		body := map[string]interface{}{
			"expirationTime":             defaultExpirationTime,
			"refreshTokenExpirationTime": defaultRefreshExpirationTime,
		}

		headers := http.Header{
			"Authorization": {defaultAuthorization},
			"Id":            {defaultId},
			"Refresh":       {defaultRefreshToken},
		}

		jsonBody, err := json.Marshal(&body)
		So(err, ShouldBeNil)

		httpClient := newMockHTTPClient(
			&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(jsonBody)),
				Header:     headers,
			},
			nil)

		identityAPIClient := newIdentityAPIClient(t, httpClient)

		Convey("When GetToken is called", func() {
			tokenResponse, err := identityAPIClient.GetToken(ctx, defaultCredentials)
			So(err, ShouldBeNil)

			Convey("Then the expected identity token is returned", func() {

				expectedTokenResponse := TokenResponse{
					Token:                      defaultAuthorization,
					RefreshToken:               defaultRefreshToken,
					ExpirationTime:             defaultExpirationTime,
					RefreshTokenExpirationTime: defaultRefreshExpirationTime,
				}

				So(*tokenResponse, ShouldResemble, expectedTokenResponse)

				Convey("And no error is returned", func() {
					So(err, ShouldBeNil)

					Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						So(doCalls, ShouldHaveLength, 1)
						So(doCalls[0].Req.URL.Path, ShouldEqual, tokensEndpoint)
						So(doCalls[0].Req.Method, ShouldEqual, http.MethodPost)
						expectedBody, _ := json.Marshal(defaultCredentials)
						actualBody, _ := io.ReadAll(doCalls[0].Req.Body)

						So(actualBody, ShouldResemble, expectedBody)
					})
				})
			})
		})
	})

	Convey("Given Get Token is returned with an error", t, func() {

		httpClient := newMockHTTPClient(
			&http.Response{
				StatusCode: http.StatusInternalServerError,
			},
			nil)

		identityAPIClient := newIdentityAPIClient(t, httpClient)

		Convey("When GetToken is called", func() {
			tokenResponse, err := identityAPIClient.GetToken(ctx, defaultCredentials)

			Convey("Then an error should be returned", func() {
				So(err.Status(), ShouldEqual, http.StatusInternalServerError)

				Convey("And the returned token response should be nil", func() {
					So(tokenResponse, ShouldBeNil)

					Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						So(doCalls, ShouldHaveLength, 1)
						So(doCalls[0].Req.URL.Path, ShouldEqual, tokensEndpoint)
						So(doCalls[0].Req.Method, ShouldEqual, http.MethodPost)
						expectedBody, _ := json.Marshal(defaultCredentials)
						actualBody, _ := io.ReadAll(doCalls[0].Req.Body)

						So(actualBody, ShouldResemble, expectedBody)
					})
				})
			})
		})
	})

	Convey("Given Get Token is returned with unauthorised", t, func() {

		httpClient := newMockHTTPClient(
			&http.Response{
				StatusCode: http.StatusUnauthorized,
			},
			nil)

		identityAPIClient := newIdentityAPIClient(t, httpClient)

		Convey("When GetToken is called", func() {
			tokenResponse, err := identityAPIClient.GetToken(ctx, defaultCredentials)

			Convey("Then an error should be returned", func() {
				So(err.Status(), ShouldEqual, http.StatusUnauthorized)

				Convey("And the returned token response should be nil", func() {
					So(tokenResponse, ShouldBeNil)

					Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						So(doCalls, ShouldHaveLength, 1)
						So(doCalls[0].Req.URL.Path, ShouldEqual, tokensEndpoint)
						So(doCalls[0].Req.Method, ShouldEqual, http.MethodPost)
						expectedBody, _ := json.Marshal(defaultCredentials)
						actualBody, _ := io.ReadAll(doCalls[0].Req.Body)

						So(actualBody, ShouldResemble, expectedBody)
					})
				})
			})
		})
	})

	Convey("Given Get Token is returned with a corrupted repsonse", t, func() {

		body := []byte("{}[")

		httpClient := newMockHTTPClient(
			&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(body)),
			},
			nil)

		identityAPIClient := newIdentityAPIClient(t, httpClient)

		Convey("When GetToken is called", func() {
			tokenResponse, err := identityAPIClient.GetToken(ctx, defaultCredentials)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)

				Convey("And the returned token response should be nil", func() {
					So(tokenResponse, ShouldBeNil)

					Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						So(doCalls, ShouldHaveLength, 1)
						So(doCalls[0].Req.URL.Path, ShouldEqual, tokensEndpoint)
						So(doCalls[0].Req.Method, ShouldEqual, http.MethodPost)
						expectedBody, _ := json.Marshal(defaultCredentials)
						actualBody, _ := io.ReadAll(doCalls[0].Req.Body)

						So(actualBody, ShouldResemble, expectedBody)
					})
				})
			})
		})
	})
}
