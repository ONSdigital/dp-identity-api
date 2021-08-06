package api

import (
	"context"
	"encoding/csv"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/ONSdigital/dp-identity-api/utilities"
	"github.com/ONSdigital/log.go/log"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

const (
	usersGroupsFileName = "users_groups_export_"
	backupFilePath      = "/Users/mk-dv-m0047/workspace/test/"
	backupFileName      = "groups_export_2021-08-06_13_27_11.csv"
)

type usersGroups struct {
	groupName string
	users []*cognitoidentityprovider.UserType
}

func (api *API) ListUsersGroupsHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	var (
		awsErr error
		listUsersGroupsResp, result *cognitoidentityprovider.ListUsersInGroupOutput
		listUsersGroupsInput = models.BuildListUsersGroupsRequest(api.UserPoolId, "", "")
		usersGroupsError *models.ErrorResponse
		backOffSchedule = []time.Duration{
			1 * time.Second,
			3 * time.Second,
			10 * time.Second,
		}
		userGroupsList = []usersGroups{}
		debug = true
	)

	// read data from backup - only a POC so breaking any modeling not too problematic
	allUsersGroupsDataBackup, err := utilities.LoadDataFromCSV(ctx, backupFilePath+backupFileName)
	if err != nil {
		log.Event(ctx, "fatal runtime error", log.Error(err), log.FATAL)
		os.Exit(1)
	}

	// traverse all users in groups data
	for _, groupData := range allUsersGroupsDataBackup {
		// set group name for request
		listUsersGroupsInput.GroupName = &groupData[0]
		listUsersGroupsInput.NextToken = nil
		listUsersGroupsResp, awsErr = api.generateListUsersGroupsRequest(listUsersGroupsInput)
		// continue oif no users in group
		if len(listUsersGroupsResp.Users) == 0 {
			continue
		}
		if awsErr != nil {
			err := models.NewCognitoError(ctx, awsErr, "Cognito ListUsersInGroup request from list users in groups endpoint")
			usersGroupsError = models.NewErrorResponse(http.StatusInternalServerError, nil, err)
		} else {
			if listUsersGroupsResp.NextToken != nil {
				listUsersGroupsInput.NextToken = listUsersGroupsResp.NextToken
				// set `loadingInProgress` to control requesting new list data
				loadingInProgress := true
				for loadingInProgress {
					for _, backoff := range backOffSchedule {
						result, awsErr = api.generateListUsersGroupsRequest(listUsersGroupsInput)
						if awsErr == nil {
							listUsersGroupsResp.Users = append(listUsersGroupsResp.Users, result.Users...)
							if result.NextToken != nil {
								listUsersGroupsInput.NextToken = result.NextToken
								break
							} else {
								loadingInProgress = false
								break
							}
						} else {
							err := models.NewCognitoError(ctx, awsErr, "Cognito ListUsers request from signout all users from group endpoint")
							if err.Code != models.TooManyRequestsError {
								usersGroupsError = models.NewErrorResponse(http.StatusInternalServerError, nil, err)
								loadingInProgress = false
								break
							}
						}
						time.Sleep(backoff)
					}
				}
			}
		}
		// debug
		if debug {
			println("user count for group: "+groupData[0]+" is : "+strconv.Itoa(len(listUsersGroupsResp.Users)))
		}
		// add to userGroupsList
		userGroupsList = append(
			userGroupsList,
			usersGroups{
				groupName: groupData[0],
				users: listUsersGroupsResp.Users,
			},
		)
	}
	if usersGroupsError != nil {
		return nil, usersGroupsError
	} else {
		t := time.Now()
		file, _ := os.Create(filePath+usersGroupsFileName+t.Format(dateLayout)+".csv")
		defer file.Close()
	
		writer := csv.NewWriter(file)
		defer writer.Flush()

		for _, userGroupsData := range userGroupsList {
			// build a string using user names
			var userNamesList string
			for _, user := range userGroupsData.users {
				userNamesList += ","+*user.Username
			}
			userGroupsData := userGroupsData.groupName+userNamesList
			_ = writer.Write(strings.Split(userGroupsData, ","))
		}
	}
	return models.NewSuccessResponse(nil, http.StatusOK, nil), nil
}

// generateListUsersGroupsRequest - local routine to generate users in groups request
func (api *API) generateListUsersGroupsRequest(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
	result, err := api.CognitoClient.ListUsersInGroup(input)
	return result, err
}
