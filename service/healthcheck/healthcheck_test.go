package healthcheck_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-identity-api/cognito/mock"
	healthcheck "github.com/ONSdigital/dp-identity-api/service/healthcheck"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetHealthCheck(t *testing.T) {
	ctx := context.Background()

	m := &mock.MockCognitoIdentityProviderClient{}

	m.DescribeUserPoolFunc = func(poolInputData *cognitoidentityprovider.DescribeUserPoolInput) (*cognitoidentityprovider.DescribeUserPoolOutput, error) {
		exsitingPoolID := "us-west-2_aaaaaaaaa"
		if *poolInputData.UserPoolId != exsitingPoolID {
			return nil, errors.New("Failed to load user pool data")
		} 
		return nil, nil
	}

	Convey("dp-identity-api healthchecker reports healthy", t, func() {
		awsUserPoolID := "us-west-2_aaaaaaaaa"
		checkState := health.NewCheckState("dp-identity-api-test")
	
		checker := healthcheck.CognitoHealthCheck(m, &awsUserPoolID)
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
	
		checker := healthcheck.CognitoHealthCheck(m, &awsUserPoolID)
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