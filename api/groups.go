package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gorilla/mux"
)

//AddUserToGroupHandler adds a user to the specified group
func (api *API) AddUserToGroupHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	vars := mux.Vars(req)
	group := models.Group{Name: vars["id"]}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, handleBodyReadError(ctx, err)
	}

	var bodyJson map[string]string
	err = json.Unmarshal(body, &bodyJson)
	if err != nil {
		return nil, handleBodyUnmarshalError(ctx, err)
	}
	userId := bodyJson["user_id"]

	validationErrs := group.ValidateAddRemoveUser(ctx, userId)
	if len(validationErrs) != 0 {
		return nil, models.NewErrorResponse(http.StatusBadRequest, nil, validationErrs...)
	}

	userAddToGroupInput := group.BuildAddUserToGroupRequest(api.UserPoolId, userId)
	_, err = api.CognitoClient.AdminAddUserToGroup(userAddToGroupInput)
	if err != nil {
		cognitoErr := models.NewCognitoError(ctx, err, "Cognito AddUserToGroup request from add user to group endpoint")
		if cognitoErr.Code == models.UserNotFoundError || cognitoErr.Code == models.NotFoundError {
			return nil, models.NewErrorResponse(http.StatusBadRequest, nil, cognitoErr)
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, cognitoErr)
	}

	groupGetRequest := group.BuildGetGroupRequest(api.UserPoolId)
	groupGetResponse, err := api.CognitoClient.GetGroup(groupGetRequest)
	if err != nil {
		cognitoErr := models.NewCognitoError(ctx, err, "Cognito GetGroup request from add user to group endpoint")
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, cognitoErr)
	}

	group.MapCognitoDetails(groupGetResponse.Group)

	groupMembersRequest := group.BuildListUsersInGroupRequest(api.UserPoolId)
	groupMembersResponse, err := api.CognitoClient.ListUsersInGroup(groupMembersRequest)
	if err != nil {
		cognitoErr := models.NewCognitoError(ctx, err, "Cognito ListUsersInGroup request from add user to group endpoint")
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, cognitoErr)
	}

	group.MapMembers(&groupMembersResponse.Users)

	jsonResponse, responseErr := group.BuildSuccessfulJsonResponse(ctx)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	return models.NewSuccessResponse(jsonResponse, http.StatusOK, nil), nil
}

//ListUsersInGroupHandler list the users in the specified group
func (api *API) ListUsersInGroupHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {

	vars := mux.Vars(req)
	group := models.Group{Name: vars["id"]}

	listOfUsersInput := []*cognitoidentityprovider.UserType{}

	listUsers, err := api.getUsersInAGroup(listOfUsersInput, group)
	if err != nil {
		cognitoErr := models.NewCognitoError(ctx, err, "Cognito ListUsersInGroup request from list users in group endpoint")
		if cognitoErr.Code == models.NotFoundError {
			return nil, models.NewErrorResponse(http.StatusBadRequest, nil, cognitoErr)
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, cognitoErr)
	}

	listOfUsers := models.UsersList{}
	listOfUsers.MapCognitoUsers(&listUsers)

	jsonResponse, responseErr := listOfUsers.BuildSuccessfulJsonResponse(ctx)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	return models.NewSuccessResponse(jsonResponse, http.StatusOK, nil), nil
}

func (api *API) getUsersInAGroup(listOfUsers []*cognitoidentityprovider.UserType, group models.Group) ([]*cognitoidentityprovider.UserType, error) {
	firstTimeCheck := false
	var nextToken string
	for {
		if firstTimeCheck && nextToken == "" {
			break
		}
		firstTimeCheck = true

		groupMembersRequest := group.BuildListUsersInGroupRequestWithNextToken(api.UserPoolId, nextToken)
		groupMembersResponse, err := api.CognitoClient.ListUsersInGroup(groupMembersRequest)
		if err != nil {
			return nil, err
		}

		listOfUsers = append(listOfUsers, groupMembersResponse.Users...)
		nextToken = ""
		if groupMembersResponse.NextToken != nil {
			nextToken = *groupMembersResponse.NextToken
		}
	}
	return listOfUsers, nil
}

