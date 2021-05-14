package api

import (
	"context"
	"errors"

	"github.com/ONSdigital/dp-identity-api/api/content"
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

	// Return an error if empty userPoolId was passed.
	if userPoolId == "" {
		log.Event(ctx, content.UserPoolIdNotFoundMessage, log.ERROR)
		return nil, errors.New(content.UserPoolIdNotFoundMessage)
	}
	
	api := &API{
		Router:         r,
		CognitoClient:  cognitoClient,
		UserPoolId:     userPoolId,
		ClientId:       clientId,
		ClientSecret:   clientSecret,
		ClientAuthFlow: clientAuthFlow,
	}

	r.HandleFunc("/hello", HelloHandler(ctx)).Methods("GET")
	r.HandleFunc("/tokens", api.TokensHandler(ctx)).Methods("POST")
	r.HandleFunc("/tokens/self", api.SignOutHandler(ctx)).Methods("DELETE")
	r.HandleFunc("/users", api.CreateUserHandler(ctx)).Methods("POST")
	return api, nil
}
