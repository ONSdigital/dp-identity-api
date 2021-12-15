package main

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/ONSdigital/dp-identity-api/scripts/utils"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"io"
	"strconv"
)

type Config struct {
	UserFileName         string `envconfig:"FILENAME" required:"true"`
	S3Bucket             string `envconfig:"S3_BUCKET" required:"true"`
	S3BaseDir            string `envconfig:"S3_BASE_DIR" required:"true"`
	S3Region             string `envconfig:"S3_REGION" required:"true"`
	AWSCognitoUserPoolID string `envconfig:"USER_POOL_ID" required:"true"`
}

func (c Config) getS3UsersFilePath() string {
	return fmt.Sprintf("%s%s", c.S3BaseDir, c.UserFileName)
}

func getConfig() *Config {
	conf := &Config{}
	envconfig.Process("", conf)
	return conf
}

func main() {
	ctx := context.Background()
	config := getConfig()

	importUsersFromS3(ctx, config)
}

func importUsersFromS3(ctx context.Context, config *Config) {
	log.Info(ctx, fmt.Sprintf("started restoring users to cognito from s3 file: %v", config.getS3UsersFilePath()))

	s3FileReader := utils.S3Reader{}
	reader := s3FileReader.GetS3Reader(ctx, config.S3Region, config.S3Bucket, config.getS3UsersFilePath())
	defer s3FileReader.Close()

	client := utils.GetCognitoClient(config.S3Region)

	for {
		line, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		createUser(ctx, client, line, config)
	}
	log.Info(ctx, "Successfully processed all the users in S3 file")
}

func createUser(ctx context.Context, client *cognitoidentityprovider.CognitoIdentityProvider, line []string, config *Config) {
	isActive, _ := strconv.ParseBool(line[12])
	userInfo := models.UserParams{Forename: line[2], Lastname: line[3], Email: line[10], Active: isActive}
	userInfo.GeneratePassword(ctx)

	_, err := client.AdminCreateUser(userInfo.BuildCreateUserRequest(uuid.NewString(), config.AWSCognitoUserPoolID))
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed to create user: %+v", userInfo), err)
	}
}
