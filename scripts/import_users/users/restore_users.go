package users

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"io"

	"github.com/ONSdigital/dp-identity-api/v2/models"
	"github.com/ONSdigital/dp-identity-api/v2/scripts/import_users/config"
	"github.com/ONSdigital/dp-identity-api/v2/scripts/utils"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/pkg/errors"
)

func ImportUsersFromS3(ctx context.Context, cfg *config.Config) error {
	log.Info(ctx, fmt.Sprintf("started restoring users to cognito from s3 file: %v", cfg.GetS3UsersFilePath()))

	s3FileReader := utils.S3Reader{}
	responseBody := s3FileReader.GetS3Reader(ctx, cfg.AWSProfile, cfg.S3Region, cfg.S3Bucket, cfg.GetS3UsersFilePath())
	defer responseBody.Close()
	reader := csv.NewReader(responseBody)
	client := utils.GetCognitoClient(ctx, cfg.AWSProfile, cfg.S3Region)

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
			return fmt.Errorf("column '%s' not found in users backup file", wanted)
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
		count++
	}
	log.Info(ctx, "Successfully processed all the users in S3 file")
	return nil
}

func createUser(ctx context.Context, client *cognitoidentityprovider.Client, line []string, lineNumber int, cfg *config.Config, colsMap map[string]int) {
	userInfo := models.UserParams{
		Forename: line[colsMap["given_name"]],
		Lastname: line[colsMap["family_name"]],
		Email:    line[colsMap["email"]],
		ID:       line[colsMap["cognito:username"]],
	}

	var err error
	_, err = userInfo.GeneratePassword(ctx)
	if err != nil {
		log.Error(ctx, "failed to generate password", err)
	}

	createUserRequest := userInfo.BuildCreateUserRequest(userInfo.ID, cfg.AWSCognitoUserPoolID)
	_, errCreateUser := client.AdminCreateUser(ctx, createUserRequest)
	if errCreateUser != nil {
		log.Error(ctx, fmt.Sprintf("failed to process line %v user: %+v", lineNumber, userInfo), err)
	}

	// Disable user if it's 'enabled' column is not TRUE
	enabledCol, ok := colsMap["enabled"]
	if ok && len(line) > enabledCol && line[enabledCol] == "false" {
		userDisableRequest := userInfo.BuildDisableUserRequest(cfg.AWSCognitoUserPoolID)
		if _, err = client.AdminDisableUser(ctx, userDisableRequest); err != nil {
			log.Error(ctx, fmt.Sprintf("failed to disable user: %+v", userInfo), err)
		}
	}
}
