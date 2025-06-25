package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ONSdigital/dp-authorisation/v2/authorisation"
	"github.com/ONSdigital/dp-identity-api/v2/cognito"
	"github.com/ONSdigital/dp-identity-api/v2/jwks"
	"github.com/ONSdigital/dp-identity-api/v2/models"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/gorilla/mux"
)

var (
	IDTokenHeaderName      = "ID"
	AccessTokenHeaderName  = "Authorization"
	RefreshTokenHeaderName = "Refresh"
	WWWAuthenticateName    = "WWW-Authenticate"
	ONSRealm               = "Florence publishing platform"
	Charset                = "UTF-8"
	NewPasswordChallenge   = types.ChallengeNameTypeNewPasswordRequired
	DefaultBackOffSchedule = []time.Duration{
		1 * time.Second,
		3 * time.Second,
		10 * time.Second,
	}
)

// API provides a struct to wrap the api around
type API struct {
	Router              *mux.Router
	CognitoClient       cognito.Client
	UserPoolID          string
	ClientID            string
	ClientSecret        string
	ClientAuthFlow      types.AuthFlowType
	AWSRegion           string
	AllowedDomains      []string
	APIRequestFilter    map[string]map[string]string
	JWKSManager         jwks.Manager
	BlockPlusAddressing bool
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

// Setup function sets up the api and returns an api
func Setup(ctx context.Context,
	r *mux.Router,
	cognitoClient cognito.Client,
	userPoolID, clientID, clientSecret, awsRegion string,
	clientAuthFlow types.AuthFlowType,
	blockPlusAddressing bool,
	allowedDomains []string,
	auth authorisation.Middleware,
	jwksManager jwks.Manager) (*API, error) {
	// Return an error if empty required parameter was passed.
	if userPoolID == "" || clientID == "" || clientSecret == "" || awsRegion == "" || clientAuthFlow == "" || allowedDomains == nil || len(allowedDomains) == 0 || jwksManager == nil {
		return nil, models.NewError(ctx, nil, models.MissingConfigError, models.MissingConfigDescription)
	}

	if err := initialiseRoleGroups(ctx, cognitoClient, userPoolID); err != nil {
		return nil, err
	}

	api := &API{
		Router:              r,
		CognitoClient:       cognitoClient,
		UserPoolID:          userPoolID,
		ClientID:            clientID,
		ClientSecret:        clientSecret,
		ClientAuthFlow:      clientAuthFlow,
		AWSRegion:           awsRegion,
		BlockPlusAddressing: blockPlusAddressing,
		AllowedDomains:      allowedDomains,
		APIRequestFilter: map[string]map[string]string{
			"/v1/users": {
				"active=true":  "status=\"Enabled\"",
				"active=false": "status=\"Disabled\"",
			},
		},
		JWKSManager: jwksManager,
	}

	r.HandleFunc("/v1/tokens", contextAndErrors(api.TokensHandler)).Methods(http.MethodPost)
	r.HandleFunc("/v1/tokens", auth.Require(UsersUpdatePermission, contextAndErrors(api.SignOutAllUsersHandler))).
		Methods(http.MethodDelete)
	// self used in paths rather than identifier as the identifier is JWT tokens passed in the request headers
	r.HandleFunc("/v1/tokens/self", contextAndErrors(api.SignOutHandler)).Methods(http.MethodDelete)
	r.HandleFunc("/v1/tokens/self", contextAndErrors(api.RefreshHandler)).Methods(http.MethodPut)
	r.HandleFunc("/v1/users", auth.Require(UsersCreatePermission, contextAndErrors(api.CreateUserHandler))).
		Methods(http.MethodPost)
	r.HandleFunc("/v1/users", auth.Require(UsersReadPermission, contextAndErrors(api.ListUsersHandler))).
		Methods(http.MethodGet)
	r.HandleFunc("/v1/users/{id}", auth.Require(UsersReadPermission, contextAndErrors(api.GetUserHandler))).
		Methods(http.MethodGet)
	r.HandleFunc("/v1/users/{id}", auth.Require(UsersUpdatePermission, contextAndErrors(api.UpdateUserHandler))).
		Methods(http.MethodPut)
	r.HandleFunc("/v1/users/{id}/password", auth.Require(UsersUpdatePermission, contextAndErrors(api.UserSetPasswordHandler))).
		Methods(http.MethodPost)
	r.HandleFunc("/v1/users/{id}/groups", auth.Require(UsersReadPermission, contextAndErrors(api.ListUserGroupsHandler))).
		Methods(http.MethodGet)
	// self used in paths rather than identifier as the identifier is a Cognito Session string in change password requests
	// the user id is not yet available from the previous responses
	r.HandleFunc("/v1/users/self/password", contextAndErrors(api.ChangePasswordHandler)).
		Methods(http.MethodPut)
	r.HandleFunc("/v1/password-reset", contextAndErrors(api.PasswordResetHandler)).
		Methods(http.MethodPost)
	r.HandleFunc("/v1/groups", auth.Require(GroupsReadPermission, contextAndErrors(api.ListGroupsHandler))).
		Methods(http.MethodGet)
	r.HandleFunc("/v1/groups", auth.Require(GroupsCreatePermission, contextAndErrors(api.CreateGroupHandler))).
		Methods(http.MethodPost)
	r.HandleFunc("/v1/groups/{id}", auth.Require(GroupsReadPermission, contextAndErrors(api.GetGroupHandler))).
		Methods(http.MethodGet)
	r.HandleFunc("/v1/groups-report", auth.Require(UsersReadPermission, contextAndErrors(api.ListGroupsUsersHandler))).
		Methods(http.MethodGet)
	r.HandleFunc("/v1/groups/{id}", auth.Require(GroupsEditPermission, contextAndErrors(api.UpdateGroupHandler))).
		Methods(http.MethodPut)
	r.HandleFunc("/v1/groups/{id}", auth.Require(GroupsDeletePermission, contextAndErrors(api.DeleteGroupHandler))).
		Methods(http.MethodDelete)
	r.HandleFunc("/v1/groups/{id}/members", auth.Require(GroupsEditPermission, contextAndErrors(api.AddUserToGroupHandler))).
		Methods(http.MethodPost)
	r.HandleFunc("/v1/groups/{id}/members", auth.Require(GroupsEditPermission, contextAndErrors(api.SetGroupUsersHandler))).
		Methods(http.MethodPut)
	r.HandleFunc("/v1/groups/{id}/members", auth.Require(GroupsReadPermission, contextAndErrors(api.ListUsersInGroupHandler))).
		Methods(http.MethodGet)
	r.HandleFunc("/v1/groups/{id}/members/{user_id}", auth.Require(GroupsEditPermission, contextAndErrors(api.RemoveUserFromGroupHandler))).
		Methods(http.MethodDelete)
	r.HandleFunc("/v1/jwt-keys", contextAndErrors(api.CognitoPoolJWKSHandler)).
		Methods(http.MethodGet)
	return api, nil
}

func writeErrorResponse(ctx context.Context, w http.ResponseWriter, errorResponse *models.ErrorResponse) {
	var jsonResponse []byte
	var err error
	w.Header().Set("Content-Type", "application/json")
	// process custom headers
	if errorResponse.Headers != nil {
		for key := range errorResponse.Headers {
			w.Header().Set(key, errorResponse.Headers[key])
		}
	}
	w.WriteHeader(errorResponse.Status)
	if errorResponse.Status == http.StatusInternalServerError {
		jsonResponse, err = json.Marshal(models.Error{Code: models.InternalError, Description: models.InternalErrorDescription})
	} else {
		jsonResponse, err = json.Marshal(errorResponse)
	}
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
	return models.NewErrorResponse(http.StatusBadRequest,
		nil,
		models.NewError(ctx, err, models.BodyReadError, models.BodyReadFailedDescription),
	)
}

func handleBodyUnmarshalError(ctx context.Context, err error) *models.ErrorResponse {
	return models.NewErrorResponse(http.StatusBadRequest,
		nil,
		models.NewError(ctx, err, models.JSONUnmarshalError, models.ErrorUnmarshalFailedDescription),
	)
}

func initialiseRoleGroups(ctx context.Context, cognitoClient cognito.Client, userPoolID string) error {
	adminGroup := models.NewAdminRoleGroup()
	adminGroupCreateInput := adminGroup.BuildCreateGroupRequest(userPoolID)
	_, err := cognitoClient.CreateGroup(ctx, adminGroupCreateInput)
	if err != nil && !models.IsGroupExistsError(err) {
		cognitoErr := models.NewCognitoError(ctx, err, "CreateGroup request for admin group from API start up")
		if cognitoErr.Code != models.GroupExistsError {
			return cognitoErr
		}
	}

	publisherGroup := models.NewPublisherRoleGroup()
	publisherGroupCreateInput := publisherGroup.BuildCreateGroupRequest(userPoolID)
	_, err = cognitoClient.CreateGroup(ctx, publisherGroupCreateInput)
	if err != nil && !models.IsGroupExistsError(err) {
		cognitoErr := models.NewCognitoError(ctx, err, "CreateGroup request for publisher group from API start up")
		if cognitoErr.Code != models.GroupExistsError {
			return cognitoErr
		}
	}

	return nil
}
