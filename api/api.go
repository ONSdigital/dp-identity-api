package api

import (
	"context"
	"errors"

	"github.com/ONSdigital/dp-identity-api/apierrorsdeprecated"
	"github.com/ONSdigital/dp-identity-api/cognito"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

//API provides a struct to wrap the api around
type API struct {
	Router         *mux.Router
	CognitoClient  cognito.Client
	UserPoolId     string
	ClientId       string
	ClientSecret   string
	ClientAuthFlow string
}

//Setup function sets up the api and returns an api
func Setup(ctx context.Context, r *mux.Router, cognitoClient cognito.Client, userPoolId string, clientId string, clientSecret string, clientAuthFlow string) (*API, error) {

	// Return an error if empty required parameter was passed.
	if userPoolId == "" || clientId == "" || clientSecret == "" || clientAuthFlow == "" {
		log.Event(ctx, apierrorsdeprecated.RequiredParameterNotFoundDescription, log.ERROR)
		return nil, errors.New(apierrorsdeprecated.RequiredParameterNotFoundDescription)
	}

	api := &API{
		Router:         r,
		CognitoClient:  cognitoClient,
		UserPoolId:     userPoolId,
		ClientId:       clientId,
		ClientSecret:   clientSecret,
		ClientAuthFlow: clientAuthFlow,
	}

	r.HandleFunc("/tokens", api.TokensHandler(ctx)).Methods("POST")
	// self used in paths rather than identifier as the identifier is JWT tokens passed in the request headers
	r.HandleFunc("/tokens/self", api.SignOutHandler(ctx)).Methods("DELETE")
	r.HandleFunc("/tokens/self", api.RefreshHandler(ctx)).Methods("PUT")
	r.HandleFunc("/users", api.CreateUserHandler(ctx)).Methods("POST")
	return api, nil
}
