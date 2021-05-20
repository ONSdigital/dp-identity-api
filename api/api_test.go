package api

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-identity-api/apierrorsdeprecated"
	"github.com/ONSdigital/dp-identity-api/cognito/mock"
)

func TestSetup(t *testing.T) {
	Convey("Given an API instance", t, func() {
		r := mux.NewRouter()
		ctx := context.Background()
		api, err := Setup(ctx, r, &mock.MockCognitoIdentityProviderClient{}, "us-west-2_aaaaaaaaa", "client-aaa-bbb", "secret-ccc-ddd", "authflow")

		Convey("When created the following route(s) should have been added", func() {
			So(hasRoute(api.Router, "/tokens", "POST"), ShouldBeTrue)
			So(hasRoute(api.Router, "/users", "POST"), ShouldBeTrue)
		})

		Convey("No error returned when user pool id supplied", func() {
			So(err, ShouldBeNil)
		})

		Convey("Ensure cognito client has been added to api", func() {
			So(api.CognitoClient, ShouldNotBeNil)
		})
	})

	Convey("Given an API instance with an empty required parameter passed", t, func() {
		paramCheckTests := []struct {
			testName       string
			userPoolId     string
			clientId       string
			clientSecret   string
			clientAuthFlow string
		}{
			// missing userPoolId
			{
				"missing userPoolId",
				"",
				"client-aaa-bbb",
				"secret-ccc-ddd",
				"authflow",
			},
			// missing clientId
			{
				"missing clientId",
				"eu-west-22_bdsjhids2",
				"",
				"secret-ccc-ddd",
				"authflow",
			},
			// missing clientSecret
			{
				"missing clientSecret",
				"eu-west-22_bdsjhids2",
				"client-aaa-bbb",
				"",
				"authflow",
			},
			// missing clientAuthFlow
			{
				"missing clientAuthFlow",
				"eu-west-22_bdsjhids2",
				"client-aaa-bbb",
				"secret-ccc-ddd",
				"",
			},
		}

		for _, tt := range paramCheckTests {
			r := mux.NewRouter()
			ctx := context.Background()
			_, err := Setup(ctx, r, &mock.MockCognitoIdentityProviderClient{}, tt.userPoolId, tt.clientId, tt.clientSecret, tt.clientAuthFlow)

			Convey("Error should not be nil if require parameter is empty: "+tt.testName, func() {
				So(err.Error(), ShouldEqual, apierrorsdeprecated.RequiredParameterNotFoundDescription)
			})
		}
	})
}

func hasRoute(r *mux.Router, path, method string) bool {
	req := httptest.NewRequest(method, path, nil)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}
