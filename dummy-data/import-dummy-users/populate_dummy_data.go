package main

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"

	"github.com/ONSdigital/dp-identity-api/v2/cognito"
	"github.com/ONSdigital/dp-identity-api/v2/config"
	"github.com/ONSdigital/dp-identity-api/v2/models"
	"github.com/ONSdigital/dp-identity-api/v2/service"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/pkg/errors"
)

const LocalUserPoolName = "local-florence-users"

const UserCount = 2000
const GroupCount = 200
const GroupUserCount = 500

func main() {
	ctx := context.Background()
	if err := runUserAndGroupsPopulate(ctx); err != nil {
		log.Fatal(ctx, "fatal runtime error", err)
	}
}

func runUserAndGroupsPopulate(ctx context.Context) error {
	svcList := service.NewServiceList(&service.Init{})
	cfg, err := config.Get()
	if err != nil {
		return errors.Wrap(err, "error getting configuration")
	}
	cognitoClient := svcList.GetCognitoClient(ctx, cfg.AWSRegion)

	err = checkPoolExistsAndIsLocal(ctx, cognitoClient, cfg.AWSCognitoUserPoolID)
	if err != nil {
		return errors.Wrap(err, "error checking user pool details")
	}

	backoffSchedule := []time.Duration{
		1 * time.Second,
		3 * time.Second,
		10 * time.Second,
	}

	createUsers(ctx, cognitoClient, cfg.AWSCognitoUserPoolID, backoffSchedule)
	confirmUsers(ctx, cognitoClient, cfg.AWSCognitoUserPoolID, backoffSchedule)
	disableUsers(ctx, cognitoClient, cfg.AWSCognitoUserPoolID, backoffSchedule)
	createGroups(ctx, cognitoClient, cfg.AWSCognitoUserPoolID, backoffSchedule)
	addUsersToGroups(ctx, cognitoClient, cfg.AWSCognitoUserPoolID, backoffSchedule)

	return nil
}

func checkPoolExistsAndIsLocal(ctx context.Context, client cognito.Client, userPoolID string) error {
	input := cognitoidentityprovider.DescribeUserPoolInput{
		UserPoolId: aws.String(userPoolID),
	}
	userPoolDetails, err := client.DescribeUserPool(ctx, &input)
	if err != nil {
		return models.NewCognitoError(ctx, err, "loading User Pool details for dummy data population")
	}
	if *userPoolDetails.UserPool.Name != LocalUserPoolName {
		return models.NewValidationError(ctx, models.InvalidUserPoolError, models.InvalidUserPoolDescription)
	}
	return nil
}

