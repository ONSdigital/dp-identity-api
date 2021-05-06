package api

import (
	"context"
	"encoding/json"
	"net/http"

	"io/ioutil"

	"github.com/ONSdigital/dp-identity-api/apierrors"
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
			apierrors.HandleUnexpectedError(ctx, w, err, "failed to generate password", "", "")
			return
		}

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			apierrors.HandleUnexpectedError(ctx, w, err, "api endpoint POST user returned an error reading request body", "", "")
			return
		}

		user := models.UserParams{}
		err = json.Unmarshal(body, &user)
		if err != nil{
			apierrors.HandleUnexpectedError(ctx, w, err, "api endpoint POST user returned an error unmarshalling request body", "", "")
			return
		}

		username := user.UserName
		email := user.Email

		newUser, err := CreateNewUserModel(ctx, username, tempPassword, email, api.UserPoolId)
		if err != nil {
			apierrors.HandleUnexpectedError(ctx, w, err, "Failed to create new user model", "", "")
			return
		}

		//Create user in cognito
		resultUser, err := api.CognitoClient.AdminCreateUser(newUser)
		if err != nil {
			apierrors.HandleUnexpectedError(ctx, w, err, "Failed to create new user in user pool", "", "")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, err := json.Marshal(resultUser)
		if err != nil {
			apierrors.HandleUnexpectedError(ctx, w, err, "Failed to marshall json response", "", "")
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(jsonResponse)
		if err != nil {
			apierrors.HandleUnexpectedError(ctx, w, err, "Failed to write http response", "", "")
			return
		}
	}
}

//CreateNewUserModel creates and returns AdminCreateUserInput
func CreateNewUserModel(ctx context.Context, username string, tempPass string, emailId string, userPoolId string) (*cognitoidentityprovider.AdminCreateUserInput, error) {
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
