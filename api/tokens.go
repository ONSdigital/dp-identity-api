package api

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/ONSdigital/dp-identity-api/emailvalidation"

	"github.com/ONSdigital/log.go/log"
)

type ErrorStructure struct {
	Errors []IndividualError `json:"errors"`
}

type IndividualError struct {
	SpecificError string `json:"error"`
	Message       string `json:"message"`
	Source        Source `json:"source"`
}

type Source struct {
	Field string `json:"field"`
	Param string `json:"param"`
}

func TokensHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		field := ""
		param := ""

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			errorMessage := "api endpoint POST login returned an error reading the request body"
			handleUnexpectedError(ctx, w, err, errorMessage, field, param)
			return
		}

		defer req.Body.Close()

		authParams := make(map[string]string)
		err = json.Unmarshal(body, &authParams)
		if err != nil {
			errorMessage := "api endpoint POST login returned an error unmarshalling the body"
			handleUnexpectedError(ctx, w, err, errorMessage, field, param)
			return
		}

		validPasswordRequest := passwordValidation(authParams)
		validEmailRequest := emailValidation(authParams)
		if err != nil {
			errorMessage := "api endpoint POST login returned an error validating the email"
			handleUnexpectedError(ctx, w, err, errorMessage, field, param)
			return
		}

		invalidPasswordError := errors.New("Invalid password")
		invalidPasswordMessage := "Unable to validate the password in the request"
		invalidEmailError := errors.New("Invalid email")
		invalidErrorMessage := "Unable to validate the email in the request"

		invalidPasswordErrorBody := individualErrorBuilder(invalidPasswordError, invalidPasswordMessage, field, param)
		invalidEmailErrorBody := individualErrorBuilder(invalidEmailError, invalidErrorMessage, field, param)

		var errorList []IndividualError

		if !(validPasswordRequest) && !(validEmailRequest) {

			errorList = nil
			errorList = append(errorList, invalidPasswordErrorBody)
			errorList = append(errorList, invalidEmailErrorBody)

			errorResponseBody := errorResponseBodyBuilder(errorList)
			writeErrorResponse(ctx, w, 400, errorResponseBody)
			return
		}

		if !validPasswordRequest {

			errorList = nil
			errorList = append(errorList, invalidPasswordErrorBody)
			errorResponseBody := errorResponseBodyBuilder(errorList)

			writeErrorResponse(ctx, w, 400, errorResponseBody)
			return
		}

		if !validEmailRequest {

			errorList = nil
			errorList = append(errorList, invalidEmailErrorBody)
			errorResponseBody := errorResponseBodyBuilder(errorList)

			writeErrorResponse(ctx, w, 400, errorResponseBody)
			return
		}
	}
}

func passwordValidation(requestBody map[string]string) (isPasswordValid bool) {

	isPasswordValid = false

	if len(requestBody["password"]) != 0 {
		isPasswordValid = true
	}

	return isPasswordValid
}

//emailValidation checks for both a valid email address and an empty email address
func emailValidation(requestBody map[string]string) (isEmailValid bool) {

	isEmailValid = false
	isEmailValid = emailvalidation.IsEmailValid(requestBody["email"])

	return isEmailValid
}

func individualErrorBuilder(err error, message, sourceField, sourceParam string) (individualError IndividualError) {

	individualError = IndividualError{
		SpecificError: error.Error(err),
		Message:       message,
		Source: Source{
			Field: sourceField,
			Param: sourceParam},
	}

	return individualError
}

func errorResponseBodyBuilder(listOfErrors []IndividualError) (errorResponseBody ErrorStructure) {

	errorResponseBody = ErrorStructure{
		Errors: listOfErrors,
	}

	return errorResponseBody
}

func writeErrorResponse(ctx context.Context, w http.ResponseWriter, status int, errorResponseBody interface{}) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	jsonResponse, _ := json.Marshal(errorResponseBody)
	_, err := w.Write(jsonResponse)
	if err != nil {

		log.Event(ctx, "writing response failed", log.Error(err), log.ERROR)
		http.Error(w, "Failed to write http response", http.StatusInternalServerError)
		return
	}

	return
}

func handleUnexpectedError(ctx context.Context, w http.ResponseWriter, err error, message, sourceField, sourceParam string) {

	var errorList []IndividualError

	errorList = nil
	statusCode := 500
	internalServerErrorBody := individualErrorBuilder(err, message, sourceField, sourceParam)
	errorList = append(errorList, internalServerErrorBody)
	errorResponseBody := errorResponseBodyBuilder(errorList)

	log.Event(ctx, message, log.ERROR, log.Error(err))
	writeErrorResponse(ctx, w, statusCode, errorResponseBody)
	return
}
