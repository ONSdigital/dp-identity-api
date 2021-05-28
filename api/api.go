package api

import (
	"context"
	"encoding/json"
	"github.com/ONSdigital/dp-identity-api/models"
	"net/http"

	"github.com/ONSdigital/dp-identity-api/cognito"
	"github.com/gorilla/mux"
)

var (
	IdTokenHeaderName      = "ID"
	AccessTokenHeaderName  = "Authorization"
	RefreshTokenHeaderName = "Refresh"
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

type baseHandler func(w http.ResponseWriter, r *http.Request, ctx context.Context) (*models.SuccessResponse, *models.ErrorResponse)

func contextAndErrors(h baseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		response, err := h(w, req, ctx)
		if err != nil {
			writeErrorResponse(ctx, w, err)
			return
		}
		writeSuccessResponse(ctx, w, response)
	}
}

//Setup function sets up the api and returns an api
func Setup(ctx context.Context, r *mux.Router, cognitoClient cognito.Client, userPoolId string, clientId string, clientSecret string, clientAuthFlow string) (*API, error) {

	// Return an error if empty required parameter was passed.
	if userPoolId == "" || clientId == "" || clientSecret == "" || clientAuthFlow == "" {
		return nil, models.NewError(ctx, nil, models.MissingConfigError, models.MissingConfigDescription)
	}

	api := &API{
		Router:         r,
		CognitoClient:  cognitoClient,
		UserPoolId:     userPoolId,
		ClientId:       clientId,
		ClientSecret:   clientSecret,
		ClientAuthFlow: clientAuthFlow,
	}

	r.HandleFunc("/tokens", contextAndErrors(api.TokensHandler)).Methods("POST")
	// self used in paths rather than identifier as the identifier is JWT tokens passed in the request headers
	r.HandleFunc("/tokens/self", contextAndErrors(api.SignOutHandler)).Methods("DELETE")
	r.HandleFunc("/tokens/self", contextAndErrors(api.RefreshHandler)).Methods("PUT")
	r.HandleFunc("/users", contextAndErrors(api.CreateUserHandler)).Methods("POST")
	return api, nil
}

func writeErrorResponse(ctx context.Context, w http.ResponseWriter, errorResponse *models.ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorResponse.Status)

	jsonResponse, err := json.Marshal(errorResponse)
	if err != nil {
		responseErr := models.NewError(ctx, err, models.JSONMarshalError, models.ErrorMarshalFailedDescription)
		http.Error(w, responseErr.Description, http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		responseErr := models.NewError(ctx, err, models.WriteResponseError, models.WriteResponseFailedDescription)
		http.Error(w, responseErr.Description, http.StatusInternalServerError)
		return
	}
}

func writeSuccessResponse(ctx context.Context, w http.ResponseWriter, successResponse *models.SuccessResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(successResponse.Status)

	_, err := w.Write(successResponse.Body)
	if err != nil {
		responseErr := models.NewError(ctx, err, models.WriteResponseError, models.WriteResponseFailedDescription)
		http.Error(w, responseErr.Description, http.StatusInternalServerError)
		return
	}
}

func handleBodyReadError(ctx context.Context, err error) *models.ErrorResponse {
	return models.NewErrorResponse([]error{models.NewError(ctx, err, models.BodyReadError, models.BodyReadFailedDescription)},
		http.StatusInternalServerError)
}

func handleBodyUnmarshalError(ctx context.Context, err error) *models.ErrorResponse {
	return models.NewErrorResponse([]error{models.NewError(ctx, err, models.JSONUnmarshalError, models.ErrorUnmarshalFailedDescription)},
		http.StatusInternalServerError)
}
