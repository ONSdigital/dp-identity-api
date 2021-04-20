package api

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
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

var errorList []IndividualError

func LoginHandler(ctx context.Context) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

		body, _ := ioutil.ReadAll(req.Body)

		defer req.Body.Close()

		authParams := make(map[string]string)
		_ = json.Unmarshal(body, &authParams)

		validPasswordRequest := passwordValidation(authParams)
		validEmailRequest := emailValidation(authParams)

		invalidPasswordError := errors.New("Invalid password")
		invalidPasswordMessage := "Unable to validate the password in the request"
		invalidEmailError := errors.New("Invalid email")
		invalidErrorMessage := "Unable to validate the email in the request"
		field := ""
		param := ""

		invalidPasswordErrorBody := individualErrorBuilder(invalidPasswordError, invalidPasswordMessage, field, param)
		invalidEmailErrorBody := individualErrorBuilder(invalidEmailError, invalidErrorMessage, field, param)

		if !(validPasswordRequest) && !(validEmailRequest) {

			errorList := append(errorList, invalidPasswordErrorBody)
			errorList = append(errorList, invalidEmailErrorBody)

			errorResponseBody := errorResponseBodyBuilder(errorList)
			writeErrorResponse(w, 400, errorResponseBody)
			return
		}

		if !validPasswordRequest {

			errorList := append(errorList, invalidPasswordErrorBody)
			errorResponseBody := errorResponseBodyBuilder(errorList)

			writeErrorResponse(w, 400, errorResponseBody)
			return
		}

		if !validEmailRequest {

			errorList := append(errorList, invalidEmailErrorBody)
			errorResponseBody := errorResponseBodyBuilder(errorList)

			writeErrorResponse(w, 400, errorResponseBody)
			return
		}
	}
}

func passwordValidation(requestBody map[string]string) (passwordResponse bool) {

	passwordResponse = false

	if len(requestBody["password"]) != 0 {
		passwordResponse = true
	}

	return passwordResponse
}

//emailValidation checks for both a valid email address and an empty email address
func emailValidation(requestBody map[string]string) (emailResponse bool) {

	emailResponse = false

	emailResponse, _ = regexp.MatchString("^[a-zA-Z0-9.]+@(ext.)?ons.gov.uk$", requestBody["email"])

	return emailResponse
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

func writeErrorResponse(w http.ResponseWriter, status int, errorResponseBody interface{}) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	jsonResponse, _ := json.Marshal(errorResponseBody)
	_, _ = w.Write(jsonResponse)

	return
}
