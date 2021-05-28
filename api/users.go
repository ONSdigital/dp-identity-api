package api

import (
	"context"
	"encoding/json"
	"github.com/ONSdigital/dp-identity-api/apierrorsdeprecated"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/ONSdigital/log.go/log"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
)

//CreateUserHandler creates a new user and returns a http handler interface
func (api *API) CreateUserHandler(w http.ResponseWriter, req *http.Request, ctx context.Context) (*models.SuccessResponse, *models.ErrorResponse) {
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, handleBodyReadError(ctx, err)
	}

	user := models.UserParams{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, handleBodyUnmarshalError(ctx, err)
	}

	err = user.GeneratePassword(ctx)
	if err != nil {
		return nil, models.NewErrorResponse([]error{err}, http.StatusInternalServerError)
	}

	validationErrs := user.ValidateRegistration(ctx)

	listUserInput := user.BuildListUserRequest("email = \""+user.Email+"\"", "email", int64(1), &api.UserPoolId)
	listUserResp, err := api.CognitoClient.ListUsers(listUserInput)
	if err != nil {
		return nil, models.NewErrorResponse([]error{models.NewCognitoError(ctx, err, "Cognito ListUsers request from create users endpoint")}, http.StatusInternalServerError)
	}
	duplicateEmailErr := user.CheckForDuplicateEmail(ctx, listUserResp)
	if duplicateEmailErr != nil {
		log.Event(ctx, apierrorsdeprecated.ListUsersErrorDescription, log.ERROR)
		validationErrs = append(validationErrs, duplicateEmailErr)
	}

	if len(validationErrs) != 0 {
		return nil, models.NewErrorResponse(validationErrs, http.StatusBadRequest)
	}

	createUserRequest := user.BuildCreateUserRequest(uuid.NewString(), api.UserPoolId)

	resultUser, err := api.CognitoClient.AdminCreateUser(createUserRequest)
	if err != nil {
		responseErr := models.NewCognitoError(ctx, err, "AdminCreateUser request from create user endpoint")
		if responseErr.Code == models.InternalError {
			return nil, models.NewErrorResponse([]error{responseErr}, http.StatusInternalServerError)
		} else {
			return nil, models.NewErrorResponse([]error{responseErr}, http.StatusBadRequest)
		}
	}

	jsonResponse, responseErr := user.BuildSuccessfulJsonResponse(ctx, resultUser)
	if responseErr != nil {
		return nil, models.NewErrorResponse([]error{responseErr}, http.StatusInternalServerError)
	}

	return models.NewSuccessResponse(jsonResponse, http.StatusCreated), nil
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
