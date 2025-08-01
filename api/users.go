package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"

	"github.com/ONSdigital/dp-identity-api/v2/models"
	"github.com/ONSdigital/dp-identity-api/v2/query"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	UsersCreatePermission = "users:create"
	UsersReadPermission   = "users:read"
	UsersUpdatePermission = "users:update"
)

// CreateUserHandler creates a new user and returns a http handler interface
func (api *API) CreateUserHandler(ctx context.Context, _ http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	defer func() {
		if err := req.Body.Close(); err != nil {
			_ = models.NewError(ctx, err, models.BodyCloseError, models.BodyClosedFailedDescription)
		}
	}()

	body, err := io.ReadAll(req.Body)
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

	validationErrs := user.ValidateRegistration(ctx, api.AllowedDomains, api.BlockPlusAddressing)

	listUserInput := models.UsersList{}.BuildListUserRequest("email = \""+user.Email+"\"", "email", int32(1), nil, &api.UserPoolID)
	listUserResp, err := api.CognitoClient.ListUsers(ctx, listUserInput)
	if err != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, models.NewCognitoError(ctx, err, "ListUsers request from create users endpoint"))
	}
	duplicateEmailErr := user.CheckForDuplicateEmail(ctx, listUserResp)
	if duplicateEmailErr != nil {
		validationErrs = append(validationErrs, duplicateEmailErr)
	}

	if len(validationErrs) != 0 {
		return nil, models.NewErrorResponse(http.StatusBadRequest, nil, validationErrs...)
	}

	createUserRequest := user.BuildCreateUserRequest(uuid.NewString(), api.UserPoolID)

	resultUser, err := api.CognitoClient.AdminCreateUser(ctx, createUserRequest)
	if err != nil {
		responseErr := models.NewCognitoError(ctx, err, "AdminCreateUser request from create user endpoint")
		if responseErr.Code == models.InternalError {
			return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
		}
		return nil, models.NewErrorResponse(http.StatusBadRequest, nil, responseErr)
	}
	resultUser.User.Enabled = true
	createdUser := models.UserParams{}.MapCognitoDetails(*resultUser.User)
	jsonResponse, responseErr := createdUser.BuildSuccessfulJSONResponse(ctx)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	return models.NewSuccessResponse(jsonResponse, http.StatusCreated, nil), nil
}

// ListUsersHandler lists the users in the user pool
func (api *API) ListUsersHandler(ctx context.Context, _ http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	var (
		filterString   = aws.String("")
		validationErrs error
	)

	usersList := models.UsersList{}

	if req.URL.Query().Get("active") != "" {
		queryStr := fmt.Sprintf("%s%s", "active=", req.URL.Query().Get("active"))
		*filterString, validationErrs = api.GetFilterStringAndValidate(req.URL.Path, queryStr)
		if validationErrs != nil {
			return nil, models.NewErrorResponse(http.StatusBadRequest, nil, validationErrs)
		}
	}

	listUserResp, errResponse := api.ListUsersWorker(req.Context(), filterString, DefaultBackOffSchedule)
	if errResponse != nil {
		return nil, errResponse
	}

	usersList.SetUsers(listUserResp)

	if req.URL.Query().Get("sort") != "" {
		requestSortQueryErrs := query.SortBy(req.URL.Query().Get("sort"), usersList.Users)
		if requestSortQueryErrs != nil {
			return nil, models.NewErrorResponse(http.StatusBadRequest, nil, requestSortQueryErrs)
		}
	}

	jsonResponse, responseErr := usersList.BuildSuccessfulJSONResponse(ctx)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	return models.NewSuccessResponse(jsonResponse, http.StatusOK, nil), nil
}