//RemoveUserFromGroupHandler adds a user to the specified group
func (api *API) RemoveUserFromGroupHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	vars := mux.Vars(req)
	group := models.Group{Name: vars["id"]}

	userId := vars["user_id"]

	validationErrs := group.ValidateAddRemoveUser(ctx, userId)
	if len(validationErrs) != 0 {
		return nil, models.NewErrorResponse(http.StatusBadRequest, nil, validationErrs...)
	}

	userRemoveFromGroupInput := group.BuildRemoveUserFromGroupRequest(api.UserPoolId, userId)
	_, err := api.CognitoClient.AdminRemoveUserFromGroup(userRemoveFromGroupInput)
	if err != nil {
		cognitoErr := models.NewCognitoError(ctx, err, "Cognito RemoveUserFromGroup request from remove user from group endpoint")
		if cognitoErr.Code == models.UserNotFoundError || cognitoErr.Code == models.NotFoundError {
			return nil, models.NewErrorResponse(http.StatusBadRequest, nil, cognitoErr)
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, cognitoErr)
	}

	groupGetRequest := group.BuildGetGroupRequest(api.UserPoolId)
	groupGetResponse, err := api.CognitoClient.GetGroup(groupGetRequest)
	if err != nil {
		cognitoErr := models.NewCognitoError(ctx, err, "Cognito GetGroup request from remove user from group endpoint")
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, cognitoErr)
	}

	group.MapCognitoDetails(groupGetResponse.Group)

	groupMembersRequest := group.BuildListUsersInGroupRequest(api.UserPoolId)
	groupMembersResponse, err := api.CognitoClient.ListUsersInGroup(groupMembersRequest)
	if err != nil {
		cognitoErr := models.NewCognitoError(ctx, err, "Cognito ListUsersInGroup request from remove user from group endpoint")
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, cognitoErr)
	}

	group.MapMembers(&groupMembersResponse.Users)

	jsonResponse, responseErr := group.BuildSuccessfulJsonResponse(ctx)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	return models.NewSuccessResponse(jsonResponse, http.StatusOK, nil), nil
}

//List Groups for user pagination allows first call and then any other call if nextToken is not ""
func (api *API) getListGroups() (*cognitoidentityprovider.ListGroupsOutput, error) {
	firstTimeCheck := false
	var nextToken string
	group := models.ListUserGroupType{}

	listOfGroups := cognitoidentityprovider.ListGroupsOutput{}
	for {
		if firstTimeCheck && nextToken == "" {
			break
		}
		firstTimeCheck = true

		listGroupsRequest := group.BuildListGroupsRequest(api.UserPoolId, nextToken)
		listGroupsResponse, err := api.CognitoClient.ListGroups(listGroupsRequest)
		if err != nil {
			return nil, err
		}

		listOfGroups.Groups = append(listOfGroups.Groups, listGroupsResponse.Groups...)
		nextToken = ""
		if listGroupsResponse.NextToken != nil {
			nextToken = *listGroupsResponse.NextToken
		}
	}
	return &listOfGroups, nil
}

//ListGroupsHandler lists the users in the user pool
func (api *API) ListGroupsHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {

	// vars := mux.Vars(req)
	finalGroupsResponse := models.ListUserGroups{}

	listOfGroups, err := api.getListGroups()
	if err != nil {
		cognitoErr := models.NewCognitoError(ctx, err, "Cognito ListofUserGroups request from list user groups endpoint")
		if cognitoErr.Code == models.NotFoundError {
			return nil, models.NewErrorResponse(http.StatusNotFound, nil, cognitoErr)
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, cognitoErr)
	}

	jsonResponse, responseErr := finalGroupsResponse.BuildListGroupsSuccessfulJsonResponse(ctx, listOfGroups)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}
	return models.NewSuccessResponse(jsonResponse, http.StatusOK, nil), nil

}
