package api

import (
	"context"
	"os"

	"github.com/ONSdigital/dp-identity-api/config"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

//API provides a struct to wrap the api around
type API struct {
	Router *mux.Router
	CognitoClient *cognito.CognitoIdentityProvider
}

//Setup function sets up the api and returns an api
func Setup(ctx context.Context, r *mux.Router) *API {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	cfg, err := config.Get()
	if err != nil {
		log.Event(ctx, "Failed to initialise config", log.FATAL, log.Error(err))
		os.Exit(1)
	}

	cognitoClient := cognito.New(sess, &aws.Config{Region: &cfg.AWSRegion})

	api := &API{
		Router: r,
		CognitoClient: cognitoClient,
	}

	r.HandleFunc("/hello", HelloHandler(ctx)).Methods("GET")
	return api
}
