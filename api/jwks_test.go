package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-identity-api/jwks"
	"github.com/ONSdigital/dp-identity-api/jwks/mock"
	"github.com/ONSdigital/dp-identity-api/models"
)

const(
	jwksEndpoint = "/v1/jwt-keys"
	errorMessage = "Random Error!"
)

func TestCognitoPoolJWKSHandler(t *testing.T) {
	api, w, _ := apiSetup()
	r := httptest.NewRequest(http.MethodGet, jwksEndpoint, nil)

	resp, err := api.CognitoPoolJWKSHandler(context.Background(), w, r)

	Convey("Request json web key set - success", t, func() {
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		So(resp.Status, ShouldEqual, http.StatusOK)
		So(string(resp.Body), ShouldResemble, mock.JWKSResponseData)
	})
}

func TestCognitoPoolJWKSHandlerErrors404(t *testing.T) {
	api, w, _ := apiSetup()

	// override default jwks handler
	api.JWKSHandler = &mock.JWKSIntMock{
		JWKSGetKeysetFunc: func(awsRegion string, poolId string) (*jwks.JWKS, error) {
			return nil, errors.New(errorMessage)
		},
	}

	r := httptest.NewRequest(http.MethodGet, jwksEndpoint, nil)

	resp, err := api.CognitoPoolJWKSHandler(context.Background(), w, r)

	Convey("Request json web key set - 404 response", t, func() {
		So(resp, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Status, ShouldEqual, http.StatusNotFound)
		So(err.Errors[0].Error(), ShouldResemble, errorMessage)
	})
}

func TestCognitoPoolJWKSHandlerErrors500(t *testing.T) {
	api, w, _ := apiSetup()

	// override default jwks handler
	api.JWKSHandler = &mock.JWKSIntMock{
		JWKSGetKeysetFunc: func(awsRegion string, poolId string) (*jwks.JWKS, error) {
			return &jwks.JWKS{
				Keys: []jwks.JsonKey{
					mock.KeySetOne,
					mock.KeySetTwo,
				},
			}, nil
		},
		JWKSToRSAJSONResponseFunc: func(jwksMoqParam *jwks.JWKS) ([]byte, error) {
			return nil, errors.New(errorMessage)
		},
	}

	r := httptest.NewRequest(http.MethodGet, jwksEndpoint, nil)

	resp, err := api.CognitoPoolJWKSHandler(context.Background(), w, r)

	Convey("Request json web key set - 500 response", t, func() {
		So(resp, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Status, ShouldEqual, http.StatusInternalServerError)
		// only one error
		castErr := err.Errors[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.JWKSParseError)
		So(castErr.Description, ShouldEqual, models.JWKSParseErrorDescription)
	})
}
