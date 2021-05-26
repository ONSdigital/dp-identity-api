package api

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ONSdigital/dp-identity-api/apierrorsdeprecated"
	"github.com/ONSdigital/dp-identity-api/cognito"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/ONSdigital/dp-identity-api/utilities"
	"github.com/ONSdigital/dp-identity-api/validation"
	"github.com/ONSdigital/log.go/log"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
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
		return nil, &models.ErrorResponse{
			Errors: []error{validationErr},
			Status: http.StatusBadRequest,
		}
	}

	_, err := api.CognitoClient.GlobalSignOut(accessToken.GenerateSignOutRequest())

	if err != nil {
		responseErr := models.NewCognitoError(ctx, err, "Cognito GlobalSignOut request for signout")
		if responseErr.Code == models.NotFoundError || responseErr.Code == models.NotAuthorisedError {
			return nil, &models.ErrorResponse{
				Errors: []error{responseErr},
				Status: http.StatusBadRequest,
			}
		} else {
			return nil, &models.ErrorResponse{
				Errors: []error{responseErr},
				Status: http.StatusInternalServerError,
			}
		}
	}

	return models.NewSuccessResponse(nil, http.StatusNoContent), nil
}

//RefreshHandler refreshes a users access token and returns new access and ID tokens, expiration time and the refresh token
func (api *API) RefreshHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var errorList []apierrorsdeprecated.Error

		refreshToken := models.RefreshToken{
			TokenString: req.Header.Get(RefreshTokenHeaderName),
		}
		errorList = refreshToken.Validate(ctx, errorList)

		idToken := models.IdToken{
			TokenString: req.Header.Get(IdTokenHeaderName),
		}
		errorList = idToken.Validate(ctx, errorList)

		if len(errorList) > 0 {
			errorResponseBody := apierrorsdeprecated.ErrorResponseBodyBuilder(errorList)
			apierrorsdeprecated.WriteErrorResponse(ctx, w, http.StatusBadRequest, errorResponseBody)
			return
		}

		authInput := refreshToken.GenerateRefreshRequest(api.ClientSecret, idToken.Claims.CognitoUser, api.ClientId)
		result, authErr := api.CognitoClient.InitiateAuth(authInput)

		if authErr != nil {
			log.Event(ctx, "Cognito InitiateAuth request for token refresh - "+authErr.Error(), log.ERROR)
			if authErr.Error() == "NotAuthorizedException: Refresh Token has expired" {
				expiredTokenError := apierrorsdeprecated.IndividualErrorBuilder(authErr, apierrorsdeprecated.TokenExpiredMessage)
				errorList = append(errorList, expiredTokenError)
				errorResponseBody := apierrorsdeprecated.ErrorResponseBodyBuilder(errorList)
				apierrorsdeprecated.WriteErrorResponse(ctx, w, http.StatusForbidden, errorResponseBody)
			} else {
				apierrorsdeprecated.HandleUnexpectedError(ctx, w, authErr, apierrorsdeprecated.InternalErrorMessage)
			}
			return
		}

		buildSuccessfulResponse(result, w, ctx)

		return
	}
}

func terminateExistingSession(authParams AuthParams, userPoolId string, cognitoClient cognito.Client) (err error) {

	adminUserGlobalSignOutInput := &cognitoidentityprovider.AdminUserGlobalSignOutInput{
		Username:   &authParams.Email,
		UserPoolId: &userPoolId,
	}

	_, err = cognitoClient.AdminUserGlobalSignOut(adminUserGlobalSignOutInput)
	if err != nil {
		return err
	}
	return nil
}

func passwordValidation(requestBody AuthParams) (isPasswordValid bool) {

	return len(requestBody.Password) > 0
}

//emailValidation checks for both a valid email address and an empty email address
func emailValidation(requestBody AuthParams) (isEmailValid bool) {

	isEmailValid = validation.IsEmailValid(requestBody.Email)

	return isEmailValid
}

func buildCognitoRequest(authParams AuthParams, clientId string, clientSecret string, clientAuthFlow string) (authInput *cognitoidentityprovider.InitiateAuthInput) {

	secretHash := utilities.ComputeSecretHash(clientSecret, authParams.Email, clientId)

	authParameters := map[string]*string{
		"USERNAME":    &authParams.Email,
		"PASSWORD":    &authParams.Password,
		"SECRET_HASH": &secretHash,
	}

	authInput = &cognitoidentityprovider.InitiateAuthInput{
		AnalyticsMetadata: &cognitoidentityprovider.AnalyticsMetadataType{},
		AuthFlow:          &clientAuthFlow,
		AuthParameters:    authParameters,
		ClientId:          &clientId,
		ClientMetadata:    map[string]*string{},
		UserContextData:   &cognitoidentityprovider.UserContextDataType{},
	}

	return authInput
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
