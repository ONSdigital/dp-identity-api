package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ONSdigital/dp-identity-api/apierrorsdeprecated"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/ONSdigital/dp-identity-api/validation"
	"github.com/ONSdigital/log.go/log"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/google/uuid"
	"github.com/sethvargo/go-password/password"
)

//CreateUserHandler creates a new user and returns a http handler interface
func (api *API) CreateUserHandler(ctx context.Context) http.HandlerFunc {
	log.Event(ctx, "starting to generate a new user", log.INFO)
	return func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		var errorList []models.Error

		tempPassword, err := password.Generate(14, 1, 1, false, false)
		if err != nil {
			log.Event(ctx, apierrorsdeprecated.PasswordErrorDescription, log.ERROR)
			apierrorsdeprecated.HandleUnexpectedError(ctx, w, err, apierrorsdeprecated.PasswordErrorDescription)
			return
		}

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Event(ctx, apierrorsdeprecated.RequestErrorDescription, log.ERROR)
			apierrorsdeprecated.HandleUnexpectedError(ctx, w, err, apierrorsdeprecated.RequestErrorDescription)
			return
		}

		user := models.UserParams{}
		err = json.Unmarshal(body, &user)
		if err != nil {
			log.Event(ctx, apierrorsdeprecated.UnmarshallingErrorDescription, log.ERROR)
			apierrorsdeprecated.HandleUnexpectedError(ctx, w, err, apierrorsdeprecated.UnmarshallingErrorDescription)
			return
		}

		forename := user.Forename
		// validate forename
		if forename == "" {
			log.Event(ctx, apierrorsdeprecated.InvalidForenameErrorDescription, log.ERROR)
			errorList = append(errorList, apierrorsdeprecated.IndividualErrorBuilder(apierrorsdeprecated.ErrInvalidForename, apierrorsdeprecated.InvalidForenameErrorDescription))
		}

		surname := user.Surname
		// validate forename
		if surname == "" {
			log.Event(ctx, apierrorsdeprecated.InvalidSurnameErrorDescription, log.ERROR)
			errorList = append(errorList, apierrorsdeprecated.IndividualErrorBuilder(apierrorsdeprecated.ErrInvalidSurname, apierrorsdeprecated.InvalidSurnameErrorDescription))
		}

		email := user.Email
		// validate email
		if !validation.ValidateONSEmail(email) {
			log.Event(ctx, apierrorsdeprecated.InvalidErrorDescription, log.ERROR)
			errorList = append(errorList, apierrorsdeprecated.IndividualErrorBuilder(apierrorsdeprecated.ErrInvalidEmail, apierrorsdeprecated.InvalidErrorDescription))
		}

		// duplicate email check
		emailCheckInput, _ := ListUsersModel("email = \""+email+"\"", "email", int64(1), &api.UserPoolId)
		emailCheckResp, err := api.CognitoClient.ListUsers(emailCheckInput)
		if err != nil {
			log.Event(ctx, apierrorsdeprecated.ListUsersErrorDescription, log.ERROR)
			apierrorsdeprecated.HandleUnexpectedError(ctx, w, err, apierrorsdeprecated.ListUsersErrorDescription)
			return
		}
		if len(emailCheckResp.Users) > 0 {
			log.Event(ctx, apierrorsdeprecated.ListUsersErrorDescription, log.ERROR)
			errorList = append(errorList, apierrorsdeprecated.IndividualErrorBuilder(apierrorsdeprecated.ErrDuplicateEmail, apierrorsdeprecated.DuplicateEmailFound))
		}

		// report validation errors in response
		if len(errorList) != 0 {
			apierrorsdeprecated.WriteErrorResponse(ctx, w, http.StatusBadRequest, apierrorsdeprecated.ErrorResponseBodyBuilder(errorList))
			return
		}

		newUser, err := CreateNewUserModel(forename, surname, uuid.NewString(), email, tempPassword, api.UserPoolId)
		if err != nil {
			log.Event(ctx, apierrorsdeprecated.NewUserModelErrorDescription, log.ERROR)
			apierrorsdeprecated.HandleUnexpectedError(ctx, w, err, apierrorsdeprecated.NewUserModelErrorDescription)
			return
		}

		//Create user in cognito - and handle any errors returned
		resultUser, err := api.CognitoClient.AdminCreateUser(newUser)
		if err != nil {
			if strings.Contains(err.Error(), apierrorsdeprecated.InternalErrorException) {
				log.Event(ctx, apierrorsdeprecated.AdminCreateUserErrorDescription, log.ERROR)
				apierrorsdeprecated.HandleUnexpectedError(ctx, w, err, apierrorsdeprecated.AdminCreateUserErrorDescription)
				return
			} else {
				log.Event(ctx, err.Error(), log.ERROR)
				errorList = append(errorList, apierrorsdeprecated.IndividualErrorBuilder(err, err.Error()))
				apierrorsdeprecated.WriteErrorResponse(ctx, w, http.StatusBadRequest, apierrorsdeprecated.ErrorResponseBodyBuilder(errorList))
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, err := json.Marshal(resultUser)
		if err != nil {
			log.Event(ctx, apierrorsdeprecated.MarshallingNewUserErrorDescription, log.ERROR)
			apierrorsdeprecated.HandleUnexpectedError(ctx, w, err, apierrorsdeprecated.MarshallingNewUserErrorDescription)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(jsonResponse)
		if err != nil {
			log.Event(ctx, apierrorsdeprecated.HttpResponseErrorDescription, log.ERROR)
			apierrorsdeprecated.HandleUnexpectedError(ctx, w, err, apierrorsdeprecated.HttpResponseErrorDescription)
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