func createUsers(ctx context.Context, client cognito.Client, userPoolID string, backoffSchedule []time.Duration) {
	var (
		baseFirstName, baseLastName, emailDomain = "test", "user-", "@ons.gov.uk"
	)
	for i := range [UserCount]int{} {
		for _, backoff := range backoffSchedule {
			user := models.UserParams{}
			var passwordError error
			user.Password, passwordError = user.GeneratePassword(ctx)
			if passwordError != nil {
				break
			}
			lastName := baseLastName + fmt.Sprint(i)
			userID := baseFirstName + "-" + lastName
			userCreationInput := cognitoidentityprovider.AdminCreateUserInput{
				UserAttributes: []types.AttributeType{
					{
						Name:  aws.String("given_name"),
						Value: &baseFirstName,
					},
					{
						Name:  aws.String("family_name"),
						Value: aws.String(lastName),
					},
					{
						Name:  aws.String("email"),
						Value: aws.String(baseFirstName + "." + lastName + emailDomain),
					},
					{
						Name:  aws.String("email_verified"),
						Value: aws.String("true"),
					},
				},
				MessageAction:     types.MessageActionTypeSuppress,
				TemporaryPassword: &user.Password,
				UserPoolId:        &userPoolID,
				Username:          &userID,
			}
			_, awsErr := client.AdminCreateUser(ctx, &userCreationInput)
			if awsErr != nil {
				err := models.NewCognitoError(ctx, awsErr, "AdminCreateUser during dummy data creation")
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

func confirmUsers(ctx context.Context, client cognito.Client, userPoolID string, backoffSchedule []time.Duration) {
	var (
		baseEmailPrefix, emailDomain = "test.user-", "@ons.gov.uk"
	)
	for i := range [UserCount]int{} {
		if math.Mod(float64(i), float64(11)) == 0 {
			continue
		}
		for _, backoff := range backoffSchedule {
			user := models.UserParams{}
			var passwordError error
			user.Password, passwordError = user.GeneratePassword(ctx)
			if passwordError != nil {
				break
			}
			userSetPasswordInput := cognitoidentityprovider.AdminSetUserPasswordInput{
				Password:   &user.Password,
				Permanent:  true,
				UserPoolId: &userPoolID,
				Username:   aws.String(baseEmailPrefix + fmt.Sprint(i) + emailDomain),
			}
			_, awsErr := client.AdminSetUserPassword(ctx, &userSetPasswordInput)
			if awsErr != nil {
				err := models.NewCognitoError(ctx, awsErr, "AdminSetUserPassword during dummy data creation")
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

func disableUsers(ctx context.Context, client cognito.Client, userPoolID string, backoffSchedule []time.Duration) {
	var (
		baseFirstName, baseLastName, emailDomain = "test", "user-", "@ons.gov.uk"
	)
	for i := range [UserCount]int{} {
		if math.Mod(float64(i), float64(51)) != 0 {
			continue
		}
		for _, backoff := range backoffSchedule {
			lastName := baseLastName + fmt.Sprint(i)
			userDisableInput := cognitoidentityprovider.AdminDisableUserInput{
				UserPoolId: &userPoolID,
				Username:   aws.String(baseFirstName + "." + lastName + emailDomain),
			}
			_, awsErr := client.AdminDisableUser(ctx, &userDisableInput)
			if awsErr != nil {
				err := models.NewCognitoError(ctx, awsErr, "AdminDisableUser during dummy data creation")
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

func createGroups(ctx context.Context, client cognito.Client, userPoolID string, backoffSchedule []time.Duration) {
	var (
		baseGroupName, baseDescription = "test-group-", "Test Group "
	)
	for i := range [GroupCount]int{} {
		for _, backoff := range backoffSchedule {
			groupCreationInput := cognitoidentityprovider.CreateGroupInput{
				Description: aws.String(baseDescription + fmt.Sprint(i)),
				GroupName:   aws.String(baseGroupName + fmt.Sprint(i)),
				Precedence:  aws.Int32(3),
				UserPoolId:  &userPoolID,
			}
			_, awsErr := client.CreateGroup(ctx, &groupCreationInput)
			if awsErr != nil {
				err := models.NewCognitoError(ctx, awsErr, "CreateGroup during dummy data creation")
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

func addUsersToGroups(ctx context.Context, client cognito.Client, userPoolID string, backoffSchedule []time.Duration) {
	for userNumber := range [UserCount]int{} {
		if userNumber > GroupUserCount {
			break
		}
		for groupNumber := range [GroupCount]int{} {
			if userNumber == groupNumber {
				addUserToGroup(ctx, fmt.Sprint(userNumber), fmt.Sprint(groupNumber), userPoolID, client, backoffSchedule)
			}
			if math.Mod(float64(userNumber), float64(2)) == 0 && math.Mod(float64(groupNumber), float64(2)) == 0 {
				addUserToGroup(ctx, fmt.Sprint(userNumber), fmt.Sprint(groupNumber), userPoolID, client, backoffSchedule)
			}
			if math.Mod(float64(userNumber), float64(3)) == 0 && math.Mod(float64(groupNumber), float64(3)) == 0 {
				addUserToGroup(ctx, fmt.Sprint(userNumber), fmt.Sprint(groupNumber), userPoolID, client, backoffSchedule)
			}
		}
	}
}

func addUserToGroup(ctx context.Context, userNumber, groupNumber, userPoolID string, client cognito.Client, backoffSchedule []time.Duration) {
	var (
		baseFirstName, baseLastName, emailDomain = "test", "user-", "@ons.gov.uk"
		baseGroupName                            = "test-group-"
	)
	for _, backoff := range backoffSchedule {
		lastName := baseLastName + userNumber
		userAddToGroupInput := cognitoidentityprovider.AdminAddUserToGroupInput{
			GroupName:  aws.String(baseGroupName + groupNumber),
			UserPoolId: &userPoolID,
			Username:   aws.String(baseFirstName + "." + lastName + emailDomain),
		}
		_, awsErr := client.AdminAddUserToGroup(ctx, &userAddToGroupInput)
		if awsErr != nil {
			err := models.NewCognitoError(ctx, awsErr, "AdminAddUserToGroup during dummy data creation")
			if err.Code != models.TooManyRequestsError {
				break
			}
		} else {
			break
		}
		time.Sleep(backoff)
	}
}
