package healthcheck_test

import (
	"context"
	"net/http"
	"testing"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	cognitomock "github.com/ONSdigital/dp-identity-api/cognitoclient/mock"
	healthcheck "github.com/ONSdigital/dp-identity-api/service/healthcheck"

	. "github.com/smartystreets/goconvey/convey"
)


func TestGetHealthCheck(t *testing.T) {
	ctx := context.Background()
	Convey("dp-identity-api healthchecker reports healthy", t, func() {
		awsUserPoolID := "us-west-2_aaaaaaaaa"
		checkState := health.NewCheckState("dp-identity-api-test")
	
		checker := healthcheck.CognitoHealthCheck(&cognitomock.MockCognitoIdentityProviderClient{}, &awsUserPoolID)
		err := checker(ctx,checkState)
		Convey("When GetHealthCheck is called", func() {
			Convey("Then the HealthCheck flag is set to true and HealthCheck is returned", func() {
				So(checkState.StatusCode(), ShouldEqual, http.StatusOK)
				So(checkState.Status(), ShouldEqual, health.StatusOK)
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("dp-identity-api healthchecker reports critical", t, func() {
		awsUserPoolID := "us-west-2_aaaaaaaab"
			
		checkState := health.NewCheckState("dp-identity-api-test")
	
		checker := healthcheck.CognitoHealthCheck(&cognitomock.MockCognitoIdentityProviderClient{}, &awsUserPoolID)
		err := checker(ctx,checkState)
		Convey("When GetHealthCheck is called", func() {
			Convey("Then the HealthCheck flag is set to true and HealthCheck is returned", func() {
				So(checkState.StatusCode(), ShouldEqual, http.StatusTooManyRequests)
				So(checkState.Status(), ShouldEqual, health.StatusCritical)
				So(err, ShouldNotBeNil)
			})
		})
	})
}