package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ONSdigital/dp-identity-api/apierrors"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/ONSdigital/dp-identity-api/validation"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/ONSdigital/log.go/log"
)

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (api *API) TokensHandler(ctx context.Context) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		field := ""
		param := ""

		defer req.Body.Close()

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			errorMessage := "api endpoint POST login returned an error reading the request body"
			apierrors.HandleUnexpectedError(ctx, w, err, errorMessage, field, param)
			return
		}

		var authParams AuthParams
		err = json.Unmarshal(body, &authParams)
		if err != nil {
			errorMessage := "api endpoint POST login returned an error unmarshalling the body"
			apierrors.HandleUnexpectedError(ctx, w, err, errorMessage, field, param)
			return
		}

		validPasswordRequest := passwordValidation(authParams)
		validEmailRequest := emailValidation(authParams)

		invalidPasswordErrorBody := apierrors.IndividualErrorBuilder(apierrors.ErrInvalidPassword, apierrors.InvalidPasswordMessage, field, param)
		invalidEmailErrorBody := apierrors.IndividualErrorBuilder(apierrors.ErrInvalidEmail, apierrors.InvalidErrorMessage, field, param)

		var errorList []models.IndividualError

		if !validPasswordRequest {
			errorList = append(errorList, invalidPasswordErrorBody)
		}

		if !validEmailRequest {
			errorList = append(errorList, invalidEmailErrorBody)
		}

		errorResponseBody := apierrors.ErrorResponseBodyBuilder(errorList)
		apierrors.WriteErrorResponse(ctx, w, http.StatusBadRequest, errorResponseBody)

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
			invalidTokenErrorBody := apierrors.IndividualErrorBuilder(apierrors.InvalidTokenError, apierrors.MissingTokenMessage, field, param)
			errorList = append(errorList, invalidTokenErrorBody)
			log.Event(ctx, "no authorization header provided", log.ERROR)
			errorResponseBody := apierrors.ErrorResponseBodyBuilder(errorList)
			apierrors.WriteErrorResponse(ctx, w, http.StatusBadRequest, errorResponseBody)
			return
		}

		authComponents := strings.Split(authString, " ")
		if len(authComponents) != 2 {
			log.Event(ctx, "malformed authorization header provided", log.ERROR)
			invalidTokenErrorBody := apierrors.IndividualErrorBuilder(apierrors.InvalidTokenError, apierrors.MalformedTokenMessage, field, param)
			errorList = append(errorList, invalidTokenErrorBody)
			errorResponseBody := apierrors.ErrorResponseBodyBuilder(errorList)
			apierrors.WriteErrorResponse(ctx, w, http.StatusBadRequest, errorResponseBody)
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
				apierrors.WriteErrorResponse(ctx, w, http.StatusInternalServerError, errorResponseBody)
			} else {
				apierrors.WriteErrorResponse(ctx, w, http.StatusBadRequest, errorResponseBody)
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