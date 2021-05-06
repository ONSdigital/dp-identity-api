package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"io/ioutil"

	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/ONSdigital/log.go/log"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/sethvargo/go-password/password"
)

//CreateUserHandler creates a new user and returns a http handler interface
func (api *API) CreateUserHandler(ctx context.Context) http.HandlerFunc {
	log.Event(ctx, "starting to generate a new user", log.INFO)
	return func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		tempPassword, err := password.Generate(32, 10, 10, false, false)
		if err != nil {
			log.Event(ctx, "failed to generate password", log.ERROR, log.Error(err))
			http.Error(w, "Failed to generate password", http.StatusInternalServerError)
			return
		}

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Event(ctx, "api endpoint POST user returned an error reading the request body", log.Error(err), log.ERROR)
			http.Error(w, "Failed to read the request body", http.StatusInternalServerError)
			return
		}

		user := models.UserParams{}
		err = json.Unmarshal(body, &user)
		if err != nil {
			log.Event(ctx, "api endpoint POST user returned an error unmarshalling the body", log.Error(err), log.ERROR)
			http.Error(w, "Failed to unmarshall the body", http.StatusInternalServerError)
			return
		}

		username := user.UserName
		email := user.Email

		newUser, err := CreateNewUserModel(ctx, username, tempPassword, email, api.UserPoolId)
		if err != nil {
			log.Event(ctx, "creating new user failed model", log.Error(err), log.ERROR)
			http.Error(w, "Failed to create user model", http.StatusInternalServerError)
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

//CreateNewUserModel creates and returns AdminCreateUserInput
func CreateNewUserModel(ctx context.Context, username string, tempPass string, emailId string, userPoolId string) (*cognitoidentityprovider.AdminCreateUserInput, error) {
	// Return an error if empty id was passed.
	if userPoolId == "" {
		return nil, errors.New("userPoolId must not be an empty string")
	}
	emailAttrName := "email"
	deliveryMethod := "EMAIL"

	user := &models.CreateUserInput{
		UserInput: &cognitoidentityprovider.AdminCreateUserInput{
			UserAttributes: []*cognitoidentityprovider.AttributeType{
				{
					Name:  &emailAttrName,
					Value: &emailId,
				},
			},
			DesiredDeliveryMediums: []*string{
				&deliveryMethod,
			},
			TemporaryPassword: &tempPass,
			UserPoolId:        &userPoolId,
			Username:          &username,
		},
	}
	return user.UserInput, nil
}
