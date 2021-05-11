package api

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/ONSdigital/dp-identity-api/apierrors"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/ONSdigital/dp-identity-api/validation"
	"github.com/ONSdigital/log.go/log"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/sethvargo/go-password/password"
)

//CreateUserHandler creates a new user and returns a http handler interface
func (api *API) CreateUserHandler(ctx context.Context) http.HandlerFunc {
	log.Event(ctx, "starting to generate a new user", log.INFO)
	return func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		
		var errorList []models.IndividualError

		tempPassword, err := password.Generate(14, 1, 1, false, false)
		if err != nil {
			log.Event(ctx, passwordErrorMessage, log.ERROR)
			apierrors.HandleUnexpectedError(ctx, w, err, passwordErrorMessage, passwordErrorField, passwordErrorParam)
			return
		}

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Event(ctx, requestErrorMessage, log.ERROR)
			apierrors.HandleUnexpectedError(ctx, w, err, requestErrorMessage, requestErrorField, requestErrorParam)
			return
		}

		user := models.UserParams{}
		err = json.Unmarshal(body, &user)
		if err != nil {
			log.Event(ctx, unmarshallingErrorMessage, log.ERROR)
			apierrors.HandleUnexpectedError(ctx, w, err, unmarshallingErrorMessage, unmarshallingErrorField, unmarshallingErrorParam)
			return
		}

		username := user.UserName
		// validate username
		if len(username) == 0  {
			log.Event(ctx, validUserNameErrorField, log.ERROR)
			errorList = append(errorList, apierrors.IndividualErrorBuilder(apierrors.ErrInvalidUserName, apierrors.InvalidUserNameMessage, validUserNameErrorField, validUserNameErrorParam))
		}

		email := user.Email
		// validate email
		if !validation.IsEmailValid(email) {
			log.Event(ctx, validEmailErrorField, log.ERROR)
			errorList = append(errorList, apierrors.IndividualErrorBuilder(apierrors.ErrInvalidEmail, apierrors.InvalidErrorMessage, validEmailErrorField, validEmailErrorParam))
		}

		// report validation errors in response
		if len(errorList) != 0 {
			apierrors.WriteErrorResponse(ctx, w, http.StatusBadRequest, apierrors.ErrorResponseBodyBuilder(errorList))
			return
		}

		newUser, err := CreateNewUserModel(ctx, username, tempPassword, email, api.UserPoolId)
		if err != nil {
			log.Event(ctx, newUserModelErrorMessage, log.ERROR)
			apierrors.HandleUnexpectedError(ctx, w, err, newUserModelErrorMessage, newUserModelErrorField, newUserModelErrorParam)
			return
		}

		//Create user in cognito
		resultUser, err := api.CognitoClient.AdminCreateUser(newUser)
		if err != nil {
			log.Event(ctx, adminCreateUserErrorMessage, log.ERROR)
			apierrors.HandleUnexpectedError(ctx, w, err, adminCreateUserErrorMessage, adminCreateUserErrorField, adminCreateUserErrorParam)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, err := json.Marshal(resultUser)
		if err != nil {
			log.Event(ctx, marshallingNewUserErrorMessage, log.ERROR)
			apierrors.HandleUnexpectedError(ctx, w, err, marshallingNewUserErrorMessage, marshallingNewUserErrorField, marshallingNewUserErrorParam)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(jsonResponse)
		if err != nil {
			log.Event(ctx, httpResponseErrorMessage, log.ERROR)
			apierrors.HandleUnexpectedError(ctx, w, err, httpResponseErrorMessage, httpResponseErrorField, httpResponseErrorParam)
			return
		}
	}
}

//CreateNewUserModel creates and returns AdminCreateUserInput
func CreateNewUserModel(ctx context.Context, username string, tempPass string, emailId string, userPoolId string) (*cognitoidentityprovider.AdminCreateUserInput, error) {
	// Return an error if empty id was passed.
	if userPoolId == "" {
		return nil, errors.New(userPoolIdNotFoundMessage)
	}

	var (
		deliveryMethod, emailAttrName string = "EMAIL", "email"
	)

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

// message, field and param error constants
const passwordErrorMessage           = "failed to generate password"
const passwordErrorField             = "temp password"
const passwordErrorParam             = "error creating temp password"

const requestErrorMessage            = "api endpoint POST user returned an error reading request body"
const requestErrorField              = "request body"
const requestErrorParam              = "error in reading request body"

const unmarshallingErrorMessage      = "api endpoint POST user returned an error unmarshalling request body"
const unmarshallingErrorField        = "unmarshalling"
const unmarshallingErrorParam        = "error unmarshalling request body"

const validUserNameErrorField        = "validating username"
const validUserNameErrorParam        = "error validating username"

const validEmailErrorField           = "validating email"
const validEmailErrorParam           = "error validating email"

const newUserModelErrorMessage       = "Failed to create new user model"
const newUserModelErrorField         = "create new user model"
const newUserModelErrorParam         = "error creating new user model"

const adminCreateUserErrorMessage    = "Failed to create new user in user pool"
const adminCreateUserErrorField      = "create new user pool user"
const adminCreateUserErrorParam      = "error creating new user pool user"

const marshallingNewUserErrorMessage = "Failed to marshall json response"
const marshallingNewUserErrorField   = "marshalling"
const marshallingNewUserErrorParam   = "error marshalling new user response"

const httpResponseErrorMessage       = "Failed to write http response"
const httpResponseErrorField         = "response"
const httpResponseErrorParam         = "error writing response"

const userPoolIdNotFoundMessage      = "userPoolId must not be an empty string"