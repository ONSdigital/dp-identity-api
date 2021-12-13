package main

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/ONSdigital/dp-identity-api/scripts/utils"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/google/uuid"
	"io"
	"os"
	"strconv"
	"strings"
)

type userRestoreConfig struct {
	userFileName                  string
	s3Bucket, s3BaseDir, s3Region string
	awsCognitoUserPoolID          string
}

func (c userRestoreConfig) getS3UsersFilePath() string {
	return fmt.Sprintf("%s%s", c.s3BaseDir, c.userFileName)
}

func getConfig() *userRestoreConfig {
	conf := &userRestoreConfig{}
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		switch pair[0] {
		case "filename":
			conf.userFileName = pair[1]
		case "s3_bucket":
			conf.s3Bucket = pair[1]
		case "s3_base_dir":
			conf.s3BaseDir = pair[1]
		case "s3_region":
			conf.s3Region = pair[1]
		case "user_pool_id":
			conf.awsCognitoUserPoolID = pair[1]
		}
	}
	if conf.userFileName == "" {
		fmt.Println("Please set Environment Variables ")
		os.Exit(1)
	}

	return conf
}

func main() {
	ctx := context.Background()
	config := getConfig()

	importUsersFromS3(ctx, config)
}

func importUsersFromS3(ctx context.Context, config *userRestoreConfig) {
	log.Info(ctx, fmt.Sprintf("started restoring users to cognito from s3 file: %v", config.getS3UsersFilePath()))

	s3FileReader := utils.S3Reader{}
	reader := s3FileReader.GetS3Reader(ctx, config.s3Region, config.s3Bucket, config.getS3UsersFilePath())
	defer s3FileReader.Close()

	client := utils.GetCognitoClient(config.s3Region)

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

func createUser(ctx context.Context, client *cognitoidentityprovider.CognitoIdentityProvider, line []string, config *userRestoreConfig) {
	isActive, _ := strconv.ParseBool(line[12])
	userInfo := models.UserParams{Forename: line[2], Lastname: line[3], Email: line[10], Active: isActive}
	userInfo.GeneratePassword(ctx)

	_, err := client.AdminCreateUser(userInfo.BuildCreateUserRequest(uuid.NewString(), config.awsCognitoUserPoolID))
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed to create user: %+v", userInfo), err)
	}
}
