package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ONSdigital/dp-identity-api/models"
	"net/http"

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

type baseHandler func(w http.ResponseWriter, r *http.Request, errorList *models.ErrorList)

func (handler baseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var errorList models.ErrorList
	handler(w, r, &errorList)

	if len(errorList.Errors) > 0 {
		ctx := r.Context()
		WriteErrorResponse(ctx, w, errorList)
	}
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
	r.HandleFunc("/tokens/self", baseHandler(api.signOutHandler).ServeHTTP).Methods("DELETE")
	r.HandleFunc("/tokens/self", api.RefreshHandler(ctx)).Methods("PUT")
	r.HandleFunc("/users", api.CreateUserHandler(ctx)).Methods("POST")
	return api, nil
}

func WriteErrorResponse(ctx context.Context, w http.ResponseWriter, errorList models.ErrorList) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorList.Status)

	jsonResponse, err := json.Marshal(errorList)
	if err != nil {
		log.Event(ctx, "failed to marshal the error", log.Error(err), log.ERROR)
		http.Error(w, "failed to marshal the error", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Event(ctx, "writing response failed", log.Error(err), log.ERROR)
		http.Error(w, "failed to write http response", http.StatusInternalServerError)
		return
	}
}
