package api

import (
	"context"
	"encoding/json"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
)

//CreateUserHandler creates a new user and returns a http handler interface
func (api *API) CreateUserHandler(w http.ResponseWriter, req *http.Request, ctx context.Context) (*models.SuccessResponse, *models.ErrorResponse) {
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, handleBodyReadError(ctx, err)
	}

	user := models.UserParams{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, handleBodyUnmarshalError(ctx, err)
	}

	err = user.GeneratePassword(ctx)
	if err != nil {
		return nil, models.NewErrorResponse([]error{err}, http.StatusInternalServerError)
	}

	validationErrs := user.ValidateRegistration(ctx)

	listUserInput := user.BuildListUserRequest("email = \""+user.Email+"\"", "email", int64(1), &api.UserPoolId)
	listUserResp, err := api.CognitoClient.ListUsers(listUserInput)
	if err != nil {
		return nil, models.NewErrorResponse([]error{models.NewCognitoError(ctx, err, "Cognito ListUsers request from create users endpoint")}, http.StatusInternalServerError)
	}
	duplicateEmailErr := user.CheckForDuplicateEmail(ctx, listUserResp)
	if duplicateEmailErr != nil {
		validationErrs = append(validationErrs, duplicateEmailErr)
	}

	if len(validationErrs) != 0 {
		return nil, models.NewErrorResponse(validationErrs, http.StatusBadRequest)
	}

	createUserRequest := user.BuildCreateUserRequest(uuid.NewString(), api.UserPoolId)

	resultUser, err := api.CognitoClient.AdminCreateUser(createUserRequest)
	if err != nil {
		responseErr := models.NewCognitoError(ctx, err, "AdminCreateUser request from create user endpoint")
		if responseErr.Code == models.InternalError {
			return nil, models.NewErrorResponse([]error{responseErr}, http.StatusInternalServerError)
		} else {
			return nil, models.NewErrorResponse([]error{responseErr}, http.StatusBadRequest)
		}
	}

	jsonResponse, responseErr := user.BuildSuccessfulJsonResponse(ctx, resultUser)
	if responseErr != nil {
		return nil, models.NewErrorResponse([]error{responseErr}, http.StatusInternalServerError)
	}

	return models.NewSuccessResponse(jsonResponse, http.StatusCreated), nil
}
