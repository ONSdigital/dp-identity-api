package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ONSdigital/dp-identity-api/apierrors"
	"github.com/ONSdigital/dp-identity-api/validation"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ONSdigital/log.go/log"
)

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var invalidPasswordError = errors.New("Invalid password")
var invalidPasswordMessage = "Unable to validate the password in the request"
var invalidEmailError = errors.New("Invalid email")
var invalidErrorMessage = "Unable to validate the email in the request"
var invalidTokenError = errors.New("Invalid token")
var invalidTokenMessage = "The provided token does not correspond to an active session"

func TokensHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		field := ""
		param := ""

		defer req.Body.Close()

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			errorMessage := "api endpoint POST login returned an error reading the request body"
			handleUnexpectedError(ctx, w, err, errorMessage, field, param)
			return
		}

		var authParams AuthParams
		err = json.Unmarshal(body, &authParams)
		if err != nil {
			errorMessage := "api endpoint POST login returned an error unmarshalling the body"
			handleUnexpectedError(ctx, w, err, errorMessage, field, param)
			return
		}

		validPasswordRequest := passwordValidation(authParams)
		validEmailRequest := emailValidation(authParams)

		invalidPasswordErrorBody := apierrors.IndividualErrorBuilder(invalidPasswordError, invalidPasswordMessage, field, param)
		invalidEmailErrorBody := apierrors.IndividualErrorBuilder(invalidEmailError, invalidErrorMessage, field, param)

		var errorList []apierrors.IndividualError

		if !validPasswordRequest {
			errorList = append(errorList, invalidPasswordErrorBody)
		}

		if !validEmailRequest {
			errorList = append(errorList, invalidEmailErrorBody)
		}

		errorResponseBody := apierrors.ErrorResponseBodyBuilder(errorList)
		writeErrorResponse(ctx, w, http.StatusBadRequest, errorResponseBody)
		return
	}
}

func (api *API) SignOutHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		var errorList []apierrors.IndividualError
		field := ""
		param := ""

		authString := req.Header.Get("Authorization")
		if authString == "" {
			invalidTokenErrorBody := apierrors.IndividualErrorBuilder(invalidTokenError, invalidTokenMessage, field, param)
			errorList = append(errorList, invalidTokenErrorBody)
			log.Event(ctx, "no authorization header provided", log.ERROR)
			errorResponseBody := apierrors.ErrorResponseBodyBuilder(errorList)
			writeErrorResponse(ctx, w, http.StatusBadRequest, errorResponseBody)
			return
		}

		authComponents := strings.Split(authString, " ")
		if len(authComponents) != 2 {
			log.Event(ctx, "malformed authorization header provided", log.ERROR)
			invalidTokenErrorBody := apierrors.IndividualErrorBuilder(invalidTokenError, invalidTokenMessage, field, param)
			errorList = append(errorList, invalidTokenErrorBody)
			errorResponseBody := apierrors.ErrorResponseBodyBuilder(errorList)
			writeErrorResponse(ctx, w, http.StatusBadRequest, errorResponseBody)
			return
		}

		_, err := api.CognitoClient.GlobalSignOut(
			&cognitoidentityprovider.GlobalSignOutInput{
				AccessToken: &authComponents[1]})

		if err != nil {
			log.Event(ctx, "From Cognito - "+err.Error(), log.ERROR)
			invalidTokenErrorBody := apierrors.IndividualErrorBuilder(err, "", field, param)
			errorList = append(errorList, invalidTokenErrorBody)
			errorResponseBody := apierrors.ErrorResponseBodyBuilder(errorList)
			if strings.Contains(err.Error(), "InternalErrorException") {
				writeErrorResponse(ctx, w, http.StatusInternalServerError, errorResponseBody)
			} else {
				writeErrorResponse(ctx, w, http.StatusBadRequest, errorResponseBody)
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

func writeErrorResponse(ctx context.Context, w http.ResponseWriter, status int, errorResponseBody interface{}) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	jsonResponse, err := json.Marshal(errorResponseBody)
	if err != nil {
		log.Event(ctx, "failed to marshal the error", log.Error(err), log.ERROR)
		http.Error(w, "failed to marshal the error", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Event(ctx, "writing response failed", log.Error(err), log.ERROR)
		http.Error(w, "failed to write http response", http.StatusInternalServerError)
		return
	}
}

func handleUnexpectedError(ctx context.Context, w http.ResponseWriter, err error, message, sourceField, sourceParam string) {

	var errorList []apierrors.IndividualError

	internalServerErrorBody := apierrors.IndividualErrorBuilder(err, message, sourceField, sourceParam)
	errorList = append(errorList, internalServerErrorBody)
	errorResponseBody := apierrors.ErrorResponseBodyBuilder(errorList)

	log.Event(ctx, message, log.ERROR, log.Error(err))
	writeErrorResponse(ctx, w, http.StatusInternalServerError, errorResponseBody)
	return
}
