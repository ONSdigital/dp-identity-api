package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

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

	err = user.GeneratePassword(ctx)
	if err != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, err)
	}

	validationErrs := user.ValidateRegistration(ctx)

	listUserInput := models.UsersList{}.BuildListUserRequest("email = \""+user.Email+"\"", "email", int64(1), &api.UserPoolId)
	listUserResp, err := api.CognitoClient.ListUsers(listUserInput)
	if err != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, models.NewCognitoError(ctx, err, "Cognito ListUsers request from create users endpoint"))
	}
	duplicateEmailErr := user.CheckForDuplicateEmail(ctx, listUserResp)
	if duplicateEmailErr != nil {
		validationErrs = append(validationErrs, duplicateEmailErr)
	}

	if len(validationErrs) != 0 {
		return nil, models.NewErrorResponse(http.StatusBadRequest, nil, validationErrs...)
	}

	createUserRequest := user.BuildCreateUserRequest(uuid.NewString(), api.UserPoolId)

	resultUser, err := api.CognitoClient.AdminCreateUser(createUserRequest)
	if err != nil {
		responseErr := models.NewCognitoError(ctx, err, "AdminCreateUser request from create user endpoint")
		if responseErr.Code == models.InternalError {
			return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
		} else {
			return nil, models.NewErrorResponse(http.StatusBadRequest, nil, responseErr)
		}
	}

	createdUser := models.UserParams{}.MapCognitoDetails(resultUser.User)
	jsonResponse, responseErr := createdUser.BuildSuccessfulJsonResponse(ctx)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	return models.NewSuccessResponse(jsonResponse, http.StatusCreated, nil), nil
}

//ListUsersHandler lists the users in the user pool
func (api *API) ListUsersHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	usersList := models.UsersList{}
	listUserInput := usersList.BuildListUserRequest("", "", int64(0), &api.UserPoolId)
	listUserResp, err := api.CognitoClient.ListUsers(listUserInput)
	if err != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, models.NewCognitoError(ctx, err, "Cognito ListUsers request from create users endpoint"))
	}

	usersList.MapCognitoUsers(listUserResp)

	jsonResponse, responseErr := usersList.BuildSuccessfulJsonResponse(ctx)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	return models.NewSuccessResponse(jsonResponse, http.StatusOK, nil), nil
}

//GetUserHandler lists the users in the user pool
func (api *API) GetUserHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	vars := mux.Vars(req)
	user := models.UserParams{ID: vars["id"]}
	userInput := user.BuildAdminGetUserRequest(api.UserPoolId)
	userResp, err := api.CognitoClient.AdminGetUser(userInput)
	if err != nil {
		responseErr := models.NewCognitoError(ctx, err, "Cognito ListUsers request from create users endpoint")
		if responseErr.Code == models.UserNotFoundError {
			return nil, models.NewErrorResponse(http.StatusNotFound, nil, responseErr)
		} else {
			return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
		}
	}

	user.MapCognitoGetResponse(userResp)

	jsonResponse, responseErr := user.BuildSuccessfulJsonResponse(ctx)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	return models.NewSuccessResponse(jsonResponse, http.StatusOK, nil), nil
}

//UpdateUserHandler updates a users details in Cognito and returns a http handler interface
func (api *API) UpdateUserHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	defer req.Body.Close()
	vars := mux.Vars(req)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, handleBodyReadError(ctx, err)
	}

	user := models.UserParams{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, handleBodyUnmarshalError(ctx, err)
	}
	user.ID = vars["id"]

	validationErrs := user.ValidateUpdate(ctx)

	if len(validationErrs) != 0 {
		return nil, models.NewErrorResponse(http.StatusBadRequest, nil, validationErrs...)
	}

	if user.Active {
		userEnableRequest := user.BuildEnableUserRequest(api.UserPoolId)
		if _, err = api.CognitoClient.AdminEnableUser(userEnableRequest); err != nil {
			return nil, processUpdateCognitoError(ctx, err, "AdminEnableUser request from update user endpoint")
		}
	} else {
		userDisableRequest := user.BuildDisableUserRequest(api.UserPoolId)
		if _, err = api.CognitoClient.AdminDisableUser(userDisableRequest); err != nil {
			return nil, processUpdateCognitoError(ctx, err, "AdminDisableUser request from update user endpoint")
		}
	}

	userUpdateRequest := user.BuildUpdateUserRequest(api.UserPoolId)

	_, err = api.CognitoClient.AdminUpdateUserAttributes(userUpdateRequest)
	if err != nil {
		return nil, processUpdateCognitoError(ctx, err, "AdminUpdateUserAttributes request from update user endpoint")
	}

	userDetailsRequest := user.BuildAdminGetUserRequest(api.UserPoolId)
	userDetailsResponse, err := api.CognitoClient.AdminGetUser(userDetailsRequest)
	if err != nil {
		responseErr := models.NewCognitoError(ctx, err, "AdminGetUser request from update user endpoint")
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	user.MapCognitoGetResponse(userDetailsResponse)

	jsonResponse, responseErr := user.BuildSuccessfulJsonResponse(ctx)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	return models.NewSuccessResponse(jsonResponse, http.StatusOK, nil), nil
}

