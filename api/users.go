package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ONSdigital/dp-identity-api/api/content"
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
			log.Event(ctx, content.PasswordErrorMessage, log.ERROR)
			apierrors.HandleUnexpectedError(ctx, w, err, content.PasswordErrorMessage, content.PasswordErrorField, content.PasswordErrorParam)
			return
		}

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Event(ctx, content.RequestErrorMessage, log.ERROR)
			apierrors.HandleUnexpectedError(ctx, w, err, content.RequestErrorMessage, content.RequestErrorField, content.RequestErrorParam)
			return
		}

		user := models.UserParams{}
		err = json.Unmarshal(body, &user)
		if err != nil {
			log.Event(ctx, content.UnmarshallingErrorMessage, log.ERROR)
			apierrors.HandleUnexpectedError(ctx, w, err, content.UnmarshallingErrorMessage, content.UnmarshallingErrorField, content.UnmarshallingErrorParam)
			return
		}

		username := user.UserName
		// validate username
		if len(username) == 0  {
			log.Event(ctx, content.ValidUserNameErrorField, log.ERROR)
			errorList = append(errorList, apierrors.IndividualErrorBuilder(apierrors.ErrInvalidUserName, apierrors.InvalidUserNameMessage, content.ValidUserNameErrorField, content.ValidUserNameErrorParam))
		}

		email := user.Email
		// validate email
		if !validation.IsEmailValid(email) {
			log.Event(ctx, content.ValidEmailErrorField, log.ERROR)
			errorList = append(errorList, apierrors.IndividualErrorBuilder(apierrors.ErrInvalidEmail, apierrors.InvalidErrorMessage, content.ValidEmailErrorField, content.ValidEmailErrorParam))
		}

		// report validation errors in response
		if len(errorList) != 0 {
			apierrors.WriteErrorResponse(ctx, w, http.StatusBadRequest, apierrors.ErrorResponseBodyBuilder(errorList))
			return
		}

		newUser, err := CreateNewUserModel(ctx, username, tempPassword, email, api.UserPoolId)
		if err != nil {
			log.Event(ctx, content.NewUserModelErrorMessage, log.ERROR)
			apierrors.HandleUnexpectedError(ctx, w, err, content.NewUserModelErrorMessage, content.NewUserModelErrorField, content.NewUserModelErrorParam)
			return
		}

		//Create user in cognito - and handle any errors returned
		resultUser, err := api.CognitoClient.AdminCreateUser(newUser)
		if err != nil {
			if strings.Contains(err.Error(), content.InternalErrorException) {
				log.Event(ctx, content.AdminCreateUserErrorMessage, log.ERROR)
				apierrors.HandleUnexpectedError(ctx, w, err, content.AdminCreateUserErrorMessage, content.AdminCreateUserErrorField, content.AdminCreateUserErrorParam)
				return
			} else {
				log.Event(ctx, err.Error(), log.ERROR)
				errorList = append(errorList, apierrors.IndividualErrorBuilder(err, err.Error(), content.AdminCreateUserErrorField, content.AdminCreateUserErrorParam))
				apierrors.WriteErrorResponse(ctx, w, http.StatusBadRequest, apierrors.ErrorResponseBodyBuilder(errorList))
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, err := json.Marshal(resultUser)
		if err != nil {
			log.Event(ctx, content.MarshallingNewUserErrorMessage, log.ERROR)
			apierrors.HandleUnexpectedError(ctx, w, err, content.MarshallingNewUserErrorMessage, content.MarshallingNewUserErrorField, content.MarshallingNewUserErrorParam)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(jsonResponse)
		if err != nil {
			log.Event(ctx, content.HttpResponseErrorMessage, log.ERROR)
			apierrors.HandleUnexpectedError(ctx, w, err, content.HttpResponseErrorMessage, content.HttpResponseErrorField, content.HttpResponseErrorParam)
			return
		}
	}
}

//CreateNewUserModel creates and returns AdminCreateUserInput
func CreateNewUserModel(ctx context.Context, username string, tempPass string, emailId string, userPoolId string) (*cognitoidentityprovider.AdminCreateUserInput, error) {
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
