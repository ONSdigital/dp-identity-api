package groups

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/ONSdigital/dp-identity-api/scripts/import_users/config"
	"io"
	"strconv"
	"strings"

	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/ONSdigital/dp-identity-api/scripts/utils"
	"github.com/ONSdigital/log.go/v2/log"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/pkg/errors"
)

func ImportGroupsFromS3(ctx context.Context, config *config.Config) error {
	log.Info(ctx, fmt.Sprintf("started restoring groups to cognito from s3 file: %v", config.GetS3GroupsFilePath()))

	s3FileReader := utils.S3Reader{}
	responseBody := s3FileReader.GetS3Reader(ctx, config.S3Region, config.S3Bucket, config.GetS3GroupsFilePath())
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
	return nil
}

func ImportGroupsMembersFromS3(ctx context.Context, config *config.Config) error {
	log.Info(ctx, fmt.Sprintf("started restoring group members to cognito from s3 file: %v", config.GetS3GroupsFilePath()))

	s3FileReader := utils.S3Reader{}
	responseBody := s3FileReader.GetS3Reader(ctx, config.S3Region, config.S3Bucket, config.GetS3GroupUsersFilePath())
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
	return nil
}

func createGroup(lineNumber int, line []string, client *cognito.CognitoIdentityProvider, ctx context.Context) {
	if len(line) <= 5 {
		log.Error(ctx, "", errors.New(fmt.Sprintf("line:%v - %+v is not in required format", lineNumber, line)))
		return
	}
	createGroup := models.CreateUpdateGroup{}
	createGroup.ID = &line[0]
	precedence, _ := strconv.Atoi(line[4])
	precedence1 := int64(precedence)
	createGroup.Precedence = &precedence1
	createGroup.Name = &line[2]
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
