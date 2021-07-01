package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/dp-identity-api/models"

	"github.com/ONSdigital/dp-identity-api/cognito"
	"github.com/gorilla/mux"
)

var (
	IdTokenHeaderName      = "ID"
	AccessTokenHeaderName  = "Authorization"
	RefreshTokenHeaderName = "Refresh"
	WWWAuthenticateName    = "WWW-Authenticate"
	ONSRealm               = "Florence publishing platform"
	Charset                = "UTF-8"
	NewPasswordChallenge   = "NEW_PASSWORD_REQUIRED"
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

type baseHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request) (*models.SuccessResponse, *models.ErrorResponse)

func contextAndErrors(h baseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		response, err := h(ctx, w, req)
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

	if err := initialiseRoleGroups(ctx, cognitoClient, userPoolId); err != nil {
		return nil, err
	}

	api := &API{
		Router:         r,
		CognitoClient:  cognitoClient,
		UserPoolId:     userPoolId,
		ClientId:       clientId,
		ClientSecret:   clientSecret,
		ClientAuthFlow: clientAuthFlow,
	}

	r.HandleFunc("/v1/tokens", contextAndErrors(api.TokensHandler)).Methods("POST")
	// self used in paths rather than identifier as the identifier is JWT tokens passed in the request headers
	r.HandleFunc("/v1/tokens/self", contextAndErrors(api.SignOutHandler)).Methods("DELETE")
	r.HandleFunc("/v1/tokens/self", contextAndErrors(api.RefreshHandler)).Methods("PUT")
	r.HandleFunc("/v1/users", contextAndErrors(api.CreateUserHandler)).Methods("POST")
	r.HandleFunc("/v1/users", contextAndErrors(api.ListUsersHandler)).Methods("GET")
	r.HandleFunc("/v1/users/{id}", contextAndErrors(api.GetUserHandler)).Methods("GET")
	r.HandleFunc("/v1/users/{id}", contextAndErrors(api.UpdateUserHandler)).Methods("PUT")
	// self used in paths rather than identifier as the identifier is a Cognito Session string in change password requests
	// the user id is not yet available from the previous responses
	r.HandleFunc("/v1/users/self/password", contextAndErrors(api.ChangePasswordHandler)).Methods("PUT")
	r.HandleFunc("/v1/password-reset", contextAndErrors(api.PasswordResetHandler)).Methods("POST")
	return api, nil
}

func writeErrorResponse(ctx context.Context, w http.ResponseWriter, errorResponse *models.ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	// process custom headers
	if errorResponse.Headers != nil {
		for key := range errorResponse.Headers {
			w.Header().Set(key, errorResponse.Headers[key])
		}
	}
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
	// process custom headers
	if successResponse.Headers != nil {
		for key := range successResponse.Headers {
			w.Header().Set(key, successResponse.Headers[key])
		}
	}
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
		http.StatusInternalServerError,
		nil,
	)
}

func handleBodyUnmarshalError(ctx context.Context, err error) *models.ErrorResponse {
	return models.NewErrorResponse([]error{models.NewError(ctx, err, models.JSONUnmarshalError, models.ErrorUnmarshalFailedDescription)},
		http.StatusInternalServerError,
		nil,
	)
}

func initialiseRoleGroups(ctx context.Context, cognitoClient cognito.Client, userPoolId string) error {
	adminGroup := models.NewAdminRoleGroup()
	adminGroupCreateInput := adminGroup.BuildCreateGroupRequest(userPoolId)
	_, err := cognitoClient.CreateGroup(adminGroupCreateInput)
	if err != nil {
		cognitoErr := models.NewCognitoError(ctx, err, "CreateGroup request for admin group from API start up")
		if cognitoErr.Code != models.GroupExistsError {
			return cognitoErr
		}
	}

	publisherGroup := models.NewPublisherRoleGroup()
	publisherGroupCreateInput := publisherGroup.BuildCreateGroupRequest(userPoolId)
	_, err = cognitoClient.CreateGroup(publisherGroupCreateInput)
	if err != nil {
		cognitoErr := models.NewCognitoError(ctx, err, "CreateGroup request for publisher group from API start up")
		if cognitoErr.Code != models.GroupExistsError {
			return cognitoErr
		}
	}

	return nil
}
