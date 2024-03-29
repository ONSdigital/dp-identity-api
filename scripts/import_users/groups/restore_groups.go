package groups

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ONSdigital/dp-identity-api/scripts/import_users/config"

	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/ONSdigital/dp-identity-api/scripts/utils"
	"github.com/ONSdigital/log.go/v2/log"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/pkg/errors"
)

func ImportGroupsFromS3(ctx context.Context, config *config.Config) error {
	log.Info(ctx, fmt.Sprintf("started restoring groups to cognito from s3 file: %v", config.GetS3GroupsFilePath()))

	s3FileReader := utils.S3Reader{}
	responseBody := s3FileReader.GetS3Reader(ctx, config.AWSProfile, config.S3Region, config.S3Bucket, config.GetS3GroupsFilePath())
	defer responseBody.Close()
	reader := csv.NewReader(responseBody)
	client := utils.GetCognitoClient(config.AWSProfile, config.S3Region)

	// Extract column indexes from header line
	cols, err := reader.Read()
	if err != nil {
		return errors.New("unable to read header from groups backup file")
	}
	colsMap := make(map[string]int, len(cols))
	for i, col := range cols {
		colsMap[col] = i
	}
	for _, wanted := range []string{"groupname", "precedence", "description"} {
		if _, ok := colsMap[wanted]; !ok {
			return fmt.Errorf("column '%s' not found in groups backup file", wanted)
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
		createGroup(count, line, client, ctx, config.AWSCognitoUserPoolID, colsMap)
		count += 1
	}
	log.Info(ctx, "Successfully processed all the groups in S3 file")
	return nil
}

func ImportGroupsMembersFromS3(ctx context.Context, config *config.Config) error {
	log.Info(ctx, fmt.Sprintf("started restoring group members to cognito from s3 file: %v", config.GetS3GroupsFilePath()))

	s3FileReader := utils.S3Reader{}
	responseBody := s3FileReader.GetS3Reader(ctx, config.AWSProfile, config.S3Region, config.S3Bucket, config.GetS3GroupUsersFilePath())
	defer responseBody.Close()
	reader := csv.NewReader(responseBody)
	client := utils.GetCognitoClient(config.AWSProfile, config.S3Region)

	// Extract column indexes from header line
	cols, err := reader.Read()
	if err != nil {
		return errors.New("unable to read header from users_groups backup file")
	}
	colsMap := make(map[string]int, len(cols))
	for i, col := range cols {
		colsMap[col] = i
	}
	for _, wanted := range []string{"user_name", "groups"} {
		if _, ok := colsMap[wanted]; !ok {
			return fmt.Errorf("column '%s' not found in users_groups backup file", wanted)
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
		adduserToGroup(ctx, client, line, count, config.AWSCognitoUserPoolID, colsMap)
		count += 1
	}
	log.Info(ctx, "Successfully processed all the group members in S3 file")
	return nil
}

func createGroup(lineNumber int, line []string, client *cognito.CognitoIdentityProvider, ctx context.Context, userPoolID string, colsMap map[string]int) {
	createGroup := models.CreateUpdateGroup{}
	createGroup.ID = &line[colsMap["groupname"]]
	precedence, _ := strconv.Atoi(line[colsMap["precedence"]])
	precedence1 := int64(precedence)
	createGroup.Precedence = &precedence1
	createGroup.Name = &line[colsMap["description"]]
	input := createGroup.BuildCreateGroupInput(&userPoolID)

	_, err := client.CreateGroup(input)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed to process line: %d (incl. header) with name:%q group: %+v", lineNumber+1, *createGroup.Name, createGroup), err)
	}
}

func adduserToGroup(ctx context.Context, client *cognito.CognitoIdentityProvider, line []string, lineNumber int, userPoolId string, colsMap map[string]int) {
	userId := line[colsMap["user_name"]]
	groups := strings.Split(line[colsMap["groups"]], ", ")
	for _, group := range groups {
		if group == "" {
			continue
		}
		_, err := client.AdminAddUserToGroup(&cognito.AdminAddUserToGroupInput{GroupName: &group, UserPoolId: &userPoolId, Username: &userId})
		if err != nil {
			log.Error(ctx, fmt.Sprintf("failed to process line: %d (incl. header) - user %v group: %v", lineNumber+1, userId, group), err)
		}
	}
}
