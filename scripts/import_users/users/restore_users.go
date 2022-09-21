package users

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"

	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/ONSdigital/dp-identity-api/scripts/import_users/config"
	"github.com/ONSdigital/dp-identity-api/scripts/utils"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/pkg/errors"
)

const suppressEmails = "SUPPRESS"

func ImportUsersFromS3(ctx context.Context, config *config.Config) error {
	log.Info(ctx, fmt.Sprintf("started restoring users to cognito from s3 file: %v", config.GetS3UsersFilePath()))

	s3FileReader := utils.S3Reader{}
	responseBody := s3FileReader.GetS3Reader(ctx, config.AWSProfile, config.S3Region, config.S3Bucket, config.GetS3UsersFilePath())
	defer responseBody.Close()
	reader := csv.NewReader(responseBody)
	client := utils.GetCognitoClient(config.AWSProfile, config.S3Region)

	// Extract column indexes from header line
	cols, err := reader.Read()
	if err != nil {
		return errors.New("unable to read header from users backup file")
	}
	colsMap := make(map[string]int, len(cols))
	for i, col := range cols {
		colsMap[col] = i
	}
	for _, wanted := range []string{"given_name", "family_name", "email", "cognito:username"} {
		if _, ok := colsMap[wanted]; !ok {
			return errors.New(fmt.Sprintf("column '%s' not found in users backup file", wanted))
		}
	}

	var count int
	for {
		line, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		count += 1

		userInfo, err := createUser(ctx, client, line, count, config, colsMap)
		if err != nil {
			log.Error(ctx, fmt.Sprintf("failed to create user. Processline %v user: %v", count, userInfo), err)
			continue
		}

		// Update password to be permanent if permanent passwrod env variable set to true,
		// allows user to reset password
		if config.PermanentPassword {
			if err = makeUserPasswordPermanent(ctx, config, client, userInfo); err != nil {
				log.Error(ctx, fmt.Sprintf("failed to make user password permanent. Processline %v user: %v", count, userInfo), err)
				continue
			}
		}

		// Disable user if it's 'enabled' column is not TRUE
		enabledCol, ok := colsMap["enabled"]
		if ok && len(line) > enabledCol && line[enabledCol] == "false" {
			if err = disableUser(ctx, config, client, userInfo); err != nil {
				log.Error(ctx, fmt.Sprintf("failed to disable user. Processline %v user: %v", count, userInfo), err)
				continue
			}
		}
	}

	log.Info(ctx, "Successfully processed all the users in S3 file")

	return nil
}

func createUser(ctx context.Context, client *cognitoidentityprovider.CognitoIdentityProvider, line []string, lineNumber int, config *config.Config, colsMap map[string]int) (*models.UserParams, error) {
	userInfo := models.UserParams{
		Forename: line[colsMap["given_name"]],
		Lastname: line[colsMap["family_name"]],
		Email:    line[colsMap["email"]],
		ID:       line[colsMap["cognito:username"]],
	}
	userInfo.GeneratePassword(ctx)

	createUserRequest := userInfo.BuildCreateUserRequest(userInfo.ID, config.AWSCognitoUserPoolID, config.MessageAction)
	_, err := client.AdminCreateUser(createUserRequest)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func makeUserPasswordPermanent(ctx context.Context, config *config.Config, client *cognitoidentityprovider.CognitoIdentityProvider, userInfo *models.UserParams) (err error) {
	perm := true
	user := &cognitoidentityprovider.AdminSetUserPasswordInput{
		Password:   &userInfo.Password,
		Permanent:  &perm,
		UserPoolId: &config.AWSCognitoUserPoolID,
		Username:   &userInfo.ID,
	}
	_, err = client.AdminSetUserPassword(user)

	return
}

func disableUser(ctx context.Context, config *config.Config, client *cognitoidentityprovider.CognitoIdentityProvider, userInfo *models.UserParams) (err error) {
	userDisableRequest := userInfo.BuildDisableUserRequest(config.AWSCognitoUserPoolID)
	_, err = client.AdminDisableUser(userDisableRequest)

	return
}
