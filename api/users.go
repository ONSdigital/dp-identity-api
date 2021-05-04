package api

import (
	"context"
	"encoding/json"
	"net/http"

	models "github.com/ONSdigital/dp-identity-api/models"
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

		id := NewID()

		if err := req.ParseForm(); err != nil {
			log.Event(ctx, "failed to parse request form", log.ERROR, log.Error(err))
			return
		}
		defer req.Body.Close()
		
		username := req.Form.Get("username")
		tempPassword := req.Form.Get("password")
		email := req.Form.Get("email")
		
		newUser,err := CreateNewUserModel(ctx,id,username,tempPassword,email)
		if err != nil {
			log.Event(ctx, "creating new user model failed", log.Error(err), log.ERROR)
			http.Error(w, "Failed to create new user model", http.StatusInternalServerError)
			return
		}
		
		//Create user in cognito
		resultUser, err := api.CognitoClient.AdminCreateUser(newUser)
		if err != nil {
			log.Event(ctx, "creating user failed", log.Error(err), log.ERROR)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, err := json.Marshal(resultUser)
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

	func CreateNewUserModel(ctx context.Context, id string, username string, tempPass string, email string) (models.CognitoUser,error){

		log.Event(ctx, "creating user", log.Data{"id": id})

		// Return an error if empty id was passed.
		if id == "" {
			return nil, errors.New("id must not be an empty string")
		}
		return models.CognitoUser{
			TemporaryPassword : &tempPass,
			UserAttributes{
				Name: 
				Value: &email
			}
			UserPoolId: id,
			Username: &username,
			DesiredDeliveryMediums: "EMAIL"

		}, nil

	}
