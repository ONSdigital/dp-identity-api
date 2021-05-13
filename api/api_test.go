package api

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-identity-api/api/content"
	"github.com/ONSdigital/dp-identity-api/cognito/mock"
)

func TestSetup(t *testing.T) {
	Convey("Given an API instance", t, func() {
		r := mux.NewRouter()
		ctx := context.Background()
		api, err := Setup(ctx, r, &mock.MockCognitoIdentityProviderClient{}, "us-west-2_aaaaaaaaa", "client-aaa-bbb", "secret-ccc-ddd")

		Convey("When created the following route(s) should have been added", func() {
			So(hasRoute(api.Router, "/hello", "GET"), ShouldBeTrue)
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

	Convey("Given an API instance with an empty user pool id", t, func() {
		r := mux.NewRouter()
		ctx := context.Background()
		api, err := Setup(ctx, r, &mock.MockCognitoIdentityProviderClient{}, "", "client-aaa-bbb", "secret-ccc-ddd")

		Convey("Error should not be nil if user pool id is empty", func() {
			So(api, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, content.UserPoolIdNotFoundMessage)
		})

	})
}

func hasRoute(r *mux.Router, path, method string) bool {
	req := httptest.NewRequest(method, path, nil)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}
