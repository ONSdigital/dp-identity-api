package api

import (
	"regexp"
)

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