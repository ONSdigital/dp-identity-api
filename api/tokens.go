package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"

	"github.com/ONSdigital/dp-identity-api/v2/models"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

// TokensHandler uses submitted email address and password to sign a user in against Cognito and returns a http handler interface
func (api *API) TokensHandler(ctx context.Context, _ http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	defer func() {
		if err := req.Body.Close(); err != nil {
			_ = models.NewError(ctx, err, models.BodyCloseError, models.BodyClosedFailedDescription)
		}
	}()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, handleBodyReadError(ctx, err)
	}
	var userSignIn models.UserSignIn
	err = json.Unmarshal(body, &userSignIn)
	if err != nil {
		return nil, handleBodyUnmarshalError(ctx, err)
	}
	validationErrs := userSignIn.ValidateCredentials(ctx)
	if validationErrs != nil {
		return nil, models.NewErrorResponse(http.StatusBadRequest, nil, *validationErrs...)
	}
	input := userSignIn.BuildCognitoRequest(api.ClientID, api.ClientSecret, api.ClientAuthFlow)
	result, authErr := api.CognitoClient.InitiateAuth(ctx, input)
	if authErr != nil {
		responseErr := models.NewCognitoError(ctx, authErr, "Cognito InitiateAuth request from sign in handler")
		if responseErr.Code == models.InternalError {
			return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
		}
		switch responseErr.Description {
		case models.SignInFailedDescription:
			// Returning `WWW-Authenticate` in header as part of http.StatusUnauthorized response
			// See here: https://datatracker.ietf.org/doc/html/rfc7235#section-4.1
			headers := map[string]string{
				WWWAuthenticateName: "Bearer realm=\"" + ONSRealm + "\", charset=\"" + Charset + "\"",
			}
			return nil, models.NewErrorResponse(http.StatusUnauthorized, headers, responseErr)
		case models.SignInAttemptsExceededDescription:
			// Cognito returns the same Code for invalid credentials and too many attempts errors, changing our Error.Code to enable differentiation in the client
			responseErr.Code = models.TooManyFailedAttemptsError
			return nil, models.NewErrorResponse(http.StatusForbidden, nil, responseErr)
		default:
			return nil, models.NewErrorResponse(http.StatusBadRequest, nil, responseErr)
		}
	}
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

	jsonResponse, responseErr := userSignIn.BuildSuccessfulJSONResponse(ctx, result, refreshTokenTTL)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	// success headers
	var headers map[string]string
	if result.AuthenticationResult != nil {
		headers = map[string]string{
			AccessTokenHeaderName:  "Bearer " + *result.AuthenticationResult.AccessToken,
			IDTokenHeaderName:      *result.AuthenticationResult.IdToken,
			RefreshTokenHeaderName: *result.AuthenticationResult.RefreshToken,
		}
	} else {
		headers = nil
	}

	// response - http.StatusCreated by default
	httpStatus := http.StatusCreated
	if result.ChallengeName == NewPasswordChallenge {
		httpStatus = http.StatusAccepted
	}

	return models.NewSuccessResponse(jsonResponse, httpStatus, headers), nil
}

