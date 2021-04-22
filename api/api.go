package api

import (
	"context"

	"github.com/ONSdigital/dp-identity-api/config"
	"github.com/gorilla/mux"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

//API provides a struct to wrap the api around
type API struct {
	Router        *mux.Router
	CognitoClient *cognito.CognitoIdentityProvider
}

//Setup function sets up the api and returns an api
func Setup(ctx context.Context, cfg *config.Config, r *mux.Router) *API {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	cognitoClient := cognito.New(sess, &aws.Config{Region: &cfg.AWSRegion})

	api := &API{
		Router:        r,
		CognitoClient: cognitoClient,
	}

	r.HandleFunc("/hello", HelloHandler(ctx)).Methods("GET")
	r.HandleFunc("/tokens", TokensHandler()).Methods("POST")
	return api
}
