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

// DescribeUserPoolHandler returns details of the user pool
func (api *API) DescribeUserPoolHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		vars := mux.Vars(req)
		userpool := vars["userpool"]
		_, err := api.CognitoClient.DescribeUserPool(&cognitoidentityprovider.DescribeUserPoolInput{UserPoolId: &userpool})

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to load user pool details", http.StatusInternalServerError)
			return
		}

		response := HealthcheckResponse{
			Message: healthyMessage,
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Event(ctx, "marshalling response failed", log.Error(err), log.ERROR)
			http.Error(w, "Failed to marshall json response", http.StatusInternalServerError)
		}

		_, err = w.Write(jsonResponse)
		if err != nil {
			log.Event(ctx, "writing response failed", log.Error(err), log.ERROR)
			http.Error(w, "Failed to write http response", http.StatusInternalServerError)
			return
		}
	}
}
