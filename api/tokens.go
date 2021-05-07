package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/ONSdigital/dp-identity-api/apierrors"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/ONSdigital/dp-identity-api/validation"
)

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var errInvalidPassword = errors.New("invalid password")
var invalidPasswordMessage = "Unable to validate the password in the request"
var errInvalidEmail = errors.New("invalid email")
var invalidErrorMessage = "Unable to validate the email in the request"

func TokensHandler() http.HandlerFunc {

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

		invalidPasswordErrorBody := apierrors.IndividualErrorBuilder(errInvalidPassword, invalidPasswordMessage, field, param)
		invalidEmailErrorBody := apierrors.IndividualErrorBuilder(errInvalidEmail, invalidErrorMessage, field, param)

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

func passwordValidation(requestBody AuthParams) (isPasswordValid bool) {

	return len(requestBody.Password) > 0
}

//emailValidation checks for both a valid email address and an empty email address
func emailValidation(requestBody AuthParams) (isEmailValid bool) {

	isEmailValid = validation.IsEmailValid(requestBody.Email)

	return isEmailValid
}
