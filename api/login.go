package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
)

type ErrorStructure struct {
	Errors []IndividualError `json:"errors,omitempty"`
}

type IndividualError struct {
	SpecificError string `json:"error,omitempty"`
	Message       string `json:"message,omitempty"`
	Source        Source `json:"source,omitempty"`
}

type Source struct {
	Field string `json:"field,omitempty"`
	Param string `json:"param,omitempty"`
}

func LoginHandler(ctx context.Context) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

		body, _ := ioutil.ReadAll(req.Body)

		//Need to close the body

		authParams := make(map[string]string)
		_ = json.Unmarshal(body, &authParams)

		emailResponse, passwordResponse := bodyValidation(authParams)
		if !(emailResponse || passwordResponse) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			fullResponse := ErrorStructure{
				Errors: []IndividualError{
					{
						SpecificError: "string, unchanging so devs can use this in code",
						Message:       "detailed explanation of error",
						Source: Source{
							Field: "reference to field like some.field or something",
							Param: "query param causing issue"}},
				},
			}

			jsonResponse, _ := json.Marshal(fullResponse)
			_, _ = w.Write(jsonResponse)

			return
		}

		validEmailResponse := emailValidation(authParams)
		if !validEmailResponse {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

}

func bodyValidation(requestBody map[string]string) (emailResponse, passwordResponse bool) {

	emailResponse = false
	passwordResponse = false

	if len(requestBody["email"]) != 0 {
		emailResponse = true
	}

	if len(requestBody["password"]) != 0 {
		passwordResponse = true
	}

	return emailResponse, passwordResponse
}

func emailValidation(requestBody map[string]string) (emailResponse bool) {

	emailResponse = false

	emailResponse, _ = regexp.MatchString("^[a-zA-Z0-9.]+@(ext.)?ons.gov.uk$", requestBody["email"])

	return emailResponse
}
