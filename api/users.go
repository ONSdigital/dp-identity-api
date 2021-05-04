package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/log.go/log"

	uuid "github.com/satori/go.uuid"
)

//NewID generates a random uuid and returns it as a string.
var NewID = func() string {
	return uuid.NewV4().String()
}

func (api *API) CreateUserHandler(ctx context.Context) http.HandlerFunc {
	log.Event(ctx, "starting to generate a new user", log.INFO)
	return func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		id := NewID()

		if err := req.ParseForm(); err != nil {
			log.Event(ctx, "failed to parse request form", log.ERROR, log.Error(err))
			return
		}

		username := req.Form.Get("username")
		tempPassword := req.Form.Get("password")
		email := req.Form.Get("email")

		//Create user in cognito
		user, err := api.CognitoClient.AdminCreateUser(ctx)
		if err != nil {
			log.Event(ctx, "creating user failed", log.Error(err), log.ERROR)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, err := json.Marshal(user)
		if err != nil {
			log.Event(ctx, "marshalling response failed", log.Error(err), log.ERROR)
			http.Error(w, "Failed to marshall json response", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(jsonResponse)
		if err != nil {
			log.Event(ctx, "writing response failed", log.Error(err), log.ERROR)
			http.Error(w, "Failed to write http response", http.StatusInternalServerError)
			return
		}
	}
}
