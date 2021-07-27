package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/ONSdigital/log.go/log"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

//TokensHandler uses submitted email address and password to sign a user in against Cognito and returns a http handler interface
func (api *API) TokensHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
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

	input := userSignIn.BuildCognitoRequest(api.ClientId, api.ClientSecret, api.ClientAuthFlow)
	result, authErr := api.CognitoClient.InitiateAuth(input)

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

	jsonResponse, responseErr := userSignIn.BuildSuccessfulJsonResponse(ctx, result)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	// success headers
	var headers map[string]string
	if result.AuthenticationResult != nil {
		headers = map[string]string{
			AccessTokenHeaderName:  "Bearer " + *result.AuthenticationResult.AccessToken,
			IdTokenHeaderName:      *result.AuthenticationResult.IdToken,
			RefreshTokenHeaderName: *result.AuthenticationResult.RefreshToken,
		}
	} else {
		headers = nil
	}

	// response - http.StatusCreated by default
	httpStatus := http.StatusCreated
	if result.ChallengeName != nil && *result.ChallengeName == NewPasswordChallenge {
		httpStatus = http.StatusAccepted
	}

	return models.NewSuccessResponse(jsonResponse, httpStatus, headers), nil
}

//SignOutHandler invalidates a users access token signing them out and returns a http handler interface
func (api *API) SignOutHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	accessToken := models.AccessToken{
		AuthHeader: req.Header.Get(AccessTokenHeaderName),
	}
	validationErr := accessToken.Validate(ctx)
	if validationErr != nil {
		return nil, models.NewErrorResponse(http.StatusBadRequest, nil, validationErr)
	}

	_, err := api.CognitoClient.GlobalSignOut(accessToken.GenerateSignOutRequest())

	if err != nil {
		responseErr := models.NewCognitoError(ctx, err, "Cognito GlobalSignOut request for sign out")
		if responseErr.Code == models.NotFoundError || responseErr.Code == models.NotAuthorisedError {
			return nil, models.NewErrorResponse(http.StatusBadRequest, nil, responseErr)
		} else {
			return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
		}
	}

	return models.NewSuccessResponse(nil, http.StatusNoContent, nil), nil
}

//RefreshHandler refreshes a users access token and returns new access and ID tokens, expiration time and the refresh token
func (api *API) RefreshHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	var validationErrs []error
	refreshToken := models.RefreshToken{TokenString: req.Header.Get(RefreshTokenHeaderName)}
	validationErr := refreshToken.Validate(ctx)
	if validationErr != nil {
		validationErrs = append(validationErrs, validationErr)
	}

	idToken := models.IdToken{TokenString: req.Header.Get(IdTokenHeaderName)}
	validationErr = idToken.Validate(ctx)
	if validationErr != nil {
		validationErrs = append(validationErrs, validationErr)
	}

	if len(validationErrs) > 0 {
		return nil, models.NewErrorResponse(http.StatusBadRequest, nil, validationErrs...)
	}

	authInput := refreshToken.GenerateRefreshRequest(api.ClientSecret, idToken.Claims.CognitoUser, api.ClientId)
	result, authErr := api.CognitoClient.InitiateAuth(authInput)

	if authErr != nil {
		responseErr := models.NewCognitoError(ctx, authErr, "Cognito InitiateAuth request for token refresh")
		if responseErr.Code == models.NotAuthorisedError {
			return nil, models.NewErrorResponse(http.StatusForbidden, nil, responseErr)
		} else {
			return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
		}
	}

	jsonResponse, responseErr := refreshToken.BuildSuccessfulJsonResponse(ctx, result)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	headers := map[string]string{
		AccessTokenHeaderName: "Bearer " + *result.AuthenticationResult.AccessToken,
		IdTokenHeaderName:     *result.AuthenticationResult.IdToken,
	}

	return models.NewSuccessResponse(jsonResponse, http.StatusCreated, headers), nil
}

//SignOutAllUsersHandler bulk refresh token invalidation for panic sign out handling
func (api *API) SignOutAllUsersHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	usersList := models.UsersList{}
	listUserInput := usersList.BuildListUserRequest("status=\"Enabled\"", "username", int64(0), &api.UserPoolId)
	listUserResp, err := api.CognitoClient.ListUsers(listUserInput)
	if err != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, err)
	}
	usersList.MapCognitoUsers(listUserResp)
	globalSignOut := &models.GlobalSignOut{
		ResultsChannel: make(chan string, len(usersList.Users)),
		BackoffSchedule: []time.Duration{
			1 * time.Second,
			3 * time.Second,
			10 * time.Second,
		},
		RetryAllowed: true,
	}
	// run api.SignOutUsersWorker concurrently
	go api.SignOutUsersWorker(ctx, globalSignOut, &usersList.Users)
	return models.NewSuccessResponse(nil, http.StatusAccepted, nil), nil
}

// SignOutUsersWorker - signs out users globally by invalidating user's refresh token
func (api *API) SignOutUsersWorker(ctx context.Context, g *models.GlobalSignOut, usersList *[]models.UserParams) {
	userSignOutRequestData := g.BuildSignOutUserRequest(usersList, &api.UserPoolId)
	for i := range userSignOutRequestData {
		for _, backoff := range g.BackoffSchedule {
			_, err := api.generateGlobalSignOutRequest(userSignOutRequestData[i])
			if err != nil {
				responseErr := models.NewCognitoError(ctx, err, "Cognito AdminUserGlobalSignOut request for sign out")
				if responseErr.Code != models.TooManyRequestsError { // 1. Process all errors other than TooManyRequestsException (429)
					if g.RetryAllowed { // 2. Attempt one more request to AdminUserGlobalSignOut if GlobalSignOut.RetryAllowed true, else break to next user
						g.RetryAllowed = false // 3. Set GlobalSignOut.RetryAllowed to false and request AdminUserGlobalSignOut
						log.Event(ctx, "Error Cognito AdminUserGlobalSignOut:", log.ERROR, log.Error(err))
						_, retryErr := api.generateGlobalSignOutRequest(userSignOutRequestData[i])
						if retryErr != nil { // 4. If error response from request received, process it
							retryResponseErr := models.NewCognitoError(ctx, err, "Cognito AdminUserGlobalSignOut request for sign out")
							if retryResponseErr.Code != models.TooManyRequestsError { // 4.1 If error is not a TooManyRequestsException (429), reset GlobalSignOut.RetryAllowed and break to next user
								g.RetryAllowed = true
								break
							}
						} else { // 5. If request successful, reset GlobalSignOut.RetryAllowed and set original response error to nil - username will be added to GlobalSignOut.ResultsChannel
							g.RetryAllowed = true
							err = nil
						}
					} else { // 6. If GlobalSignOut.RetryAllowed is already false, know that a second request already made, so reset GlobalSignOut.RetryAllowed and break to next user 
						g.RetryAllowed = true
						break
					}
				}
			}
			if err == nil {
				g.ResultsChannel <- *userSignOutRequestData[i].Username
				break
			}
			log.Event(ctx, "Error Cognito AdminUserGlobalSignOut:", log.ERROR, log.Error(err))
			time.Sleep(backoff)
		}
	}
	close(g.ResultsChannel)
}

// generateGlobalSignOutRequest - local routine to generete the global signout request per user
func (api *API) generateGlobalSignOutRequest(user *cognitoidentityprovider.AdminUserGlobalSignOutInput) (*cognitoidentityprovider.AdminUserGlobalSignOutOutput, error) {
	response, error := api.CognitoClient.AdminUserGlobalSignOut(user)
	return response, error
}
