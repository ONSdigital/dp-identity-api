package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/ONSdigital/dp-identity-api/scripts/utils"
	"github.com/ONSdigital/log.go/v2/log"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"io"
	"strconv"
	"strings"
)

type Config struct {
	GroupsFilename       string `envconfig:"GROUPS_FILENAME" required:"true"`
	GroupUsersFilename   string `envconfig:"GROUPUSERS_FILENAME" required:"true"`
	S3Bucket             string `envconfig:"S3_BUCKET" required:"true"`
	S3BaseDir            string `envconfig:"S3_BASE_DIR" required:"true"`
	S3Region             string `envconfig:"S3_REGION" required:"true"`
	AWSCognitoUserPoolID string `envconfig:"USER_POOL_ID" required:"true"`
}

func (c Config) getS3GroupsFilePath() string {
	return fmt.Sprintf("%s%s", c.S3BaseDir, c.GroupsFilename)
}

func (c Config) getS3GroupUsersFilePath() string {
	return fmt.Sprintf("%s%s", c.S3BaseDir, c.GroupUsersFilename)
}

func readConfig() *Config {
	conf := &Config{}

	envconfig.Process("", conf)

	return conf
}

func main() {
	ctx := context.Background()
	conf := readConfig()

	fmt.Printf("Config: %+v", conf)
	importGroupsFromS3(ctx, conf)
	//importGroupsMembersFromS3(ctx, conf)
}

func importGroupsFromS3(ctx context.Context, config *Config) {
	log.Info(ctx, fmt.Sprintf("started restoring groups to cognito from s3 file: %v", config.getS3GroupsFilePath()))

	s3FileReader := utils.S3Reader{}
	responseBody := s3FileReader.GetS3Reader(ctx, config.S3Region, config.S3Bucket, config.getS3GroupsFilePath())
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
		createGroup(count, line, client, ctx)
		count += 1
	}
	log.Info(ctx, "Successfully processed all the groups in S3 file")
}

func importGroupsMembersFromS3(ctx context.Context, config *Config) {
	log.Info(ctx, fmt.Sprintf("started restoring group members to cognito from s3 file: %v", config.getS3GroupsFilePath()))

	s3FileReader := utils.S3Reader{}
	responseBody := s3FileReader.GetS3Reader(ctx, config.S3Region, config.S3Bucket, config.getS3GroupUsersFilePath())
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
		adduserToGroup(ctx, client, line, count, config.AWSCognitoUserPoolID)
		count += 1
	}
	log.Info(ctx, "Successfully processed all the group members in S3 file")
}

func createGroup(lineNumber int, line []string, client *cognito.CognitoIdentityProvider, ctx context.Context) {
	if len(line) <= 5 {
		log.Error(ctx, "", errors.New(fmt.Sprintf("line:%v - %+v is not in required format", lineNumber, line)))
		return
	}
	createGroup := models.CreateUpdateGroup{}
	createGroup.GroupName = &line[0]
	precedence, _ := strconv.Atoi(line[4])
	precedence1 := int64(precedence)
	createGroup.Precedence = &precedence1
	createGroup.Description = &line[2]
	input := createGroup.BuildCreateGroupInput(&line[1])

	_, err := client.CreateGroup(input)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed to process line:%v group: %+v", lineNumber, createGroup), err)
	}
}

func adduserToGroup(ctx context.Context, client *cognito.CognitoIdentityProvider, line []string, lineNumber int, userPoolId string) {
	if len(line) <= 2 {
		log.Error(ctx, "", errors.New(fmt.Sprintf("line:%v - %+v is not in required format", lineNumber, line)))
		return
	}
	userId := line[0]
	groups := strings.Split(line[1], ", ")
	for _, group := range groups {
		_, err := client.AdminAddUserToGroup(&cognito.AdminAddUserToGroupInput{GroupName: &group, UserPoolId: &userPoolId, Username: &userId})
		if err != nil {
			log.Error(ctx, fmt.Sprintf("failed to process line: %v -  user %v group: %v", lineNumber, userId, group), err)
		}
	}
}
