package api

import (
	"context"
	"regexp"
	"net/http"
	"fmt"
)

func LoginHandler(ctx context.Context) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request){
		req.ParseForm()

		userEmail := req.Form.Get("email")
		password := req.Form.Get("password")

		authParams := map[string]string{
			"email": userEmail,
			"password": password,
		}

		fmt.Println(authParams)

		emailResponse, passwordResponse := bodyValidation(authParams) 
		if emailResponse || passwordResponse == false {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	}

}

func bodyValidation(requestBody map[string]string)(emailResponse, passwordResponse bool){

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

func emailValidation(requestBody map[string]string)(emailResponse bool){

	emailResponse = false

	emailResponse, _ = regexp.MatchString("^[a-zA-Z0-9.]+@(ext.)?ons.gov.uk$", requestBody["email"])

	return emailResponse
}