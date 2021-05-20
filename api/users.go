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
	"github.com/google/uuid"
	"github.com/sethvargo/go-password/password"
)

//CreateUserHandler creates a new user and returns a http handler interface
func (api *API) CreateUserHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Event(ctx, "starting to generate a new user", log.INFO)
		defer req.Body.Close()

		var errorList []apierrors.IndividualError

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

		forename := user.Forename
		// validate forename
		if forename == "" {
			log.Event(ctx, apierrors.InvalidForenameErrorMessage, log.ERROR)
			errorList = append(errorList, apierrors.IndividualErrorBuilder(apierrors.ErrInvalidForename, apierrors.InvalidForenameErrorMessage, content.ValidForenameErrorField, content.ValidForenameErrorParam))
		}

		surname := user.Surname
		// validate forename
		if surname == "" {
			log.Event(ctx, apierrors.InvalidSurnameErrorMessage, log.ERROR)
			errorList = append(errorList, apierrors.IndividualErrorBuilder(apierrors.ErrInvalidSurname, apierrors.InvalidSurnameErrorMessage, content.ValidSurnameErrorField, content.ValidSurnameErrorParam))
		}

		email := user.Email
		// validate email
		if !validation.ValidateONSEmail(email) {
			log.Event(ctx, apierrors.InvalidErrorMessage, log.ERROR)
			errorList = append(errorList, apierrors.IndividualErrorBuilder(apierrors.ErrInvalidEmail, apierrors.InvalidErrorMessage, content.ValidEmailErrorField, content.ValidEmailErrorParam))
		}

		// duplicate email check
		emailCheckInput, _ := ListUsersModel("email = \""+email+"\"", "email", int64(1), &api.UserPoolId)
		emailCheckResp, err := api.CognitoClient.ListUsers(emailCheckInput)
		if err != nil {
			log.Event(ctx, content.ListUsersErrorMessage, log.ERROR)
			apierrors.HandleUnexpectedError(ctx, w, err, content.ListUsersErrorMessage, content.ListUsersErrorField, content.ListUsersErrorParam)
			return
		}
		if len(emailCheckResp.Users) > 0 {
			log.Event(ctx, content.ListUsersErrorMessage, log.ERROR)
			errorList = append(errorList, apierrors.IndividualErrorBuilder(apierrors.ErrDuplicateEmail, content.DuplicateEmailFound, content.ListUsersErrorField, content.ListUsersErrorParam))
		}

		// report validation errors in response
		if len(errorList) != 0 {
			apierrors.WriteErrorResponse(ctx, w, http.StatusBadRequest, apierrors.ErrorResponseBodyBuilder(errorList))
			return
		}

		newUser, err := CreateNewUserModel(forename, surname, uuid.NewString(), email, tempPassword, api.UserPoolId)
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

//CreateNewUserModel creates and returns *AdminCreateUserInput
func CreateNewUserModel(forename string, surname, userId string, email string, tempPass string, userPoolId string) (*cognitoidentityprovider.AdminCreateUserInput, error) {
	var (
		deliveryMethod, forenameAttrName, surnameAttrName, emailAttrName string = "EMAIL", "name", "family_name", "email"
	)

	user := &models.CreateUserInput{
		UserInput: &cognitoidentityprovider.AdminCreateUserInput{
			UserAttributes: []*cognitoidentityprovider.AttributeType{
				{
					Name:  &forenameAttrName,
					Value: &forename,
				},
				{
					Name:  &surnameAttrName,
					Value: &surname,
				},
				{
					Name:  &emailAttrName,
					Value: &email,
				},
			},
			DesiredDeliveryMediums: []*string{
				&deliveryMethod,
			},
			TemporaryPassword: &tempPass,
			UserPoolId:        &userPoolId,
			Username:          &userId,
		},
	}
	return user.UserInput, nil
}

//ListUsersModel creates and returns *ListUsersInput
func ListUsersModel(filterString string, requiredAttribute string, limit int64, userPoolId *string) (*cognitoidentityprovider.ListUsersInput, error) {
	user := &models.ListUsersInput{
		ListUsersInput: &cognitoidentityprovider.ListUsersInput{
			AttributesToGet: []*string{
				&requiredAttribute,
			},
			Filter:     &filterString,
			Limit:      &limit,
			UserPoolId: userPoolId,
		},
	}

	return user.ListUsersInput, nil
}
