package api

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/cognito"
	"github.com/gorilla/mux"
)

//API provides a struct to wrap the api around
type API struct {
	Router        	*mux.Router
	CognitoClient 	cognito.Client
	UserPoolId		string
	ClientId	 	string
	ClientSecret 	string
}

//Setup function sets up the api and returns an api
func Setup(ctx context.Context, r *mux.Router, cognitoClient cognito.Client, userPoolId string, clientId string, clientSecret string) *API {
	api := &API{
		Router:        	r,
		CognitoClient: 	cognitoClient,
		UserPoolId: 	userPoolId,
		ClientId: 		clientId,
		ClientSecret: 	clientSecret,
	}

	r.HandleFunc("/hello", HelloHandler(ctx)).Methods("GET")
	r.HandleFunc("/tokens", TokensHandler()).Methods("POST")
	r.HandleFunc("/tokens/self", api.SignOutHandler(ctx)).Methods("DELETE")
	return api
}
