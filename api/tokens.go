package api

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/ONSdigital/dp-identity-api/utilities"
	"github.com/ONSdigital/dp-identity-api/validation"
	"github.com/ONSdigital/log.go/log"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/ONSdigital/dp-identity-api/apierrorsdeprecated"
)

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (api *API) TokensHandler(ctx context.Context) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

		field := ""
		param := ""

		defer req.Body.Close()

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			errorMessage := "api endpoint POST login returned an error reading the request body"
			apierrorsdeprecated.HandleUnexpectedError(ctx, w, err, errorMessage, field, param)
			return
		}

		var authParams AuthParams
		err = json.Unmarshal(body, &authParams)
		if err != nil {
			errorMessage := "api endpoint POST login returned an error unmarshalling the body"
			apierrorsdeprecated.HandleUnexpectedError(ctx, w, err, errorMessage, field, param)
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
					errorMessage := "api endpoint POST login returned an error and failed to login to cognito"
					apierrorsdeprecated.HandleUnexpectedError(ctx, w, authErr, errorMessage, field, param)
					return
				}

				var errorList []models.IndividualError
				switch authErr.Error() {
				case "NotAuthorizedException: Incorrect username or password.":
					{
						notAuthorizedMessage := "unautheticated user: Unable to autheticate request"
						notAuthorizedError := apierrorsdeprecated.IndividualErrorBuilder(authErr, notAuthorizedMessage, field, param)
						errorList = append(errorList, notAuthorizedError)
						errorResponseBody := apierrorsdeprecated.ErrorResponseBodyBuilder(errorList)
						apierrorsdeprecated.WriteErrorResponse(ctx, w, http.StatusUnauthorized, errorResponseBody)
						return
					}
				case "NotAuthorizedException: Password attempts exceeded":
					{
						forbiddenMessage := "exceeded the number of attemps to login in with the provided credentials"
						forbiddenError := apierrorsdeprecated.IndividualErrorBuilder(authErr, forbiddenMessage, field, param)
						errorList = append(errorList, forbiddenError)
						errorResponseBody := apierrorsdeprecated.ErrorResponseBodyBuilder(errorList)
						apierrorsdeprecated.WriteErrorResponse(ctx, w, http.StatusForbidden, errorResponseBody)
						return
					}
				default:
					{
						loginMessage := "something went wrong, and api endpoint POST login returned an error and failed to login to cognito. Please try again or contact an administrator."
						loginError := apierrorsdeprecated.IndividualErrorBuilder(authErr, loginMessage, field, param)
						errorList = append(errorList, loginError)
						errorResponseBody := apierrorsdeprecated.ErrorResponseBodyBuilder(errorList)
						apierrorsdeprecated.WriteErrorResponse(ctx, w, http.StatusBadRequest, errorResponseBody)
						return
					}
				}
			}

			buildSucessfulResponse(result, w, ctx)

			return
		}

		invalidPasswordErrorBody := apierrorsdeprecated.IndividualErrorBuilder(apierrorsdeprecated.ErrInvalidPassword, apierrorsdeprecated.InvalidPasswordMessage, field, param)
		invalidEmailErrorBody := apierrorsdeprecated.IndividualErrorBuilder(apierrorsdeprecated.ErrInvalidEmail, apierrorsdeprecated.InvalidErrorMessage, field, param)

		var errorList []models.IndividualError

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

func (api *API) SignOutHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		var errorList []models.IndividualError
		field := ""
		param := ""

		authString := req.Header.Get("Authorization")
		if authString == "" {
			invalidTokenErrorBody := apierrorsdeprecated.IndividualErrorBuilder(apierrorsdeprecated.InvalidTokenError, apierrorsdeprecated.MissingTokenMessage, field, param)
			errorList = append(errorList, invalidTokenErrorBody)
			log.Event(ctx, "no authorization header provided", log.ERROR)
			errorResponseBody := apierrorsdeprecated.ErrorResponseBodyBuilder(errorList)
			apierrorsdeprecated.WriteErrorResponse(ctx, w, http.StatusBadRequest, errorResponseBody)
			return
		}

		authComponents := strings.Split(authString, " ")
		if len(authComponents) != 2 {
			log.Event(ctx, "malformed authorization header provided", log.ERROR)
			invalidTokenErrorBody := apierrorsdeprecated.IndividualErrorBuilder(apierrorsdeprecated.InvalidTokenError, apierrorsdeprecated.MalformedTokenMessage, field, param)
			errorList = append(errorList, invalidTokenErrorBody)
			errorResponseBody := apierrorsdeprecated.ErrorResponseBodyBuilder(errorList)
			apierrorsdeprecated.WriteErrorResponse(ctx, w, http.StatusBadRequest, errorResponseBody)
			return
		}

		_, err := api.CognitoClient.GlobalSignOut(
			&cognitoidentityprovider.GlobalSignOutInput{
				AccessToken: &authComponents[1]})

		if err != nil {
			log.Event(ctx, "From Cognito - "+err.Error(), log.ERROR)
			invalidTokenErrorBody := apierrorsdeprecated.IndividualErrorBuilder(err, "", field, param)
			errorList = append(errorList, invalidTokenErrorBody)
			errorResponseBody := apierrorsdeprecated.ErrorResponseBodyBuilder(errorList)
			isInternalError := apierrorsdeprecated.IdentifyInternalError(err)
			if isInternalError {
				apierrorsdeprecated.WriteErrorResponse(ctx, w, http.StatusInternalServerError, errorResponseBody)
			} else {
				apierrorsdeprecated.WriteErrorResponse(ctx, w, http.StatusBadRequest, errorResponseBody)
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
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

func buildSucessfulResponse(result *cognitoidentityprovider.InitiateAuthOutput, w http.ResponseWriter, ctx context.Context) {

	if result.AuthenticationResult != nil {
		tokenDuration := time.Duration(*result.AuthenticationResult.ExpiresIn)
		expirationTime := time.Now().UTC().Add(time.Second * tokenDuration).String()

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Authorization", "Bearer "+*result.AuthenticationResult.AccessToken)
		w.Header().Set("ID", *result.AuthenticationResult.IdToken)
		w.Header().Set("Refresh", *result.AuthenticationResult.RefreshToken)
		w.WriteHeader(http.StatusCreated)

		postBody := map[string]interface{}{"expirationTime": expirationTime}

		buildjson(postBody, w, ctx)

		return
	} else {
		err := errors.New("unexpected return from cognito")
		errorMessage := "unexpected response from cognito, there was no authentication result field"
		apierrorsdeprecated.HandleUnexpectedError(ctx, w, err, errorMessage, "", "")
		return
	}
}

func buildjson(jsonInput map[string]interface{}, w http.ResponseWriter, ctx context.Context) {

	jsonResponse, err := json.Marshal(jsonInput)

	if err != nil {
		errorMessage := "failed to marshal the error"
		apierrorsdeprecated.HandleUnexpectedError(ctx, w, err, errorMessage, "", "")
		return
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		errorMessage := "writing response failed"
		apierrorsdeprecated.HandleUnexpectedError(ctx, w, err, errorMessage, "", "")

		return
	}
	return

}
