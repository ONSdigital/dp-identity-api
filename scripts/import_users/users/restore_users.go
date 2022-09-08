package users

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/ONSdigital/dp-identity-api/scripts/import_users/config"
	"github.com/ONSdigital/dp-identity-api/scripts/utils"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/pkg/errors"
	"io"
)

func ImportUsersFromS3(ctx context.Context, config *config.Config) error {
	log.Info(ctx, fmt.Sprintf("started restoring users to cognito from s3 file: %v", config.GetS3UsersFilePath()))

	s3FileReader := utils.S3Reader{}
	responseBody := s3FileReader.GetS3Reader(ctx, config.S3Region, config.S3Bucket, config.GetS3UsersFilePath())
	defer responseBody.Close()
	reader := csv.NewReader(responseBody)
	client := utils.GetCognitoClient(config.S3Region)

	// Extract column indexes from header line
	cols, err := reader.Read()
	if err != nil {
		return errors.New("unable to read header from user backup file")
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

	count := 1
	for {
		line, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		createUser(ctx, client, line, count, config, colsMap)
		count += 1
	}
	log.Info(ctx, "Successfully processed all the users in S3 file")
	return nil
}

func createUser(ctx context.Context, client *cognitoidentityprovider.CognitoIdentityProvider, line []string, lineNumber int, config *config.Config, colsMap map[string]int) {
	userInfo := models.UserParams{
		Forename: line[colsMap["given_name"]],
		Lastname: line[colsMap["family_name"]],
		Email:    line[colsMap["email"]],
		ID:       line[colsMap["cognito:username"]],
	}
	userInfo.GeneratePassword(ctx)

	createUserRequest := userInfo.BuildCreateUserRequest(userInfo.ID, config.AWSCognitoUserPoolID)
	_, err := client.AdminCreateUser(createUserRequest)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed to processline %v user: %+v", lineNumber, userInfo), err)
	}

	//Disable user if it's 'enabled' column is not TRUE
	enabledCol, ok := colsMap["enabled"]
	if ok && line[enabledCol] == "false" {
		userDisableRequest := userInfo.BuildDisableUserRequest(config.AWSCognitoUserPoolID)
		if _, err = client.AdminDisableUser(userDisableRequest); err != nil {
			log.Error(ctx, fmt.Sprintf("failed to disable user: %+v", userInfo), err)
		}
	}

}
