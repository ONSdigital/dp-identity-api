package main

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-identity-api/cognito"
	"github.com/ONSdigital/dp-identity-api/config"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/ONSdigital/dp-identity-api/service"
	"github.com/ONSdigital/log.go/log"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/pkg/errors"
	"os"
	"time"
)

func main() {
	ctx := context.Background()
	if err := runUserAndGroupsRemove(ctx); err != nil {
		log.Event(nil, "fatal runtime error", log.Error(err), log.FATAL)
		os.Exit(1)
	}
}

func runUserAndGroupsRemove(ctx context.Context) error {
	svcList := service.NewServiceList(&service.Init{})
	cfg, err := config.Get()
	if err != nil {
		return errors.Wrap(err, "error getting configuration")
	}
	cognitoClient := svcList.GetCognitoClient(cfg.AWSRegion)

	err = checkPoolExistsAndIsLocal(ctx, cognitoClient, cfg.AWSCognitoUserPoolID)
	if err != nil {
		return errors.Wrap(err, "error checking user pool details")
	}
	backoffSchedule := []time.Duration{
		1 * time.Second,
		3 * time.Second,
		10 * time.Second,
	}

	deleteUsers(ctx, cognitoClient, cfg.AWSCognitoUserPoolID, backoffSchedule)
	deleteGroups(ctx, cognitoClient, cfg.AWSCognitoUserPoolID, backoffSchedule)

	return nil
}

func deleteUsers(ctx context.Context, client cognito.Client, userPoolId string, backoffSchedule []time.Duration) {
	baseUsername := "test-user-"
	for i := range [UserCount]int{} {
		for _, backoff := range backoffSchedule {
			username := baseUsername + fmt.Sprint(i)
			userDeletionInput := cognitoidentityprovider.AdminDeleteUserInput{
				UserPoolId: &userPoolId,
				Username:   &username,
			}
			_, awsErr := client.AdminDeleteUser(&userDeletionInput)
			if awsErr != nil {
				err := models.NewCognitoError(ctx, awsErr, "AdminDeleteUser during dummy data creation")
				if err.Code != models.TooManyRequestsError {
					break
				}
			} else {
				break
			}
			time.Sleep(backoff)
		}
	}
}

func deleteGroups(ctx context.Context, client cognito.Client, userPoolId string, backoffSchedule []time.Duration) {
	baseGroupName := "test-group-"

	for i := range [GroupCount]int{} {
		for _, backoff := range backoffSchedule {
			groupName := baseGroupName + fmt.Sprint(i)
			groupDeletionInput := cognitoidentityprovider.DeleteGroupInput{
				GroupName:  &groupName,
				UserPoolId: &userPoolId,
			}
			_, awsErr := client.DeleteGroup(&groupDeletionInput)
			if awsErr != nil {
				err := models.NewCognitoError(ctx, awsErr, "AdminDeleteUser during dummy data creation")
				if err.Code != models.TooManyRequestsError {
					break
				}
			} else {
				break
			}
			time.Sleep(backoff)
		}
	}
}
