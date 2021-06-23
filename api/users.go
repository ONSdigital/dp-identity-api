package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/google/uuid"
)

//CreateUserHandler creates a new user and returns a http handler interface
func (api *API) CreateUserHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
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

	tempPassword, err := user.GeneratePassword(ctx)
	if err != nil {
		return nil, models.NewErrorResponse([]error{err}, http.StatusInternalServerError, nil)
	}
	user.Password = *tempPassword

	validationErrs := user.ValidateRegistration(ctx)

	listUserInput := models.UsersList{}.BuildListUserRequest("email = \""+user.Email+"\"", "email", int64(1), &api.UserPoolId)
	listUserResp, err := api.CognitoClient.ListUsers(listUserInput)
	if err != nil {
		return nil, models.NewErrorResponse([]error{models.NewCognitoError(ctx, err, "Cognito ListUsers request from create users endpoint")}, http.StatusInternalServerError, nil)
	}
	duplicateEmailErr := user.CheckForDuplicateEmail(ctx, listUserResp)
	if duplicateEmailErr != nil {
		validationErrs = append(validationErrs, duplicateEmailErr)
	}

	if len(validationErrs) != 0 {
		return nil, models.NewErrorResponse(validationErrs, http.StatusBadRequest, nil)
	}

	createUserRequest := user.BuildCreateUserRequest(uuid.NewString(), api.UserPoolId)

	resultUser, err := api.CognitoClient.AdminCreateUser(createUserRequest)
	if err != nil {
		responseErr := models.NewCognitoError(ctx, err, "AdminCreateUser request from create user endpoint")
		if responseErr.Code == models.InternalError {
			return nil, models.NewErrorResponse([]error{responseErr}, http.StatusInternalServerError, nil)
		} else {
			return nil, models.NewErrorResponse([]error{responseErr}, http.StatusBadRequest, nil)
		}
	}

	jsonResponse, responseErr := user.BuildSuccessfulJsonResponse(ctx, resultUser)
	if responseErr != nil {
		return nil, models.NewErrorResponse([]error{responseErr}, http.StatusInternalServerError, nil)
	}

	return models.NewSuccessResponse(jsonResponse, http.StatusCreated, nil), nil
}

//ListUsersHandler lists the users in the user pool
func (api *API) ListUsersHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	usersList := models.UsersList{}
	listUserInput := usersList.BuildListUserRequest("", "", int64(0), &api.UserPoolId)
	listUserResp, err := api.CognitoClient.ListUsers(listUserInput)
	if err != nil {
		return nil, models.NewErrorResponse([]error{models.NewCognitoError(ctx, err, "Cognito ListUsers request from create users endpoint")}, http.StatusInternalServerError, nil)
	}

	usersList.Users = usersList.MapCognitoUsers(listUserResp)
	usersList.Count = len(usersList.Users)

	jsonResponse, responseErr := usersList.BuildSuccessfulJsonResponse(ctx)
	if responseErr != nil {
		return nil, models.NewErrorResponse([]error{responseErr}, http.StatusInternalServerError, nil)
	}

	return models.NewSuccessResponse(jsonResponse, http.StatusOK, nil), nil
}

//ChangePasswordHandler processes changes to the users password
func (api *API) ChangePasswordHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	defer req.Body.Close()
	var jsonResponse []byte = nil
	var responseErr error = nil
	var headers map[string]string = nil

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, handleBodyReadError(ctx, err)
	}

	changePasswordParams := models.ChangePassword{}
	err = json.Unmarshal(body, &changePasswordParams)
	if err != nil {
		return nil, handleBodyUnmarshalError(ctx, err)
	}

	if changePasswordParams.ChangeType == models.NewPasswordRequiredType {
		validationErrs := changePasswordParams.ValidateNewPasswordRequiredRequest(ctx)
		if len(validationErrs) != 0 {
			return nil, models.NewErrorResponse(validationErrs, http.StatusBadRequest, nil)
		}

		changePasswordRequest := changePasswordParams.BuildAuthChallengeResponseRequest(api.ClientSecret, api.ClientId, NewPasswordChallenge)

		result, cognitoErr := api.CognitoClient.RespondToAuthChallenge(changePasswordRequest)

		if cognitoErr != nil {
			parsedErr := models.NewCognitoError(ctx, cognitoErr, "RespondToAuthChallenge request, NEW_PASSWORD_REQUIRED type, from change password endpoint")
			if parsedErr.Code == models.InternalError {
				return nil, models.NewErrorResponse([]error{parsedErr}, http.StatusInternalServerError, nil)
			} else if parsedErr.Code == models.InvalidPasswordError || parsedErr.Code == models.InvalidCodeError {
				return nil, models.NewErrorResponse([]error{parsedErr}, http.StatusBadRequest, nil)
			}
		} else {
			jsonResponse, responseErr = changePasswordParams.BuildAuthChallengeSuccessfulJsonResponse(ctx, result)

			if responseErr == nil {
				headers = map[string]string{
					AccessTokenHeaderName:  "Bearer " + *result.AuthenticationResult.AccessToken,
					IdTokenHeaderName:      *result.AuthenticationResult.IdToken,
					RefreshTokenHeaderName: *result.AuthenticationResult.RefreshToken,
				}
			}
		}
	} else if changePasswordParams.ChangeType == models.ForgottenPasswordType {
		// This feature is to be added in a separate ticket later
		err = models.NewValidationError(ctx, models.NotImplementedError, models.NotImplementedDescription)
		return nil, models.NewErrorResponse([]error{err}, http.StatusNotImplemented, nil)
	} else {
		err = models.NewValidationError(ctx, models.UnknownRequestTypeError, models.UnknownPasswordChangeTypeDescription)
		return nil, models.NewErrorResponse([]error{err}, http.StatusBadRequest, nil)
	}

	if responseErr != nil {
		return nil, models.NewErrorResponse([]error{responseErr}, http.StatusInternalServerError, nil)
	}

	return models.NewSuccessResponse(jsonResponse, http.StatusAccepted, headers), nil
}
