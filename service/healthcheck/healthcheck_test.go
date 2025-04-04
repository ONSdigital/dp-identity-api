package healthcheck_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/smithy-go"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-identity-api/v2/cognito/mock"
	"github.com/ONSdigital/dp-identity-api/v2/models"
	"github.com/ONSdigital/dp-identity-api/v2/service/healthcheck"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetHealthCheck(t *testing.T) {
	ctx := context.Background()

	m := &mock.MockCognitoIdentityProviderClient{}

	m.DescribeUserPoolFunc = func(_ context.Context, poolInputData *cognitoidentityprovider.DescribeUserPoolInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.DescribeUserPoolOutput, error) {
		exsitingPoolID := "us-west-2_aaaaaaaaa"
		if *poolInputData.UserPoolId != exsitingPoolID {
			return nil, errors.New("Failed to load user pool data")
		}
		return nil, nil
	}

	Convey("dp-identity-api healthchecker reports healthy", t, func() {
		awsUserPoolID := "us-west-2_aaaaaaaaa"
		checkState := health.NewCheckState("dp-identity-api-test")

		m.GetGroupFunc = func(_ context.Context, _ *cognitoidentityprovider.GetGroupInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.GetGroupOutput, error) {
			group := &cognitoidentityprovider.GetGroupOutput{
				Group: &types.GroupType{},
			}
			return group, nil
		}
		checker := healthcheck.CognitoHealthCheck(ctx, m, &awsUserPoolID)
		err := checker(ctx, checkState)
		Convey("When GetHealthCheck is called", func() {
			Convey("Then the HealthCheck flag is set to true and HealthCheck is returned", func() {
				So(checkState.StatusCode(), ShouldEqual, http.StatusOK)
				So(checkState.Status(), ShouldEqual, health.StatusOK)
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("dp-identity-api healthchecker reports critical", t, func() {
		Convey("When the user pool can't be found", func() {
			m.GetGroupFunc = func(_ context.Context, _ *cognitoidentityprovider.GetGroupInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.GetGroupOutput, error) {
				group := &cognitoidentityprovider.GetGroupOutput{
					Group: &types.GroupType{},
				}
				return group, nil
			}

			awsUserPoolID := "us-west-2_aaaaaaaab"

			checkState := health.NewCheckState("dp-identity-api-test")

			checker := healthcheck.CognitoHealthCheck(ctx, m, &awsUserPoolID)
			err := checker(ctx, checkState)
			Convey("When GetHealthCheck is called", func() {
				Convey("Then the HealthCheck flag is set to true and HealthCheck is returned", func() {
					So(checkState.StatusCode(), ShouldEqual, http.StatusTooManyRequests)
					So(checkState.Status(), ShouldEqual, health.StatusCritical)
					So(err, ShouldNotBeNil)
				})
			})
		})

		Convey("When a role group can't be found", func() {
			awsUserPoolID := "us-west-2_aaaaaaaaa"

			Convey("the admin role group is missing", func() {
				m.GetGroupFunc = func(_ context.Context, inputData *cognitoidentityprovider.GetGroupInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.GetGroupOutput, error) {
					if *inputData.GroupName == models.AdminRoleGroup {
						awsErrCode := "ResourceNotFoundException"
						awsErrMessage := "Group not found."
						awsErr := &smithy.GenericAPIError{
							Code:    awsErrCode,
							Message: awsErrMessage,
						}
						return nil, awsErr
					}
					group := &cognitoidentityprovider.GetGroupOutput{
						Group: &types.GroupType{},
					}
					return group, nil
				}

				checkState := health.NewCheckState("dp-identity-api-test")

				checker := healthcheck.CognitoHealthCheck(ctx, m, &awsUserPoolID)
				err := checker(ctx, checkState)
				Convey("When GetHealthCheck is called", func() {
					Convey("Then the HealthCheck flag is set to true and HealthCheck is returned", func() {
						So(checkState.StatusCode(), ShouldEqual, http.StatusTooManyRequests)
						So(checkState.Status(), ShouldEqual, health.StatusCritical)
						So(err, ShouldNotBeNil)
					})
				})
			})

			Convey("the publisher role group is missing", func() {
				m.GetGroupFunc = func(_ context.Context, inputData *cognitoidentityprovider.GetGroupInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.GetGroupOutput, error) {
					if *inputData.GroupName == models.PublisherRoleGroup {
						awsErrCode := "ResourceNotFoundException"
						awsErrMessage := "Group not found."
						awsErr := &smithy.GenericAPIError{
							Code:    awsErrCode,
							Message: awsErrMessage,
						}
						return nil, awsErr
					}
					group := &cognitoidentityprovider.GetGroupOutput{
						Group: &types.GroupType{},
					}
					return group, nil
				}

				checkState := health.NewCheckState("dp-identity-api-test")

				checker := healthcheck.CognitoHealthCheck(ctx, m, &awsUserPoolID)
				err := checker(ctx, checkState)
				Convey("When GetHealthCheck is called", func() {
					Convey("Then the HealthCheck flag is set to true and HealthCheck is returned", func() {
						So(checkState.StatusCode(), ShouldEqual, http.StatusTooManyRequests)
						So(checkState.Status(), ShouldEqual, health.StatusCritical)
						So(err, ShouldNotBeNil)
					})
				})
			})
		})
	})
}
