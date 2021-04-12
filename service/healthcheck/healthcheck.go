package healthcheck

import (
	"context"
	"net/http"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"

	cognitoclient "github.com/ONSdigital/dp-identity-api/cognitoclient"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

const CognitoHealthy  = "Cognito Healthy"

func CognitoHealthCheck(cognitoClient cognitoclient.Client, userPoolID *string) health.Checker {
	
	return func(ctx context.Context, state *health.CheckState) error {
		_, err := cognitoClient.DescribeUserPool(&cognito.DescribeUserPoolInput{UserPoolId: userPoolID})
	
		if err != nil {
			state.Update(health.StatusCritical, err.Error(), http.StatusTooManyRequests)
			return err
		}
		state.Update(health.StatusOK, CognitoHealthy, http.StatusOK)
	
		return nil
	}
}
