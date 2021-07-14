package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ONSdigital/dp-identity-api/models"
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
