package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ONSdigital/log.go/log"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gorilla/mux"
	"net/http"
)

// getUserProfile returns a user profile from AWS Cognito
func (api *API) GetUserProfile(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		vars := mux.Vars(req)
		userName := vars["userName"]

		result, err := api.CognitoClient.AdminGetUser(&cognitoidentityprovider.AdminGetUserInput{UserPoolId: &api.UserPoolId, Username: &userName})

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to load user profile", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, err := json.Marshal(result)
		if err != nil {
			log.Event(ctx, "marshalling response failed", log.Error(err), log.ERROR)
			http.Error(w, "Failed to marshall json response", http.StatusInternalServerError)
			return
		}

		_, err = w.Write(jsonResponse)
		if err != nil {
			log.Event(ctx, "writing response failed", log.Error(err), log.ERROR)
			http.Error(w, "Failed to write http response", http.StatusInternalServerError)
			return
		}
	}
}
