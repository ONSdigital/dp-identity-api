package api

import (
	"context"

	cognitoclient "github.com/ONSdigital/dp-identity-api/cognitoclient"
	"github.com/gorilla/mux"
)

//API provides a struct to wrap the api around
type API struct {
	Router        *mux.Router
	CognitoClient cognitoclient.Client
}

//Setup function sets up the api and returns an api
func Setup(ctx context.Context, r *mux.Router, cognitoClient cognitoclient.Client) *API {
	api := &API{	
		Router:        r,
		CognitoClient: cognitoClient,
	}

	r.HandleFunc("/hello", HelloHandler(ctx)).Methods("GET")
    return api
}
