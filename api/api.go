package api

import (
	"context"

	"github.com/ONSdigital/dp-identity-api/cognito"
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
func Setup(ctx context.Context, r *mux.Router, cognitoClient cognito.Client, userPoolId, clientId, clientSecret, clientAuthFlow string) *API {
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
	return api
}
