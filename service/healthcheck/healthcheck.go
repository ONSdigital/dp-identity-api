package healthcheck

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-identity-api/models"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/v2/log"

	cognitoclient "github.com/ONSdigital/dp-identity-api/cognito"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

const CognitoHealthy = "Cognito Healthy"

func CognitoHealthCheck(_ context.Context, cognitoClient cognitoclient.Client, userPoolID *string) health.Checker {
	return func(ctx context.Context, state *health.CheckState) error {
		_, err := cognitoClient.DescribeUserPool(&cognito.DescribeUserPoolInput{UserPoolId: userPoolID})

		if err != nil {
			if stateErr := state.Update(health.StatusCritical, err.Error(), http.StatusTooManyRequests); stateErr != nil {
				log.Error(ctx, "Error updating state during identity service healthcheck", stateErr)
			}
			// log the error
			log.Error(ctx, "Error running identity service healthcheck", err)
			return err
		}

		adminGroupDetails := models.NewAdminRoleGroup()
		adminGroupRequest := adminGroupDetails.BuildGetGroupRequest(*userPoolID)
		_, err = cognitoClient.GetGroup(adminGroupRequest)

		if err != nil {
			if stateErr := state.Update(health.StatusCritical, err.Error(), http.StatusTooManyRequests); stateErr != nil {
				log.Error(ctx, "Error updating state during identity service healthcheck", stateErr)
			}
			// log the error
			log.Error(ctx, "Error running identity service healthcheck", err)
			return err
		}

		publisherGroupDetails := models.NewPublisherRoleGroup()
		publisherGroupRequest := publisherGroupDetails.BuildGetGroupRequest(*userPoolID)
		_, err = cognitoClient.GetGroup(publisherGroupRequest)

		if err != nil {
			if stateErr := state.Update(health.StatusCritical, err.Error(), http.StatusTooManyRequests); stateErr != nil {
				log.Error(ctx, "Error updating state during identity service healthcheck", stateErr)
			}
			// log the error
			log.Error(ctx, "Error running identity service healthcheck", err)
			return err
		}

		if stateErr := state.Update(health.StatusOK, CognitoHealthy, http.StatusOK); stateErr != nil {
			log.Error(ctx, "Error updating state during identity service healthcheck", stateErr)
		}

		return nil
	}
}