func processUpdateCognitoError(ctx context.Context, err error, errContext string) *models.ErrorResponse {
	responseErr := models.NewCognitoError(ctx, err, errContext)
	if responseErr.Code == models.UserNotFoundError {
		return models.NewErrorResponse(http.StatusNotFound, nil, responseErr)
	} else if responseErr.Code == models.InvalidFieldError {
		return models.NewErrorResponse(http.StatusBadRequest, nil, responseErr)
	}
	return models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
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
			return nil, models.NewErrorResponse(http.StatusBadRequest, nil, validationErrs...)
		}

		changePasswordRequest := changePasswordParams.BuildAuthChallengeResponseRequest(api.ClientSecret, api.ClientId, NewPasswordChallenge)

		result, cognitoErr := api.CognitoClient.RespondToAuthChallenge(changePasswordRequest)

		if cognitoErr != nil {
			parsedErr := models.NewCognitoError(ctx, cognitoErr, "RespondToAuthChallenge request from change password endpoint")
			if parsedErr.Code == models.InternalError {
				return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, parsedErr)
			} else if parsedErr.Code == models.InvalidPasswordError || parsedErr.Code == models.InvalidCodeError {
				return nil, models.NewErrorResponse(http.StatusBadRequest, nil, parsedErr)
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
		validationErrs := changePasswordParams.ValidateForgottenPasswordRequiredRequest(ctx)
		if len(validationErrs) != 0 {
			return nil, models.NewErrorResponse(http.StatusBadRequest, nil, validationErrs...)
		}
		changeForgottenPasswordRequest := changePasswordParams.BuildConfirmForgotPasswordRequest(api.ClientSecret, api.ClientId)

		_, cognitoErr := api.CognitoClient.ConfirmForgotPassword(changeForgottenPasswordRequest)

		if cognitoErr != nil {
			parsedErr := models.NewCognitoError(ctx, cognitoErr, "ConfirmForgottenPassword request from change password endpoint")
			if parsedErr.Code == models.InternalError {
				return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, parsedErr)
			} else if parsedErr.Code == models.InvalidPasswordError || parsedErr.Code == models.InvalidCodeError {
				return nil, models.NewErrorResponse(http.StatusBadRequest, nil, parsedErr)
			}
		}
	} else {
		err = models.NewValidationError(ctx, models.UnknownRequestTypeError, models.UnknownPasswordChangeTypeDescription)
		return nil, models.NewErrorResponse(http.StatusBadRequest, nil, err)
	}

	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	return models.NewSuccessResponse(jsonResponse, http.StatusAccepted, headers), nil
}

//PasswordResetHandler requests a password reset email be sent to the user and returns a http handler interface
func (api *API) PasswordResetHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, handleBodyReadError(ctx, err)
	}

	passwordResetParams := models.PasswordReset{}
	err = json.Unmarshal(body, &passwordResetParams)
	if err != nil {
		return nil, handleBodyUnmarshalError(ctx, err)
	}

	validationErr := passwordResetParams.Validate(ctx)

	if validationErr != nil {
		return nil, models.NewErrorResponse(http.StatusBadRequest, nil, validationErr)
	}

	forgotPasswordRequest := passwordResetParams.BuildCognitoRequest(api.ClientSecret, api.ClientId)

	_, err = api.CognitoClient.ForgotPassword(forgotPasswordRequest)
	if err != nil {
		responseErr := models.NewCognitoError(ctx, err, "ForgotPassword request from password reset endpoint")
		if responseErr.Code == models.LimitExceededError || responseErr.Code == models.TooManyRequestsError {
			return nil, models.NewErrorResponse(http.StatusBadRequest, nil, responseErr)
		} else if responseErr.Code != models.UserNotFoundError && responseErr.Code != models.UserNotConfirmedError {
			return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
		}
	}

	return models.NewSuccessResponse(nil, http.StatusAccepted, nil), nil
}
