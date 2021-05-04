package api

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-identity-api/cognito/mock"
)

func TestSetup(t *testing.T) {
	Convey("Given an API instance", t, func() {
		r := mux.NewRouter()
		ctx := context.Background()
		api := Setup(ctx, r, &mock.MockCognitoIdentityProviderClient{})

		Convey("When created the following route(s) should have been added", func() {
			So(hasRoute(api.Router, "/hello", "GET"), ShouldBeTrue)
			So(hasRoute(api.Router, "/tokens", "POST"), ShouldBeTrue)
			So(hasRoute(api.Router, "/users", "POST"), ShouldBeTrue)
		})

		Convey("Ensure cognito client has been added to api", func() {
			So(api.CognitoClient, ShouldNotBeNil)
		})
	})
}

func hasRoute(r *mux.Router, path, method string) bool {
	req := httptest.NewRequest(method, path, nil)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}
