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
	"github.com/google/uuid"
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
	//skip header
	reader.Read()

	count := 1
	for {
		line, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		createUser(ctx, client, line, count, config)
		count += 1
	}
	log.Info(ctx, "Successfully processed all the users in S3 file")
	return nil
}

func createUser(ctx context.Context, client *cognitoidentityprovider.CognitoIdentityProvider, line []string, lineNumber int, config *config.Config) {
	if len(line) <= 21 {
		log.Error(ctx, "", errors.New(fmt.Sprintf("line:%v - %+v is not in required format", lineNumber, line)))
		return
	}
	userInfo := models.UserParams{
		Name:                line[1],
		GivenName:           line[2],
		FamilyName:          line[3],
		MiddleName:          line[4],
		Nickname:            line[5],
		PreferredUsername:   line[6],
		Profile:             line[7],
		Picture:             line[8],
		Website:             line[9],
		Email:               line[10],
		EmailVerified:       line[11],
		Gender:              line[12],
		Birthdate:           line[13],
		ZoneInfo:            line[14],
		Locale:              line[15],
		PhoneNumber:         line[16],
		PhoneNumberVerified: line[17],
		Address:             line[18],
		UpdatedAt:           line[19],
		CognitoMFAEnabled:   line[20],
		Username:            line[22],
	}
	userInfo.GeneratePassword(ctx)

	_, err := client.AdminCreateUser(userInfo.BuildCreateUserRequest(uuid.NewString(), config.AWSCognitoUserPoolID))
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed to processline %v user: %+v", lineNumber, userInfo), err)
	}
}
