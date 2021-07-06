package healthcheck

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/models"
	"net/http"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/log"

	cognitoclient "github.com/ONSdigital/dp-identity-api/cognito"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

const CognitoHealthy = "Cognito Healthy"

func CognitoHealthCheck(cognitoClient cognitoclient.Client, userPoolID *string) health.Checker {

	return func(ctx context.Context, state *health.CheckState) error {
		_, err := cognitoClient.DescribeUserPool(&cognito.DescribeUserPoolInput{UserPoolId: userPoolID})

		if err != nil {
			state.Update(health.StatusCritical, err.Error(), http.StatusTooManyRequests)
			// log the error
			log.Event(context.Background(), "Error running identity service healthcheck", log.Error(err), log.ERROR)
			return err
		}

		adminGroupDetails := models.NewAdminRoleGroup()
		adminGroupRequest := adminGroupDetails.BuildGetGroupRequest(*userPoolID)
		_, err = cognitoClient.GetGroup(adminGroupRequest)

		if err != nil {
			state.Update(health.StatusCritical, err.Error(), http.StatusTooManyRequests)
			// log the error
			log.Event(context.Background(), "Error running identity service healthcheck", log.Error(err), log.ERROR)
			return err
		}

		publisherGroupDetails := models.NewPublisherRoleGroup()
		publisherGroupRequest := publisherGroupDetails.BuildGetGroupRequest(*userPoolID)
		_, err = cognitoClient.GetGroup(publisherGroupRequest)

		if err != nil {
			state.Update(health.StatusCritical, err.Error(), http.StatusTooManyRequests)
			// log the error
			log.Event(context.Background(), "Error running identity service healthcheck", log.Error(err), log.ERROR)
			return err
		}

		state.Update(health.StatusOK, CognitoHealthy, http.StatusOK)

		return nil
	}
}
