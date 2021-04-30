package api

import (
	"context"

	"github.com/ONSdigital/dp-identity-api/cognito"
	"github.com/gorilla/mux"
)

//API provides a struct to wrap the api around
type API struct {
	Router        *mux.Router
	CognitoClient cognito.Client
	UserPoolId	string
}

//Setup function sets up the api and returns an api
func Setup(ctx context.Context, r *mux.Router, cognitoClient cognito.Client, userPoolId string) *API {
	api := &API{	
		Router:        r,
		CognitoClient: cognitoClient,
		UserPoolId: userPoolId,
	}

	r.HandleFunc("/hello", HelloHandler(ctx)).Methods("GET")
    return api
}
