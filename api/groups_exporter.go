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
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

const (
	filePath     = "/Users/mk-dv-m0047/workspace/test/"
	fileName     = "groups_export_"
	dateLayout   = "2006-01-02_15_04_05"
	headerLayout = "GroupName,Description,CreationDate,LastModifiedDate,Precedence,RoleARN,UserPoolId"
)

func (api *API) ListGroupsHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	var (
		awsErr error
		listGroupsResp, result *cognitoidentityprovider.ListGroupsOutput
		listGroupsInput = models.BuildListGroupsRequest( api.UserPoolId, "")
		usersGroupsError *models.ErrorResponse
		backOffSchedule = []time.Duration{
			1 * time.Second,
			3 * time.Second,
			10 * time.Second,
		}
	)
	listGroupsResp, awsErr = api.generateListGroupsRequest(listGroupsInput)
	if awsErr != nil {
		err := models.NewCognitoError(ctx, awsErr, "Cognito ListGroups request from list groups endpoint")
		usersGroupsError = models.NewErrorResponse(http.StatusInternalServerError, nil, err)
	} else {
		if listGroupsResp.NextToken != nil {
			listGroupsInput.NextToken = listGroupsResp.NextToken
			// set `loadingInProgress` to control requesting new list data
			loadingInProgress := true
			for loadingInProgress {
				for _, backoff := range backOffSchedule {
					result, awsErr = api.generateListGroupsRequest(listGroupsInput)
					if awsErr == nil {
						listGroupsResp.Groups = append(listGroupsResp.Groups, result.Groups...)
						if result.NextToken != nil {
							listGroupsInput.NextToken = result.NextToken
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
	if usersGroupsError != nil {
		return nil, usersGroupsError
	} else {
		t := time.Now()
		file, _ := os.Create(filePath+fileName+t.Format(dateLayout)+".csv")
		defer file.Close()
	
		writer := csv.NewWriter(file)
		defer writer.Flush()
		
		// write headers
		writer.Write(strings.Split(headerLayout, ","))

		for _, groupData := range listGroupsResp.Groups {
			creationDate := *groupData.CreationDate
			lastModifiedDate := *groupData.LastModifiedDate
			groupDataString := *groupData.GroupName+
				","+*groupData.Description+
				","+creationDate.Format(dateLayout)+
				","+lastModifiedDate.Format(dateLayout)+
				","+strconv.Itoa(int(*groupData.Precedence))+
				","+*groupData.UserPoolId
			_ = writer.Write(strings.Split(groupDataString, ","))
		}
		return models.NewSuccessResponse(nil, http.StatusAccepted, nil), nil
	}
}

// generateListUsersRequest - local routine to generate a list users request
func (api *API) generateListGroupsRequest(input *cognitoidentityprovider.ListGroupsInput) (*cognitoidentityprovider.ListGroupsOutput, error) {
	return api.CognitoClient.ListGroups(input)
}