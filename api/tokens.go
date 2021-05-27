package api

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"time"

	"errors"
	"io/ioutil"
	"net/http"

	"github.com/ONSdigital/dp-identity-api/apierrorsdeprecated"
	"github.com/ONSdigital/dp-identity-api/models"
)

//TokensHandler uses submitted email address and password to sign a user in against Cognito and returns a http handler interface
func (api *API) TokensHandler(w http.ResponseWriter, req *http.Request, ctx context.Context) (*models.SuccessResponse, *models.ErrorResponse) {
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
		return nil, models.NewErrorResponse(*validationErrs, http.StatusBadRequest)
	}

	terminationRequest := userSignIn.BuildOldSessionTerminationRequest(api.UserPoolId)
	_, err = api.CognitoClient.AdminUserGlobalSignOut(terminationRequest)

	if err != nil {
		responseErr := models.NewCognitoError(ctx, err, "Cognito AdminUserGlobalSignOut request from sign in handler")
		if responseErr.Code == models.InternalError {
			return nil, models.NewErrorResponse([]error{responseErr}, http.StatusInternalServerError)
		}

		return nil, models.NewErrorResponse([]error{responseErr}, http.StatusBadRequest)
	}

	input := userSignIn.BuildCognitoRequest(api.ClientId, api.ClientSecret, api.ClientAuthFlow)
	result, authErr := api.CognitoClient.InitiateAuth(input)

	if authErr != nil {
		responseErr := models.NewCognitoError(ctx, authErr, "Cognito InitiateAuth request from sign in handler")
		if responseErr.Code == models.InternalError {
			return nil, models.NewErrorResponse([]error{responseErr}, http.StatusInternalServerError)
		}

		switch responseErr.Description {
		case models.SignInFailedDescription:
			return nil, models.NewErrorResponse([]error{responseErr}, http.StatusUnauthorized)
		case models.SignInAttemptsExceededDescription:
			// Cognito returns the same Code for invalid credentials and too many attempts errors, changing our Error.Code to enable differentiation in the client
			responseErr.Code = models.TooManyFailedAttemptsError
			return nil, models.NewErrorResponse([]error{responseErr}, http.StatusForbidden)
		default:
			return nil, models.NewErrorResponse([]error{responseErr}, http.StatusBadRequest)
		}
	}

	jsonResponse, responseErr := userSignIn.BuildSuccessfulJsonResponse(ctx, result)
	if responseErr != nil {
		return nil, models.NewErrorResponse([]error{responseErr}, http.StatusInternalServerError)
	}

	w.Header().Set(AccessTokenHeaderName, "Bearer "+*result.AuthenticationResult.AccessToken)
	w.Header().Set(IdTokenHeaderName, *result.AuthenticationResult.IdToken)
	w.Header().Set(RefreshTokenHeaderName, *result.AuthenticationResult.RefreshToken)
	return models.NewSuccessResponse(jsonResponse, http.StatusCreated), nil
}

//SignOutHandler invalidates a users access token signing them out and returns a http handler interface
func (api *API) SignOutHandler(w http.ResponseWriter, req *http.Request, ctx context.Context) (*models.SuccessResponse, *models.ErrorResponse) {
	accessToken := models.AccessToken{
		AuthHeader: req.Header.Get(AccessTokenHeaderName),
	}
	validationErr := accessToken.Validate(ctx)
	if validationErr != nil {
		return nil, models.NewErrorResponse([]error{validationErr}, http.StatusBadRequest)
	}

	_, err := api.CognitoClient.GlobalSignOut(accessToken.GenerateSignOutRequest())

	if err != nil {
		responseErr := models.NewCognitoError(ctx, err, "Cognito GlobalSignOut request for sign out")
		if responseErr.Code == models.NotFoundError || responseErr.Code == models.NotAuthorisedError {
			return nil, models.NewErrorResponse([]error{responseErr}, http.StatusBadRequest)
		} else {
			return nil, models.NewErrorResponse([]error{responseErr}, http.StatusInternalServerError)
		}
	}

	return models.NewSuccessResponse(nil, http.StatusNoContent), nil
}

//RefreshHandler refreshes a users access token and returns new access and ID tokens, expiration time and the refresh token
func (api *API) RefreshHandler(w http.ResponseWriter, req *http.Request, ctx context.Context) (*models.SuccessResponse, *models.ErrorResponse) {
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
		return nil, models.NewErrorResponse(validationErrs, http.StatusBadRequest)
	}

	authInput := refreshToken.GenerateRefreshRequest(api.ClientSecret, idToken.Claims.CognitoUser, api.ClientId)
	result, authErr := api.CognitoClient.InitiateAuth(authInput)

	if authErr != nil {
		responseErr := models.NewCognitoError(ctx, authErr, "Cognito InitiateAuth request for token refresh")
		if responseErr.Code == models.NotAuthorisedError {
			return nil, models.NewErrorResponse([]error{responseErr}, http.StatusForbidden)
		} else {
			return nil, models.NewErrorResponse([]error{responseErr}, http.StatusInternalServerError)
		}
	}

	jsonResponse, responseErr := refreshToken.BuildSuccessfulJsonResponse(ctx, result)
	if responseErr != nil {
		return nil, models.NewErrorResponse([]error{responseErr}, http.StatusInternalServerError)
	}

	w.Header().Set(AccessTokenHeaderName, "Bearer "+*result.AuthenticationResult.AccessToken)
	w.Header().Set(IdTokenHeaderName, *result.AuthenticationResult.IdToken)
	return models.NewSuccessResponse(jsonResponse, http.StatusCreated), nil
}

func buildSuccessfulResponse(result *cognitoidentityprovider.InitiateAuthOutput, w http.ResponseWriter, ctx context.Context) {

	if result.AuthenticationResult != nil {
		tokenDuration := time.Duration(*result.AuthenticationResult.ExpiresIn)
		expirationTime := time.Now().UTC().Add(time.Second * tokenDuration).String()

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set(AccessTokenHeaderName, "Bearer "+*result.AuthenticationResult.AccessToken)
		w.Header().Set(IdTokenHeaderName, *result.AuthenticationResult.IdToken)
		if result.AuthenticationResult.RefreshToken != nil {
			w.Header().Set(RefreshTokenHeaderName, *result.AuthenticationResult.RefreshToken)
		}
		w.WriteHeader(http.StatusCreated)

		postBody := map[string]interface{}{"expirationTime": expirationTime}

		buildjson(postBody, w, ctx)

		return
	} else {
		err := errors.New("unexpected return from cognito")
		errorDescription := "unexpected response from cognito, there was no authentication result field"
		apierrorsdeprecated.HandleUnexpectedError(ctx, w, err, errorDescription)
		return
	}
}

func buildjson(jsonInput map[string]interface{}, w http.ResponseWriter, ctx context.Context) {

	jsonResponse, err := json.Marshal(jsonInput)

	if err != nil {
		errorDescription := "failed to marshal the error"
		apierrorsdeprecated.HandleUnexpectedError(ctx, w, err, errorDescription)
		return
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		errorDescription := "writing response failed"
		apierrorsdeprecated.HandleUnexpectedError(ctx, w, err, errorDescription)

		return
	}
	return

}
