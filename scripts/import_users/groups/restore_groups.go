package groups

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ONSdigital/dp-identity-api/v2/scripts/import_users/config"

	"github.com/ONSdigital/dp-identity-api/v2/models"
	"github.com/ONSdigital/dp-identity-api/v2/scripts/utils"
	"github.com/ONSdigital/log.go/v2/log"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/pkg/errors"
)

func ImportGroupsFromS3(ctx context.Context, cfg *config.Config) error {
	log.Info(ctx, fmt.Sprintf("started restoring groups to cognito from s3 file: %v", cfg.GetS3GroupsFilePath()))

	s3FileReader := utils.S3Reader{}
	responseBody := s3FileReader.GetS3Reader(ctx, cfg.AWSProfile, cfg.S3Region, cfg.S3Bucket, cfg.GetS3GroupsFilePath())
	defer responseBody.Close()
	reader := csv.NewReader(responseBody)
	client := utils.GetCognitoClient(ctx, cfg.AWSProfile, cfg.S3Region)

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
		createGroup(ctx, count, line, client, cfg.AWSCognitoUserPoolID, colsMap)
		count++
	}
	log.Info(ctx, "Successfully processed all the groups in S3 file")
	return nil
}

func ImportGroupsMembersFromS3(ctx context.Context, cfg *config.Config) error {
	log.Info(ctx, fmt.Sprintf("started restoring group members to cognito from s3 file: %v", cfg.GetS3GroupsFilePath()))

	s3FileReader := utils.S3Reader{}
	responseBody := s3FileReader.GetS3Reader(ctx, cfg.AWSProfile, cfg.S3Region, cfg.S3Bucket, cfg.GetS3GroupUsersFilePath())
	defer responseBody.Close()
	reader := csv.NewReader(responseBody)
	client := utils.GetCognitoClient(ctx, cfg.AWSProfile, cfg.S3Region)

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
		adduserToGroup(ctx, client, line, count, cfg.AWSCognitoUserPoolID, colsMap)
		count++
	}
	log.Info(ctx, "Successfully processed all the group members in S3 file")
	return nil
}

func createGroup(ctx context.Context, lineNumber int, line []string, client *cognito.Client, userPoolID string, colsMap map[string]int) {
	createGroup := models.CreateUpdateGroup{}
	createGroup.ID = &line[colsMap["groupname"]]
	precedence, _ := strconv.Atoi(line[colsMap["precedence"]])
	precedence1 := int32(precedence)
	createGroup.Precedence = &precedence1
	createGroup.Name = &line[colsMap["description"]]
	input := createGroup.BuildCreateGroupInput(&userPoolID)

	_, err := client.CreateGroup(ctx, input)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed to process line: %d (incl. header) with name:%q group: %+v", lineNumber+1, *createGroup.Name, createGroup), err)
	}
}

func adduserToGroup(ctx context.Context, client *cognito.Client, line []string, lineNumber int, userPoolID string, colsMap map[string]int) {
	userID := line[colsMap["user_name"]]
	groups := strings.Split(line[colsMap["groups"]], ", ")
	for _, group := range groups {
		if group == "" {
			continue
		}
		_, err := client.AdminAddUserToGroup(ctx, &cognito.AdminAddUserToGroupInput{GroupName: &group, UserPoolId: &userPoolID, Username: &userID})
		if err != nil {
			log.Error(ctx, fmt.Sprintf("failed to process line: %d (incl. header) - user %v group: %v", lineNumber+1, userID, group), err)
		}
	}
}
