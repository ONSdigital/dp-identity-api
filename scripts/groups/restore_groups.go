package main

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/ONSdigital/dp-identity-api/scripts/utils"
	"github.com/ONSdigital/log.go/v2/log"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"io"
	"os"
	"strconv"
	"strings"
)

type config struct {
	groupsFilename       string `envconfig:"groups_filename"`
	groupUsersFilename   string `envconfig:"groupusers_filename"`
	s3Bucket             string `envconfig:"s3_bucket"`
	s3BaseDir            string `envconfig:"s3_base_dir"`
	s3Region             string `envconfig:"s3_region"`
	awsCognitoUserPoolID string `envconfig:"user_pool_id"`
}

func (c config) getS3GroupsFilePath() string {
	return fmt.Sprintf("%s%s", c.s3BaseDir, c.groupsFilename)
}

func (c config) getS3GroupUsersFilePath() string {
	return fmt.Sprintf("%s%s", c.s3BaseDir, c.groupUsersFilename)
}

func readConfig() *config {
	conf := &config{}

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		switch pair[0] {
		case "groups_filename":
			conf.groupsFilename = pair[1]
		case "groupusers_filename":
			conf.groupUsersFilename = pair[1]
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

	if conf.groupsFilename == "" || conf.groupUsersFilename == "" {
		fmt.Println("Please set Environment Variables ")
		os.Exit(1)
	}

	return conf
}

func main() {
	ctx := context.Background()
	conf := readConfig()

	fmt.Printf("config: %+v", conf)
	importGroupsFromS3(ctx, conf)
	importGroupsMembersFromS3(ctx, conf)
}

func importGroupsFromS3(ctx context.Context, config *config) {
	log.Info(ctx, fmt.Sprintf("started restoring groups to cognito from s3 file: %v", config.getS3GroupsFilePath()))

	s3FileReader := utils.S3Reader{}
	reader := s3FileReader.GetS3Reader(ctx, config.s3Region, config.s3Bucket, config.getS3GroupsFilePath())
	defer s3FileReader.Close()
	client := utils.GetCognitoClient(config.s3Region)

	for {
		line, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		createGroup(line, client, ctx)
	}
	log.Info(ctx, "Successfully processed all the groups in S3 file")
}

func importGroupsMembersFromS3(ctx context.Context, config *config) {
	log.Info(ctx, fmt.Sprintf("started restoring group members to cognito from s3 file: %v", config.getS3GroupsFilePath()))

	s3FileReader := utils.S3Reader{}
	reader := s3FileReader.GetS3Reader(ctx, config.s3Region, config.s3Bucket, config.getS3GroupUsersFilePath())
	defer s3FileReader.Close()
	client := utils.GetCognitoClient(config.s3Region)

	for {
		line, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		adduserToGroup(ctx, client, line, config.awsCognitoUserPoolID)
		fmt.Println(line)
	}
	log.Info(ctx, "Successfully processed all the group members in S3 file")
}

func createGroup(line []string, client *cognito.CognitoIdentityProvider, ctx context.Context) {
	createGroup := models.CreateUpdateGroup{}
	createGroup.GroupName = &line[0]
	precedence, _ := strconv.Atoi(line[4])
	precedence1 := int64(precedence)
	createGroup.Precedence = &precedence1
	createGroup.Description = &line[2]
	input := createGroup.BuildCreateGroupInput(&line[1])

	_, err := client.CreateGroup(input)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed to create group: %+v", createGroup), err)
	}
}

func adduserToGroup(ctx context.Context, client *cognito.CognitoIdentityProvider, line []string, userPoolId string) {

	userId := line[0]
	groups := strings.Split(line[1], ", ")
	for _, group := range groups {
		_, err := client.AdminAddUserToGroup(&cognito.AdminAddUserToGroupInput{GroupName: &group, UserPoolId: &userPoolId, Username: &userId})
		if err != nil {
			log.Error(ctx, fmt.Sprintf("failed to add user %v to group: %v", userId, group), err)
		}
	}
}