// SignOutHandler invalidates a users access token signing them out and returns a http handler interface
func (api *API) SignOutHandler(ctx context.Context, _ http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	accessToken := models.AccessToken{
		AuthHeader: req.Header.Get(AccessTokenHeaderName),
	}
	validationErr := accessToken.Validate(ctx)
	if validationErr != nil {
		return nil, models.NewErrorResponse(http.StatusBadRequest, nil, validationErr)
	}

	_, err := api.CognitoClient.GlobalSignOut(ctx, accessToken.GenerateSignOutRequest())

	if err != nil {
		responseErr := models.NewCognitoError(ctx, err, "Cognito GlobalSignOut request for sign out")
		if responseErr.Code == models.NotFoundError || responseErr.Code == models.NotAuthorisedError {
			return nil, models.NewErrorResponse(http.StatusBadRequest, nil, responseErr)
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}
	return models.NewSuccessResponse(nil, http.StatusNoContent, nil), nil
}

// RefreshHandler refreshes a users access token and returns new access and ID tokens, expiration time and the refresh token
func (api *API) RefreshHandler(ctx context.Context, _ http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	var validationErrs []error
	refreshToken := models.RefreshToken{TokenString: req.Header.Get(RefreshTokenHeaderName)}
	validationErr := refreshToken.Validate(ctx)
	if validationErr != nil {
		validationErrs = append(validationErrs, validationErr)
	}

	idToken := models.IDToken{TokenString: req.Header.Get(IDTokenHeaderName)}
	validationErr = idToken.Validate(ctx)
	if validationErr != nil {
		validationErrs = append(validationErrs, validationErr)
	}

	if len(validationErrs) > 0 {
		return nil, models.NewErrorResponse(http.StatusBadRequest, nil, validationErrs...)
	}

	authInput := refreshToken.GenerateRefreshRequest(api.ClientSecret, idToken.Claims.CognitoUser, api.ClientID)
	result, authErr := api.CognitoClient.InitiateAuth(ctx, authInput)

	if authErr != nil {
		responseErr := models.NewCognitoError(ctx, authErr, "Cognito InitiateAuth request for token refresh")
		if responseErr.Code == models.NotAuthorisedError {
			return nil, models.NewErrorResponse(http.StatusForbidden, nil, responseErr)
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	jsonResponse, responseErr := refreshToken.BuildSuccessfulJSONResponse(ctx, result)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	headers := map[string]string{
		AccessTokenHeaderName: "Bearer " + *result.AuthenticationResult.AccessToken,
		IDTokenHeaderName:     *result.AuthenticationResult.IdToken,
	}

	return models.NewSuccessResponse(jsonResponse, http.StatusCreated, headers), nil
}

// SignOutAllUsersHandler bulk refresh token invalidation for panic sign out handling
func (api *API) SignOutAllUsersHandler(ctx context.Context, _ http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	var (
		userFilterString = `status="Enabled"`
	)
	usersList, awsErr := api.ListUsersWorker(req.Context(), &userFilterString, DefaultBackOffSchedule)
	if awsErr != nil {
		return nil, awsErr
	}
	globalSignOut := &models.GlobalSignOut{
		ResultsChannel:  make(chan string, len(*usersList)),
		BackoffSchedule: DefaultBackOffSchedule,
		RetryAllowed:    true,
	}
	// run api.SignOutUsersWorker concurrently
	go api.SignOutUsersWorker(req.Context(), globalSignOut, usersList)

	postBody, resErr := models.BuildSuccessfulSignOutAllUsersJSONResponse(ctx)
	if resErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, resErr)
	}

	return models.NewSuccessResponse(postBody, http.StatusAccepted, nil), nil
}

// ListUsersWorker - generates a list of users based on `userFilterString` filter string
func (api *API) ListUsersWorker(ctx context.Context, userFilterString *string, backoffSchedule []time.Duration) (*[]models.UserParams, *models.ErrorResponse) {
	var (
		awsErr                error
		usersList             models.UsersList
		listUsersResp, result *cognitoidentityprovider.ListUsersOutput
		listUserInput         = usersList.BuildListUserRequest(
			*userFilterString,
			"",
			int32(0),
			nil,
			&api.UserPoolID,
		)
		usersListError *models.ErrorResponse
	)
	listUsersResp, awsErr = api.generateListUsersRequest(ctx, listUserInput)
	if awsErr != nil {
		err := models.NewCognitoError(ctx, awsErr, "Cognito ListUsers request from signout all users from group endpoint")
		usersListError = models.NewErrorResponse(http.StatusInternalServerError, nil, err)
	} else if listUsersResp.PaginationToken != nil {
		listUserInput.PaginationToken = listUsersResp.PaginationToken
		// set `loadingInProgress` to control requesting new list data
		loadingInProgress := true
		for loadingInProgress {
			for _, backoff := range backoffSchedule {
				result, awsErr = api.generateListUsersRequest(ctx, listUserInput)
				if awsErr == nil {
					listUsersResp.Users = append(listUsersResp.Users, result.Users...)
					if result.PaginationToken != nil {
						listUserInput.PaginationToken = result.PaginationToken
						break
					}
					loadingInProgress = false
					break
				}
				err := models.NewCognitoError(ctx, awsErr, "Cognito ListUsers request from signout all users from group endpoint")
				if err.Code != models.TooManyRequestsError {
					usersListError = models.NewErrorResponse(http.StatusInternalServerError, nil, err)
					loadingInProgress = false
					break
				}
				time.Sleep(backoff)
			}
		}
	}
	if usersListError != nil {
		return nil, usersListError
	}
	usersList.MapCognitoUsers(&listUsersResp.Users)
	return &usersList.Users, nil
}

// SignOutUsersWorker - signs out users globally by invalidating user's refresh token
func (api *API) SignOutUsersWorker(ctx context.Context, g *models.GlobalSignOut, usersList *[]models.UserParams) {
	userSignOutRequestData := g.BuildSignOutUserRequest(usersList, &api.UserPoolID)

	for _, userSignoutRequest := range userSignOutRequestData {
		for _, backoff := range g.BackoffSchedule {
			_, err := api.generateGlobalSignOutRequest(ctx, userSignoutRequest)

			// no errors returned - add username to results channel and break to next user in list
			if err == nil {
				// the results channel is not being processed by caller currently - here for possible future extensibility
				g.ResultsChannel <- *userSignoutRequest.Username
				g.RetryAllowed = true
				break
			}

			// error returned - process it
			responseErr := models.NewCognitoError(ctx, err, "Cognito AdminUserGlobalSignOut request for sign out")

			// is error code != models.TooManyRequestsError? If so, proceed...
			if responseErr.Code != models.TooManyRequestsError {
				// if g.RetryAllowed is true, allowed to request api again
				if g.RetryAllowed {
					// attempt one more request to api
					g.RetryAllowed = false // 3. Set GlobalSignOut.RetryAllowed to false and request AdminUserGlobalSignOut
					_, retryErr := api.generateGlobalSignOutRequest(ctx, userSignoutRequest)

					if retryErr != nil {
						// if error response from request received again, process it
						retryResponseErr := models.NewCognitoError(ctx, err, "Cognito AdminUserGlobalSignOut request for sign out")

						// if error code != models.TooManyRequestsError break to next user
						if retryResponseErr.Code != models.TooManyRequestsError {
							g.RetryAllowed = true

							break
						}
					} else {
						// no error on retry, add user to results channel and break to next user
						// the results channel is not being processed by caller currently - here for possible future extensibility
						g.ResultsChannel <- *userSignoutRequest.Username
						g.RetryAllowed = true

						break
					}
				} else {
					// if GlobalSignOut.RetryAllowed already false break to next user
					g.RetryAllowed = true

					break
				}
			}

			// backoff for predetermined length of time before requesting again
			time.Sleep(backoff)
		}
	}
	close(g.ResultsChannel)
}

// generateGlobalSignOutRequest - local routine to generete the global signout request per user
func (api *API) generateGlobalSignOutRequest(ctx context.Context, user *cognitoidentityprovider.AdminUserGlobalSignOutInput) (*cognitoidentityprovider.AdminUserGlobalSignOutOutput, error) {
	return api.CognitoClient.AdminUserGlobalSignOut(ctx, user)
}

// generateListUsersRequest - local routine to generate a list users request
func (api *API) generateListUsersRequest(ctx context.Context, input *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error) {
	return api.CognitoClient.ListUsers(ctx, input)
}

// calculateTokenTTLInSeconds takes a token unit and a number, as received from Cognito, and
// returns the number of seconds.
func calculateTokenTTLInSeconds(unit types.TimeUnitsType, number int) int {
	switch unit {
	case types.TimeUnitsTypeDays:
		return number * models.SecondsInDay
	case types.TimeUnitsTypeHours:
		return number * 3600
	case types.TimeUnitsTypeMinutes:
		return number * 60
	case types.TimeUnitsTypeSeconds:
		return number
	}
	return 0
}
