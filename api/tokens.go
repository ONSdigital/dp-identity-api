package api

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ONSdigital/dp-identity-api/apierrorsdeprecated"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/ONSdigital/dp-identity-api/utilities"
	"github.com/ONSdigital/dp-identity-api/validation"
	"github.com/ONSdigital/log.go/log"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var (
	IdTokenHeaderName      = "ID"
	AccessTokenHeaderName  = "Authorization"
	RefreshTokenHeaderName = "Refresh"
)

//TokensHandler uses submitted email address and password to sign a user in against Cognito and returns a http handler interface
func (api *API) TokensHandler(ctx context.Context) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			errorDescription := "api endpoint POST login returned an error reading the request body"
			apierrorsdeprecated.HandleUnexpectedError(ctx, w, err, errorDescription)
			return
		}

		var authParams AuthParams
		err = json.Unmarshal(body, &authParams)
		if err != nil {
			errorDescription := "api endpoint POST login returned an error unmarshalling the body"
			apierrorsdeprecated.HandleUnexpectedError(ctx, w, err, errorDescription)
			return
		}

		validPasswordRequest := passwordValidation(authParams)
		validEmailRequest := emailValidation(authParams)

		if validPasswordRequest && validEmailRequest {
			input := buildCognitoRequest(authParams, api.ClientId, api.ClientSecret, api.ClientAuthFlow)
			result, authErr := api.CognitoClient.InitiateAuth(input)

			if authErr != nil {

				isInternalError := apierrorsdeprecated.IdentifyInternalError(authErr)

				if isInternalError {
					errorDescription := "api endpoint POST login returned an error and failed to login to cognito"
					apierrorsdeprecated.HandleUnexpectedError(ctx, w, authErr, errorDescription)
					return
				}

				var errorList []apierrorsdeprecated.Error

				switch authErr.Error() {
				case "NotAuthorizedException: Incorrect username or password.":
					{
						notAuthorizedDescription := "unautheticated user: Unable to autheticate request"
						notAuthorizedError := apierrorsdeprecated.IndividualErrorBuilder(authErr, notAuthorizedDescription)
						errorList = append(errorList, notAuthorizedError)
						errorResponseBody := apierrorsdeprecated.ErrorResponseBodyBuilder(errorList)
						apierrorsdeprecated.WriteErrorResponse(ctx, w, http.StatusUnauthorized, errorResponseBody)
						return
					}
				case "NotAuthorizedException: Password attempts exceeded":
					{
						forbiddenDescription := "exceeded the number of attemps to login in with the provided credentials"
						forbiddenError := apierrorsdeprecated.IndividualErrorBuilder(authErr, forbiddenDescription)
						errorList = append(errorList, forbiddenError)
						errorResponseBody := apierrorsdeprecated.ErrorResponseBodyBuilder(errorList)
						apierrorsdeprecated.WriteErrorResponse(ctx, w, http.StatusForbidden, errorResponseBody)
						return
					}
				default:
					{
						loginDescription := "something went wrong, and api endpoint POST login returned an error and failed to login to cognito. Please try again or contact an administrator."
						loginError := apierrorsdeprecated.IndividualErrorBuilder(authErr, loginDescription)
						errorList = append(errorList, loginError)
						errorResponseBody := apierrorsdeprecated.ErrorResponseBodyBuilder(errorList)
						apierrorsdeprecated.WriteErrorResponse(ctx, w, http.StatusBadRequest, errorResponseBody)
						return
					}
				}
			}

			buildSuccessfulResponse(result, w, ctx)

			return
		}

		invalidPasswordErrorBody := apierrorsdeprecated.IndividualErrorBuilder(apierrorsdeprecated.ErrInvalidPassword, apierrorsdeprecated.InvalidPasswordDescription)
		invalidEmailErrorBody := apierrorsdeprecated.IndividualErrorBuilder(apierrorsdeprecated.ErrInvalidEmail, apierrorsdeprecated.InvalidErrorDescription)

		var errorList []apierrorsdeprecated.Error

		if !validPasswordRequest {
			errorList = append(errorList, invalidPasswordErrorBody)
		}

		if !validEmailRequest {
			errorList = append(errorList, invalidEmailErrorBody)
		}

		errorResponseBody := apierrorsdeprecated.ErrorResponseBodyBuilder(errorList)
		apierrorsdeprecated.WriteErrorResponse(ctx, w, http.StatusBadRequest, errorResponseBody)

	}
}

//SignOutHandler invalidates a users access token signing them out and returns a http handler interface
func (api *API) SignOutHandler(w http.ResponseWriter, req *http.Request, ctx context.Context, errorResponse *models.ErrorResponse) {
	accessToken := models.AccessToken{
		AuthHeader: req.Header.Get(AccessTokenHeaderName),
	}
	validationErr := accessToken.Validate(ctx)
	if validationErr != nil {
		errorResponse.Errors = append(errorResponse.Errors, validationErr)
	}

	if len(errorResponse.Errors) > 0 {
		errorResponse.Status = http.StatusBadRequest
		return
	}

	_, err := api.CognitoClient.GlobalSignOut(accessToken.GenerateSignOutRequest())

	if err != nil {
		responseErr := models.NewCognitoError(ctx, err, "Cognito GlobalSignOut request for signout")
		errorResponse.Errors = append(errorResponse.Errors, &responseErr)
		if responseErr.Code == models.NotFoundError || responseErr.Code == models.NotAuthorisedError {
			errorResponse.Status = http.StatusBadRequest
		} else {
			errorResponse.Status = http.StatusInternalServerError
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
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
