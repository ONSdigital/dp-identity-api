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

func ImportUsersFromS3(ctx context.Context, cfg *config.Config) error {
	log.Info(ctx, fmt.Sprintf("started restoring users to cognito from s3 file: %v", cfg.GetS3UsersFilePath()))

	s3FileReader := utils.S3Reader{}
	responseBody := s3FileReader.GetS3Reader(ctx, cfg.AWSProfile, cfg.S3Region, cfg.S3Bucket, cfg.GetS3UsersFilePath())
	defer responseBody.Close()
	reader := csv.NewReader(responseBody)
	client := utils.GetCognitoClient(cfg.AWSProfile, cfg.S3Region)

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

	count := 1
	for {
		line, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		createUser(ctx, client, line, count, cfg, colsMap)
		count += 1
	}
	log.Info(ctx, "Successfully processed all the users in S3 file")
	return nil
}

func createUser(ctx context.Context, client *cognitoidentityprovider.CognitoIdentityProvider, line []string, lineNumber int, cfg *config.Config, colsMap map[string]int) {
	userInfo := models.UserParams{
		Forename: line[colsMap["given_name"]],
		Lastname: line[colsMap["family_name"]],
		Email:    line[colsMap["email"]],
		ID:       line[colsMap["cognito:username"]],
	}

	if err := userInfo.GeneratePassword(ctx); err != nil {
		log.Error(ctx, "failed to generate password for user", err, log.Data{
			"email": userInfo.Email,
		})
		return
	}

	createUserRequest := userInfo.BuildCreateUserRequest(userInfo.ID, cfg.AWSCognitoUserPoolID)
	_, err := client.AdminCreateUser(createUserRequest)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed to processline %v user: %+v", lineNumber, userInfo), err)
	}

	// Disable user if it's 'enabled' column is not TRUE
	enabledCol, ok := colsMap["enabled"]
	if ok && len(line) > enabledCol && line[enabledCol] == "false" {
		userDisableRequest := userInfo.BuildDisableUserRequest(cfg.AWSCognitoUserPoolID)
		if _, err = client.AdminDisableUser(userDisableRequest); err != nil {
			log.Error(ctx, fmt.Sprintf("failed to disable user: %+v", userInfo), err)
		}
	}
}
