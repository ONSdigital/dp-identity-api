package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

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