// GetUserHandler lists the users in the user pool
func (api *API) GetUserHandler(ctx context.Context, _ http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	vars := mux.Vars(req)
	user := models.UserParams{ID: vars["id"]}
	userInput := user.BuildAdminGetUserRequest(api.UserPoolID)
	userResp, err := api.CognitoClient.AdminGetUser(ctx, userInput)
	if err != nil {
		responseErr := models.NewCognitoError(ctx, err, "AdminGetUser request from get user endpoint")
		if responseErr.Code == models.UserNotFoundError {
			return nil, models.NewErrorResponse(http.StatusNotFound, nil, responseErr)
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	user.MapCognitoGetResponse(userResp)

	jsonResponse, responseErr := user.BuildSuccessfulJSONResponse(ctx)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	return models.NewSuccessResponse(jsonResponse, http.StatusOK, nil), nil
}

// UpdateUserHandler updates a users details in Cognito and returns a http handler interface
func (api *API) UpdateUserHandler(ctx context.Context, _ http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	defer func() {
		if err := req.Body.Close(); err != nil {
			_ = models.NewError(ctx, err, models.BodyCloseError, models.BodyClosedFailedDescription)
		}
	}()
	vars := mux.Vars(req)

	body, err := io.ReadAll(req.Body)
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
		userEnableRequest := user.BuildEnableUserRequest(api.UserPoolID)
		if _, err = api.CognitoClient.AdminEnableUser(ctx, userEnableRequest); err != nil {
			return nil, processUpdateCognitoError(ctx, err, "AdminEnableUser request from update user endpoint")
		}
	} else {
		userDisableRequest := user.BuildDisableUserRequest(api.UserPoolID)
		if _, err = api.CognitoClient.AdminDisableUser(ctx, userDisableRequest); err != nil {
			return nil, processUpdateCognitoError(ctx, err, "AdminDisableUser request from update user endpoint")
		}
	}

	userUpdateRequest := user.BuildUpdateUserRequest(api.UserPoolID)

	_, err = api.CognitoClient.AdminUpdateUserAttributes(ctx, userUpdateRequest)
	if err != nil {
		return nil, processUpdateCognitoError(ctx, err, "AdminUpdateUserAttributes request from update user endpoint")
	}

	userDetailsRequest := user.BuildAdminGetUserRequest(api.UserPoolID)
	userDetailsResponse, err := api.CognitoClient.AdminGetUser(ctx, userDetailsRequest)
	if err != nil {
		responseErr := models.NewCognitoError(ctx, err, "AdminGetUser request from update user endpoint")
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	user.MapCognitoGetResponse(userDetailsResponse)

	jsonResponse, responseErr := user.BuildSuccessfulJSONResponse(ctx)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	return models.NewSuccessResponse(jsonResponse, http.StatusOK, nil), nil
}

// UserSetPasswordHandler sets a users password to a generated password in Cognito and returns a http handler interface
func (api *API) UserSetPasswordHandler(ctx context.Context, _ http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	vars := mux.Vars(req)
	userID := vars["id"]

	user := models.UserParams{ID: vars["id"]}
	userInput := user.BuildAdminGetUserRequest(api.UserPoolID)

	userResp, err := api.CognitoClient.AdminGetUser(ctx, userInput)
	if err != nil {
		responseErr := models.NewCognitoError(ctx, err, "AdminGetUser request from user set password endpoint")
		if responseErr.Code == models.UserNotFoundError {
			return nil, models.NewErrorResponse(http.StatusNotFound, nil, responseErr)
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	user.MapCognitoGetResponse(userResp)

	validationErrs := user.ValidateSetPasswordRequest(ctx)
	if len(validationErrs) != 0 {
		return nil, models.NewErrorResponse(http.StatusForbidden, nil, validationErrs...)
	}

	err = user.GeneratePassword(ctx)
	if err != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, err)
	}

	userSetPasswordInput := user.BuildSetPasswordRequest(api.UserPoolID)

	_, err = api.CognitoClient.AdminSetUserPassword(ctx, userSetPasswordInput)
	if err != nil {
		responseErr := models.NewCognitoError(ctx, err, "error whilst resetting user account")

		switch responseErr.Code {
		case models.LimitExceededError, models.TooManyRequestsError:
			log.Error(ctx, "cognito request limit exceeded", responseErr, log.Data{"userID": userID})
			return nil, models.NewErrorResponse(http.StatusBadRequest, nil, responseErr)
		case models.UserNotFoundError:
			log.Error(ctx, "user not found", responseErr, log.Data{"userID": userID})
			return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
		}
	}

	log.Info(ctx, "user set password completed", log.Data{"userID": userID})

	return models.NewSuccessResponse(nil, http.StatusAccepted, nil), nil
}

func processUpdateCognitoError(ctx context.Context, err error, errContext string) *models.ErrorResponse {
	responseErr := models.NewCognitoError(ctx, err, errContext)

	switch responseErr.Code {
	case models.UserNotFoundError:
		return models.NewErrorResponse(http.StatusNotFound, nil, responseErr)
	case models.InvalidFieldError:
		return models.NewErrorResponse(http.StatusBadRequest, nil, responseErr)
	default:
		return models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}
}

// ChangePasswordHandler processes changes to the users password
func (api *API) ChangePasswordHandler(ctx context.Context, _ http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	defer func() {
		if err := req.Body.Close(); err != nil {
			_ = models.NewError(ctx, err, models.BodyCloseError, models.BodyClosedFailedDescription)
		}
	}()
	var jsonResponse []byte
	var responseErr error
	var headers map[string]string

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, handleBodyReadError(ctx, err)
	}

	changePasswordParams := models.ChangePassword{}
	err = json.Unmarshal(body, &changePasswordParams)
	if err != nil {
		return nil, handleBodyUnmarshalError(ctx, err)
	}

	//nolint:staticcheck // making this into a switch statement would not improve it
	// that much. TODO: It needs a greater level of refactoring
	if changePasswordParams.ChangeType == models.NewPasswordRequiredType {
		validationErrs := changePasswordParams.ValidateNewPasswordRequiredRequest(ctx)
		if len(validationErrs) != 0 {
			return nil, models.NewErrorResponse(http.StatusBadRequest, nil, validationErrs...)
		}

		changePasswordRequest := changePasswordParams.BuildAuthChallengeResponseRequest(api.ClientSecret, api.ClientID, NewPasswordChallenge)

		result, cognitoErr := api.CognitoClient.RespondToAuthChallenge(ctx, changePasswordRequest)

		if cognitoErr != nil {
			parsedErr := models.NewCognitoError(ctx, cognitoErr, "RespondToAuthChallenge request from change password endpoint")
			if parsedErr.Code == models.InternalError {
				return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, parsedErr)
			} else if parsedErr.Code == models.InvalidPasswordError || parsedErr.Code == models.InvalidCodeError {
				return nil, models.NewErrorResponse(http.StatusBadRequest, nil, parsedErr)
			}
		} else {
			// Determine the refresh token TTL (DescribeUserPoolClient)
			userPoolClient, err := api.CognitoClient.DescribeUserPoolClient(ctx,
				&cognitoidentityprovider.DescribeUserPoolClientInput{
					UserPoolId: &api.UserPoolID,
					ClientId:   &api.ClientID,
				},
			)
			if err != nil {
				awsErr := models.NewCognitoError(ctx, err, "Describing user pool for refresh token TTL")
				return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, awsErr)
			}

			clientTokenValidityUnits := *userPoolClient.UserPoolClient.TokenValidityUnits
			refreshTokenTTL := calculateTokenTTLInSeconds(clientTokenValidityUnits.RefreshToken, int(userPoolClient.UserPoolClient.RefreshTokenValidity))

			jsonResponse, responseErr = changePasswordParams.BuildAuthChallengeSuccessfulJSONResponse(ctx, result, refreshTokenTTL)
			if responseErr == nil {
				headers = map[string]string{
					AccessTokenHeaderName:  "Bearer " + *result.AuthenticationResult.AccessToken,
					IDTokenHeaderName:      *result.AuthenticationResult.IdToken,
					RefreshTokenHeaderName: *result.AuthenticationResult.RefreshToken,
				}
			}
		}
	} else if changePasswordParams.ChangeType == models.ForgottenPasswordType {
		validationErrs := changePasswordParams.ValidateForgottenPasswordRequest(ctx)
		if len(validationErrs) != 0 {
			return nil, models.NewErrorResponse(http.StatusBadRequest, nil, validationErrs...)
		}
		changeForgottenPasswordRequest := changePasswordParams.BuildConfirmForgotPasswordRequest(api.ClientSecret, api.ClientID)

		_, cognitoErr := api.CognitoClient.ConfirmForgotPassword(ctx, changeForgottenPasswordRequest)

		if cognitoErr != nil {
			parsedErr := models.NewCognitoError(ctx, cognitoErr, "ConfirmForgottenPassword request from change password endpoint")
			if parsedErr.Code == models.InternalError {
				return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, parsedErr)
			} else if parsedErr.Code == models.InvalidPasswordError || parsedErr.Code == models.InvalidCodeError || parsedErr.Code == models.ExpiredCodeError {
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

// PasswordResetHandler requests a password reset email be sent to the user and returns a http handler interface
func (api *API) PasswordResetHandler(ctx context.Context, _ http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	defer func() {
		if err := req.Body.Close(); err != nil {
			_ = models.NewError(ctx, err, models.BodyCloseError, models.BodyClosedFailedDescription)
		}
	}()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Error(ctx, "failed to read request body", err)
		return nil, handleBodyReadError(ctx, err)
	}

	passwordResetParams := models.PasswordReset{}
	err = json.Unmarshal(body, &passwordResetParams)
	if err != nil {
		log.Error(ctx, "failed to unmarshal password reset passwords", err)
		return nil, handleBodyUnmarshalError(ctx, err)
	}

	validationErr := passwordResetParams.Validate(ctx)

	if validationErr != nil {
		log.Error(ctx, "failed validation", validationErr, log.Data{"user_email": passwordResetParams.Email})
		return nil, models.NewErrorResponse(http.StatusBadRequest, nil, validationErr)
	}

	log.Info(ctx, "request reset parameters validated", log.Data{"user_email": passwordResetParams.Email})

	forgotPasswordRequest := passwordResetParams.BuildCognitoRequest(api.ClientSecret, api.ClientID)

	_, err = api.CognitoClient.ForgotPassword(ctx, forgotPasswordRequest)
	if err != nil {
		responseErr := models.NewCognitoError(ctx, err, "ForgotPassword request from password reset endpoint")

		if responseErr.Code == models.LimitExceededError || responseErr.Code == models.TooManyRequestsError {
			log.Error(ctx, "cognito request limit exceeded", responseErr, log.Data{"user_email": passwordResetParams.Email})
			return nil, models.NewErrorResponse(http.StatusBadRequest, nil, responseErr)
		} else if responseErr.Code != models.UserNotFoundError && responseErr.Code != models.UserNotConfirmedError {
			log.Error(ctx, "user not found or user not confirmed", responseErr, log.Data{"user_email": passwordResetParams.Email})
			return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
		}
	}

	log.Info(ctx, "password reset completed", log.Data{"user_email": passwordResetParams.Email})

	return models.NewSuccessResponse(nil, http.StatusAccepted, nil), nil
}

// List Groups for user pagination allows first call and then any other call if nextToken is not ""
func (api *API) getGroupsForUser(ctx context.Context, listOfGroups []types.GroupType, userID models.UserParams) ([]types.GroupType, error) {
	firstTimeCheck := false
	var nextToken string
	for !firstTimeCheck || nextToken != "" {
		firstTimeCheck = true

		userGroupsRequest := userID.BuildListUserGroupsRequest(api.UserPoolID, nextToken)
		userGroupsResponse, err := api.CognitoClient.AdminListGroupsForUser(ctx, userGroupsRequest)
		if err != nil {
			return nil, err
		}

		listOfGroups = append(listOfGroups, userGroupsResponse.Groups...)
		nextToken = ""
		if userGroupsResponse.NextToken != nil {
			nextToken = *userGroupsResponse.NextToken
		}
	}
	return listOfGroups, nil
}

// ListUserGroupsHandler lists the users in the user pool
func (api *API) ListUserGroupsHandler(ctx context.Context, _ http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	vars := mux.Vars(req)
	userID := models.UserParams{ID: vars["id"]}
	var listofgroupsInput []types.GroupType
	finalUserResponse := cognitoidentityprovider.AdminListGroupsForUserOutput{}
	listusergroups := models.ListUserGroups{}

	listofGroupsOutput, err := api.getGroupsForUser(ctx, listofgroupsInput, userID)
	if err != nil {
		cognitoErr := models.NewCognitoError(ctx, err, "Cognito ListofUserGroups request from list user groups endpoint")
		if cognitoErr.Code == models.NotFoundError {
			return nil, models.NewErrorResponse(http.StatusNotFound, nil, cognitoErr)
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, cognitoErr)
	}
	finalUserResponse.Groups = append(finalUserResponse.Groups, listofGroupsOutput...)
	jsonResponse, responseErr := listusergroups.BuildListUserGroupsSuccessfulJSONResponse(ctx, &finalUserResponse)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}
	return models.NewSuccessResponse(jsonResponse, http.StatusOK, nil), nil
}

func (api *API) GetFilterStringAndValidate(path, queryStr string) (string, error) {
	ctx := context.Background()

	if api.APIRequestFilter[path] != nil && api.APIRequestFilter[path][queryStr] != "" {
		return api.APIRequestFilter[path][queryStr], nil
	}
	return "", models.NewValidationError(ctx, models.InvalidFilterQuery, models.InvalidFilterQueryDescription)
}
