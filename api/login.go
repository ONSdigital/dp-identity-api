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

var passwordErrorList []IndividualError
var emailErrorList []IndividualError

func LoginHandler(ctx context.Context) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

		body, _ := ioutil.ReadAll(req.Body)

		//Need to close the body

		authParams := make(map[string]string)
		_ = json.Unmarshal(body, &authParams)

		passwordResponse := passwordValidation(authParams)
		if !(passwordResponse) {

			errorMessage := errors.New("Invalid password")
			message := "Unable to validate the password in the request"
			field := ""
			param := ""

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			errorList := errorListBuilder(passwordErrorList, individualErrorBuilder(errorMessage, message, field, param))
			errorResponseBody := errorResponseBodyBuilder(errorList)

			jsonResponse, _ := json.Marshal(errorResponseBody)
			_, _ = w.Write(jsonResponse)

			return
		}

		validEmailResponse := emailValidation(authParams)
		if !validEmailResponse {

			errorMessage := errors.New("Invalid email")
			message := "Unable to validate the email in the request"
			field := ""
			param := ""

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			errorList := errorListBuilder(emailErrorList, individualErrorBuilder(errorMessage, message, field, param))
			errorResponseBody := errorResponseBodyBuilder(errorList)

			jsonResponse, _ := json.Marshal(errorResponseBody)
			_, _ = w.Write(jsonResponse)

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

func errorListBuilder(errorList []IndividualError, individualError IndividualError) (listOfErrors []IndividualError) {

	errorList = append(errorList, individualError)

	return errorList
}

func errorResponseBodyBuilder(listOfErrors []IndividualError) (errorResponseBody ErrorStructure) {

	errorResponseBody = ErrorStructure{
		Errors: listOfErrors,
	}

	return errorResponseBody
}
