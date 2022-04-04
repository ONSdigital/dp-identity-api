package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"

	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gorilla/mux"
)

const (
	GroupsCreatePermission string = "users:create"
	GroupsReadPermission          = "users:read"
	GroupsEditPermission          = "users:update"
	GroupsDeletePermission        = "users:update"
)

//CreateGroupHandler creates a new group
func (api *API) CreateGroupHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, handleBodyReadError(ctx, err)
	}

	createGroup := models.CreateUpdateGroup{}
	err = json.Unmarshal(body, &createGroup)
	if err != nil {
		return nil, handleBodyUnmarshalError(ctx, err)
	}

	// no groupname in body, set UUID
	if createGroup.ID == nil {
		uuid := uuid.NewString()
		createGroup.ID = &uuid
	}

	createGroup.GroupsList, err = api.GetListGroups()
	if err != nil {
		cognitoErr := models.NewCognitoError(ctx, err, "Cognito ListGroups request from update group endpoint")
		if cognitoErr.Code == models.NotFoundError {
			return nil, models.NewErrorResponse(http.StatusNotFound, nil, cognitoErr)
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, cognitoErr)
	}

	validationErrs := createGroup.ValidateCreateUpdateGroupRequest(ctx)
	if len(validationErrs) != 0 {
		return nil, models.NewErrorResponse(http.StatusBadRequest, nil, validationErrs...)
	}

	// build create group input
	input := createGroup.BuildCreateGroupInput(&api.UserPoolId)
	_, err = api.CognitoClient.CreateGroup(input)
	if err != nil {
		cognitoErr := models.NewCognitoError(ctx, err, "Cognito CreateGroup request from create a new group endpoint")
		if cognitoErr.Code == models.GroupExistsError {
			return nil, models.NewErrorResponse(http.StatusBadRequest, nil, cognitoErr)
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, cognitoErr)
	}

	jsonResponse, responseErr := createGroup.BuildSuccessfulJsonResponse(ctx)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	return createGroup.NewSuccessResponse(jsonResponse, http.StatusCreated, nil), nil
}

//UpdateGroupHandler update group details for a given group by id (GroupName)
func (api *API) UpdateGroupHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	vars := mux.Vars(req)

	id := vars["id"]
	updateGroup := models.CreateUpdateGroup{
		ID: &id,
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, handleBodyReadError(ctx, err)
	}

	err = json.Unmarshal(body, &updateGroup)
	if err != nil {
		return nil, handleBodyUnmarshalError(ctx, err)
	}

	validationErrs := updateGroup.ValidateCreateUpdateGroupRequest(ctx)
	if len(validationErrs) != 0 {
		return nil, models.NewErrorResponse(http.StatusBadRequest, nil, validationErrs...)
	}

	input := updateGroup.BuildUpdateGroupInput(api.UserPoolId)
	_, err = api.CognitoClient.UpdateGroup(input)
	if err != nil {
		cognitoErr := models.NewCognitoError(ctx, err, "Cognito UpdateGroup request from update a group endpoint")
		if cognitoErr.Code == models.NotFoundError {
			return nil, models.NewErrorResponse(http.StatusNotFound, nil, cognitoErr)
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, cognitoErr)
	}

	jsonResponse, responseErr := updateGroup.BuildSuccessfulJsonResponse(ctx)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	return updateGroup.NewSuccessResponse(jsonResponse, http.StatusOK, nil), nil
}

//AddUserToGroupHandler adds a user to the specified group
func (api *API) AddUserToGroupHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	vars := mux.Vars(req)
	group := models.Group{ID: vars["id"]}

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

	jsonResponse, responseErr := group.BuildSuccessfulJsonResponse(ctx)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	return models.NewSuccessResponse(jsonResponse, http.StatusOK, nil), nil
}

//ListUsersInGroupHandler list the users in the specified group
func (api *API) ListUsersInGroupHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {

	vars := mux.Vars(req)
	group := models.Group{ID: vars["id"]}

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
	group := models.Group{ID: vars["id"]}

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

	jsonResponse, responseErr := group.BuildSuccessfulJsonResponse(ctx)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	return models.NewSuccessResponse(jsonResponse, http.StatusOK, nil), nil
}

//List Groups pagination allows first call and then any other call if nextToken is not ""
func (api *API) GetListGroups() (*cognitoidentityprovider.ListGroupsOutput, error) {
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
	finalGroupsResponse := models.ListUserGroups{}

	listOfGroups, err := api.GetListGroups()
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

//GetGroupHandler gets group details for given groups
func (api *API) GetGroupHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {

	vars := mux.Vars(req)
	group := models.Group{ID: vars["id"]}
	groupGetRequest := group.BuildGetGroupRequest(api.UserPoolId)
	groupGetResponse, err := api.CognitoClient.GetGroup(groupGetRequest)
	if err != nil {

		cognitoErr := models.NewCognitoError(ctx, err, "Cognito GetGroup request from Get group endpoint")
		if cognitoErr.Code == models.NotFoundError {
			return nil, models.NewErrorResponse(http.StatusNotFound, nil, cognitoErr)
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, cognitoErr)
	}

	group.MapCognitoDetails(groupGetResponse.Group)

	jsonResponse, responseErr := group.BuildSuccessfulJsonResponse(ctx)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	return models.NewSuccessResponse(jsonResponse, http.StatusOK, nil), nil
}

//DeleteGroupHandler deletes the group for the given group id
func (api *API) DeleteGroupHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {

	vars := mux.Vars(req)
	group := models.Group{ID: vars["id"]}
	groupDeleteRequest := group.BuildDeleteGroupRequest(api.UserPoolId)
	_, err := api.CognitoClient.DeleteGroup(groupDeleteRequest)
	if err != nil {

		cognitoErr := models.NewCognitoError(ctx, err, "Cognito DeleteGroup request from Delete group endpoint")
		if cognitoErr.Code == models.NotFoundError {
			return nil, models.NewErrorResponse(http.StatusNotFound, nil, cognitoErr)
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, cognitoErr)
	}

	return models.NewSuccessResponse(nil, http.StatusNoContent, nil), nil
}

//SetGroupUsersHandler adds a user to the specified group
func (api *API) SetGroupUsersHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	// groupUsers, groupUsersErrors := api.ListUserGroupsHandler(ctx, w, req)
	// fmt.Println(groupUsers)
	// fmt.Println(groupUsersErrors)

	return nil, nil
}
