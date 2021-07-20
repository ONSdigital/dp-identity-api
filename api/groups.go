package api

import (
	"context"
	"encoding/json"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
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

	validationErrs := group.ValidateAddUser(ctx, userId)
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
