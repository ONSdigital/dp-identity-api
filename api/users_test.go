package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"time"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/gorilla/mux"

	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	. "github.com/smartystreets/goconvey/convey"
)

const usersEndPoint = "http://localhost:25600/v1/users"
const userEndPoint = "http://localhost:25600/v1/users/abcd1234"
const changePasswordEndPoint = "http://localhost:25600/v1/users/self/password"
const requestResetEndPoint = "http://localhost:25600/v1/password-reset"
const userListGroupsEndPoint = "http://localhost:25600/v1/users/abcd1234/groups"

func TestCreateUserHandler(t *testing.T) {

	var (
		ctx                                               = context.Background()
		name, surname, status, email, invalidEmail string = "bob", "bobbings", "UNCONFIRMED", "foo_bar123@ext.ons.gov.uk", "foo_bar123@test.ons.gov.ie"
		userException                              string = "UsernameExistsException: User account already exists"
	)

	api, w, m := apiSetup()

	Convey("Admin create user - check expected responses", t, func() {
		adminCreateUsersTests := []struct {
			listUsersFunction   func(userInput *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error)
			createUsersFunction func(userInput *cognitoidentityprovider.AdminCreateUserInput) (*cognitoidentityprovider.AdminCreateUserOutput, error)
			httpResponse        int
		}{
			{
				// 200 response - no duplicate emails found
				func(userInput *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error) {
					users := &models.ListUsersOutput{
						ListUsersOutput: &cognitoidentityprovider.ListUsersOutput{
							Users: []*cognitoidentityprovider.UserType{},
						},
					}
					return users.ListUsersOutput, nil
				},
				// 201 response - user created
				func(userInput *cognitoidentityprovider.AdminCreateUserInput) (*cognitoidentityprovider.AdminCreateUserOutput, error) {
					user := &models.CreateUserOutput{
						UserOutput: &cognitoidentityprovider.AdminCreateUserOutput{
							User: &cognitoidentityprovider.UserType{
								Username:   &name,
								UserStatus: &status,
							},
						},
					}
					return user.UserOutput, nil
				},
				http.StatusCreated,
			},
			{
				// 200 response - no duplicate emails found
				func(userInput *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error) {
					users := &models.ListUsersOutput{
						ListUsersOutput: &cognitoidentityprovider.ListUsersOutput{
							Users: []*cognitoidentityprovider.UserType{},
						},
					}
					return users.ListUsersOutput, nil
				},
				// 400 response - user already exists
				func(userInput *cognitoidentityprovider.AdminCreateUserInput) (*cognitoidentityprovider.AdminCreateUserOutput, error) {
					var userExistsException cognitoidentityprovider.UsernameExistsException
					userExistsException.Message_ = &userException
					userExistsException.RespMetadata.StatusCode = http.StatusBadRequest

					return nil, &userExistsException
				},
				http.StatusBadRequest,
			},
			{
				// 400 response - duplicate email found
				func(userInput *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error) {
					users := &models.ListUsersOutput{
						ListUsersOutput: &cognitoidentityprovider.ListUsersOutput{
							Users: []*cognitoidentityprovider.UserType{
								{
									Username: &name,
								},
							},
						},
					}
					return users.ListUsersOutput, nil
				},
				nil,
				http.StatusBadRequest,
			},
			{
				// 200 response - no users found
				func(userInput *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error) {
					users := &models.ListUsersOutput{
						ListUsersOutput: &cognitoidentityprovider.ListUsersOutput{
							Users: []*cognitoidentityprovider.UserType{},
						},
					}
					return users.ListUsersOutput, nil
				},
				// 500 response - internal error exception
				func(userInput *cognitoidentityprovider.AdminCreateUserInput) (*cognitoidentityprovider.AdminCreateUserOutput, error) {
					var internalErrorException cognitoidentityprovider.InternalErrorException
					internalErrorException.Message_ = &userException
					internalErrorException.RespMetadata.StatusCode = http.StatusInternalServerError

					return nil, &internalErrorException
				},
				http.StatusInternalServerError,
			},
		}

		for _, tt := range adminCreateUsersTests {
			m.AdminCreateUserFunc = tt.createUsersFunction
			m.ListUsersFunc = tt.listUsersFunction

			postBody := map[string]interface{}{"forename": name, "lastname": surname, "email": email}
			body, _ := json.Marshal(postBody)
			r := httptest.NewRequest(http.MethodPost, usersEndPoint, bytes.NewReader(body))

			successResponse, errorResponse := api.CreateUserHandler(ctx, w, r)

			// Check whether testing a success or error case
			if tt.httpResponse > 399 {
				So(successResponse, ShouldBeNil)
				So(errorResponse.Status, ShouldEqual, tt.httpResponse)
			} else {
				So(successResponse.Status, ShouldEqual, tt.httpResponse)
				So(errorResponse, ShouldBeNil)
			}
		}
	})

	Convey("Admin create user returns 500: error unmarshalling request body", t, func() {
		r := httptest.NewRequest(http.MethodPost, usersEndPoint, bytes.NewReader(nil))

		successResponse, errorResponse := api.CreateUserHandler(ctx, w, r)

		So(successResponse, ShouldBeNil)
		castErr := errorResponse.Errors[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.JSONUnmarshalError)
		So(castErr.Description, ShouldEqual, models.ErrorUnmarshalFailedDescription)
	})

	Convey("Validation fails 400: validating email and username throws validation errors", t, func() {
		userValidationTests := []struct {
			userDetails  map[string]interface{}
			errorCodes   []string
			httpResponse int
		}{
			// missing email
			{
				map[string]interface{}{"forename": name, "lastname": surname, "email": ""},
				[]string{
					models.InvalidEmailError,
				},
				http.StatusBadRequest,
			},
			// missing both forename and surname
			{
				map[string]interface{}{"forename": "", "lastname": "", "email": email},
				[]string{
					models.InvalidForenameError,
					models.InvalidSurnameError,
				},
				http.StatusBadRequest,
			},
			// missing surname
			{
				map[string]interface{}{"forename": name, "lastname": "", "email": email},
				[]string{
					models.InvalidSurnameError,
				},
				http.StatusBadRequest,
			},
			// missing forename
			{
				map[string]interface{}{"forename": "", "lastname": surname, "email": email},
				[]string{
					models.InvalidForenameError,
				},
				http.StatusBadRequest,
			},
			// missing forename, surname and email
			{
				map[string]interface{}{"forename": "", "lastname": "", "email": ""},
				[]string{
					models.InvalidForenameError,
					models.InvalidSurnameError,
					models.InvalidEmailError,
				},
				http.StatusBadRequest,
			},
			// invalid email
			{
				map[string]interface{}{"forename": name, "lastname": surname, "email": invalidEmail},
				[]string{
					models.InvalidEmailError,
				},
				http.StatusBadRequest,
			},
		}

		for _, tt := range userValidationTests {
			body, _ := json.Marshal(tt.userDetails)
			r := httptest.NewRequest(http.MethodPost, usersEndPoint, bytes.NewReader(body))

			successResponse, errorResponse := api.CreateUserHandler(ctx, w, r)

			So(successResponse, ShouldBeNil)
			So(errorResponse.Status, ShouldEqual, tt.httpResponse)
			castErr := errorResponse.Errors[0].(*models.Error)
			So(castErr.Code, ShouldEqual, tt.errorCodes[0])
			if len(errorResponse.Errors) > 1 {
				castErr = errorResponse.Errors[1].(*models.Error)
				So(castErr.Code, ShouldEqual, tt.errorCodes[1])
			}
		}
	})
}

func TestListUserHandler(t *testing.T) {
	var ctx = context.Background()

	api, w, m := apiSetup()

	Convey("List user - check expected responses", t, func() {
		adminCreateUsersTests := []struct {
			listUsersFunction func(userInput *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error)
			httpResponse      int
		}{
			{
				// 200 response from Cognito
				func(userInput *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error) {
					users := &cognitoidentityprovider.ListUsersOutput{
						Users: []*cognitoidentityprovider.UserType{},
					}
					return users, nil
				},
				http.StatusOK,
			},
			{
				// 500 response from Cognito
				func(userInput *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error) {
					awsErrCode := "InternalErrorException"
					awsErrMessage := "Something strange happened"
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				http.StatusInternalServerError,
			},
		}

		for _, tt := range adminCreateUsersTests {
			m.ListUsersFunc = tt.listUsersFunction

			r := httptest.NewRequest(http.MethodGet, usersEndPoint, nil)

			successResponse, errorResponse := api.ListUsersHandler(ctx, w, r)

			// Check whether testing a success or error case
			if tt.httpResponse > 399 {
				So(successResponse, ShouldBeNil)
				So(errorResponse.Status, ShouldEqual, tt.httpResponse)
			} else {
				So(successResponse.Status, ShouldEqual, tt.httpResponse)
				So(errorResponse, ShouldBeNil)
			}
		}
	})
}

func TestGetUserHandler(t *testing.T) {
	var (
		ctx                                              = context.Background()
		forename, lastname, status, email, userId string = "bob", "bobbings", "UNCONFIRMED", "foo_bar123@ext.ons.gov.uk", "abcd1234"
		givenNameAttr, familyNameAttr, emailAttr  string = "given_name", "family_name", "email"
	)

	api, w, m := apiSetup()

	Convey("Get user - check expected responses", t, func() {
		adminGetUsersTests := []struct {
			getUserFunction func(userInput *cognitoidentityprovider.AdminGetUserInput) (*cognitoidentityprovider.AdminGetUserOutput, error)
			httpResponse    int
		}{
			{
				// 200 response from Cognito
				func(userInput *cognitoidentityprovider.AdminGetUserInput) (*cognitoidentityprovider.AdminGetUserOutput, error) {
					user := &cognitoidentityprovider.AdminGetUserOutput{
						UserAttributes: []*cognitoidentityprovider.AttributeType{
							{
								Name:  &givenNameAttr,
								Value: &forename,
							},
							{
								Name:  &familyNameAttr,
								Value: &lastname,
							},
							{
								Name:  &emailAttr,
								Value: &email,
							},
						},
						UserStatus: &status,
						Username:   &userId,
						Enabled:    aws.Bool(true),
					}
					return user, nil
				},
				http.StatusOK,
			},
			{
				// 500 response from Cognito
				func(userInput *cognitoidentityprovider.AdminGetUserInput) (*cognitoidentityprovider.AdminGetUserOutput, error) {
					awsErrCode := "InternalErrorException"
					awsErrMessage := "Something strange happened"
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				http.StatusInternalServerError,
			},
			{
				//404 response from Cognito
				func(userInput *cognitoidentityprovider.AdminGetUserInput) (*cognitoidentityprovider.AdminGetUserOutput, error) {
					awsErrCode := "UserNotFoundException"
					awsErrMessage := "user could not be found"
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				http.StatusNotFound,
			},
		}

		for _, tt := range adminGetUsersTests {
			m.AdminGetUserFunc = tt.getUserFunction

			r := httptest.NewRequest(http.MethodGet, userEndPoint, nil)

			successResponse, errorResponse := api.GetUserHandler(ctx, w, r)

			// Check whether testing a success or error case
			if tt.httpResponse > 399 {
				So(successResponse, ShouldBeNil)
				So(errorResponse.Status, ShouldEqual, tt.httpResponse)
			} else {
				So(successResponse.Status, ShouldEqual, tt.httpResponse)
				So(errorResponse, ShouldBeNil)
			}
		}
	})
}

func TestUpdateUserHandler(t *testing.T) {
	var (
		ctx                                              = context.Background()
		forename, lastname, email, userId, status string = "bob", "bobbings", "foo_bar123@ext.ons.gov.uk", "abcd1234", "CONFIRMED"
		givenNameAttr, familyNameAttr, emailAttr  string = "given_name", "family_name", "email"
	)

	api, w, m := apiSetup()

	successfullyGetUser := []*cognitoidentityprovider.AttributeType{
		{
			Name:  &givenNameAttr,
			Value: &forename,
		},
		{
			Name:  &familyNameAttr,
			Value: &lastname,
		},
		{
			Name:  &emailAttr,
			Value: &email,
		},
	}

	Convey("Update user - check expected responses", t, func() {
		adminCreateUsersTests := []struct {
			updateUserFunction  func(userInput *cognitoidentityprovider.AdminUpdateUserAttributesInput) (*cognitoidentityprovider.AdminUpdateUserAttributesOutput, error)
			getUserFunction     func(userInput *cognitoidentityprovider.AdminGetUserInput) (*cognitoidentityprovider.AdminGetUserOutput, error)
			enableUserFunction  func(userInput *cognitoidentityprovider.AdminEnableUserInput) (*cognitoidentityprovider.AdminEnableUserOutput, error)
			disableUserFunction func(userInput *cognitoidentityprovider.AdminDisableUserInput) (*cognitoidentityprovider.AdminDisableUserOutput, error)
			userForename        string
			userActive          bool
			assertions          func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse)
		}{
			// 200 response from Cognito
			{
				func(userInput *cognitoidentityprovider.AdminUpdateUserAttributesInput) (*cognitoidentityprovider.AdminUpdateUserAttributesOutput, error) {
					user := &cognitoidentityprovider.AdminUpdateUserAttributesOutput{}
					return user, nil
				},
				func(userInput *cognitoidentityprovider.AdminGetUserInput) (*cognitoidentityprovider.AdminGetUserOutput, error) {
					user := &cognitoidentityprovider.AdminGetUserOutput{
						UserAttributes: successfullyGetUser,
						UserStatus:     &status,
						Username:       &userId,
						Enabled:        aws.Bool(true),
					}
					return user, nil
				},
				func(userInput *cognitoidentityprovider.AdminEnableUserInput) (*cognitoidentityprovider.AdminEnableUserOutput, error) {
					return &cognitoidentityprovider.AdminEnableUserOutput{}, nil
				},
				func(userInput *cognitoidentityprovider.AdminDisableUserInput) (*cognitoidentityprovider.AdminDisableUserOutput, error) {
					return &cognitoidentityprovider.AdminDisableUserOutput{}, nil
				},
				forename,
				true,
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldNotBeNil)
					So(successResponse.Status, ShouldEqual, http.StatusOK)
					So(errorResponse, ShouldBeNil)
				},
			},
			//local validation failure
			{
				func(userInput *cognitoidentityprovider.AdminUpdateUserAttributesInput) (*cognitoidentityprovider.AdminUpdateUserAttributesOutput, error) {
					user := &cognitoidentityprovider.AdminUpdateUserAttributesOutput{}
					return user, nil
				},
				func(userInput *cognitoidentityprovider.AdminGetUserInput) (*cognitoidentityprovider.AdminGetUserOutput, error) {
					user := &cognitoidentityprovider.AdminGetUserOutput{
						UserAttributes: successfullyGetUser,
						UserStatus:     &status,
						Username:       &userId,
						Enabled:        aws.Bool(true),
					}
					return user, nil
				},
				func(userInput *cognitoidentityprovider.AdminEnableUserInput) (*cognitoidentityprovider.AdminEnableUserOutput, error) {
					return &cognitoidentityprovider.AdminEnableUserOutput{}, nil
				},
				func(userInput *cognitoidentityprovider.AdminDisableUserInput) (*cognitoidentityprovider.AdminDisableUserOutput, error) {
					return &cognitoidentityprovider.AdminDisableUserOutput{}, nil
				},
				"",
				true,
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
				},
			},
			//404 response from Cognito enable user
			{
				func(userInput *cognitoidentityprovider.AdminUpdateUserAttributesInput) (*cognitoidentityprovider.AdminUpdateUserAttributesOutput, error) {
					user := &cognitoidentityprovider.AdminUpdateUserAttributesOutput{}
					return user, nil
				},
				func(userInput *cognitoidentityprovider.AdminGetUserInput) (*cognitoidentityprovider.AdminGetUserOutput, error) {
					user := &cognitoidentityprovider.AdminGetUserOutput{
						UserAttributes: successfullyGetUser,
						UserStatus:     &status,
						Username:       &userId,
						Enabled:        aws.Bool(true),
					}
					return user, nil
				},
				func(userInput *cognitoidentityprovider.AdminEnableUserInput) (*cognitoidentityprovider.AdminEnableUserOutput, error) {
					awsErrCode := "UserNotFoundException"
					awsErrMessage := "user could not be found"
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				func(userInput *cognitoidentityprovider.AdminDisableUserInput) (*cognitoidentityprovider.AdminDisableUserOutput, error) {
					return &cognitoidentityprovider.AdminDisableUserOutput{}, nil
				},
				forename,
				true,
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusNotFound)
				},
			},
			// 500 response from Cognito enable user
			{
				func(userInput *cognitoidentityprovider.AdminUpdateUserAttributesInput) (*cognitoidentityprovider.AdminUpdateUserAttributesOutput, error) {
					user := &cognitoidentityprovider.AdminUpdateUserAttributesOutput{}
					return user, nil
				},
				func(userInput *cognitoidentityprovider.AdminGetUserInput) (*cognitoidentityprovider.AdminGetUserOutput, error) {
					user := &cognitoidentityprovider.AdminGetUserOutput{
						UserAttributes: successfullyGetUser,
						UserStatus:     &status,
						Username:       &userId,
						Enabled:        aws.Bool(true),
					}
					return user, nil
				},
				func(userInput *cognitoidentityprovider.AdminEnableUserInput) (*cognitoidentityprovider.AdminEnableUserOutput, error) {
					awsErrCode := "InternalErrorException"
					awsErrMessage := "Something strange happened"
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				func(userInput *cognitoidentityprovider.AdminDisableUserInput) (*cognitoidentityprovider.AdminDisableUserOutput, error) {
					return &cognitoidentityprovider.AdminDisableUserOutput{}, nil
				},
				forename,
				true,
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
				},
			},
			//404 response from Cognito disable user
			{
				func(userInput *cognitoidentityprovider.AdminUpdateUserAttributesInput) (*cognitoidentityprovider.AdminUpdateUserAttributesOutput, error) {
					user := &cognitoidentityprovider.AdminUpdateUserAttributesOutput{}
					return user, nil
				},
				func(userInput *cognitoidentityprovider.AdminGetUserInput) (*cognitoidentityprovider.AdminGetUserOutput, error) {
					user := &cognitoidentityprovider.AdminGetUserOutput{
						UserAttributes: successfullyGetUser,
						UserStatus:     &status,
						Username:       &userId,
						Enabled:        aws.Bool(true),
					}
					return user, nil
				},
				func(userInput *cognitoidentityprovider.AdminEnableUserInput) (*cognitoidentityprovider.AdminEnableUserOutput, error) {
					return &cognitoidentityprovider.AdminEnableUserOutput{}, nil
				},
				func(userInput *cognitoidentityprovider.AdminDisableUserInput) (*cognitoidentityprovider.AdminDisableUserOutput, error) {
					awsErrCode := "UserNotFoundException"
					awsErrMessage := "user could not be found"
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				forename,
				false,
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusNotFound)
				},
			},
			// 500 response from Cognito disable user
			{
				func(userInput *cognitoidentityprovider.AdminUpdateUserAttributesInput) (*cognitoidentityprovider.AdminUpdateUserAttributesOutput, error) {
					user := &cognitoidentityprovider.AdminUpdateUserAttributesOutput{}
					return user, nil
				},
				func(userInput *cognitoidentityprovider.AdminGetUserInput) (*cognitoidentityprovider.AdminGetUserOutput, error) {
					user := &cognitoidentityprovider.AdminGetUserOutput{
						UserAttributes: successfullyGetUser,
						UserStatus:     &status,
						Username:       &userId,
						Enabled:        aws.Bool(true),
					}
					return user, nil
				},
				func(userInput *cognitoidentityprovider.AdminEnableUserInput) (*cognitoidentityprovider.AdminEnableUserOutput, error) {
					return &cognitoidentityprovider.AdminEnableUserOutput{}, nil
				},
				func(userInput *cognitoidentityprovider.AdminDisableUserInput) (*cognitoidentityprovider.AdminDisableUserOutput, error) {
					awsErrCode := "InternalErrorException"
					awsErrMessage := "Something strange happened"
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				forename,
				false,
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
				},
			},
			//404 response from Cognito user update
			{
				func(userInput *cognitoidentityprovider.AdminUpdateUserAttributesInput) (*cognitoidentityprovider.AdminUpdateUserAttributesOutput, error) {
					awsErrCode := "UserNotFoundException"
					awsErrMessage := "user could not be found"
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				func(userInput *cognitoidentityprovider.AdminGetUserInput) (*cognitoidentityprovider.AdminGetUserOutput, error) {
					user := &cognitoidentityprovider.AdminGetUserOutput{
						UserAttributes: successfullyGetUser,
						UserStatus:     &status,
						Username:       &userId,
						Enabled:        aws.Bool(true),
					}
					return user, nil
				},
				func(userInput *cognitoidentityprovider.AdminEnableUserInput) (*cognitoidentityprovider.AdminEnableUserOutput, error) {
					return &cognitoidentityprovider.AdminEnableUserOutput{}, nil
				},
				func(userInput *cognitoidentityprovider.AdminDisableUserInput) (*cognitoidentityprovider.AdminDisableUserOutput, error) {
					return &cognitoidentityprovider.AdminDisableUserOutput{}, nil
				},
				forename,
				true,
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusNotFound)
				},
			},
			// 500 response from Cognito user update
			{
				func(userInput *cognitoidentityprovider.AdminUpdateUserAttributesInput) (*cognitoidentityprovider.AdminUpdateUserAttributesOutput, error) {
					awsErrCode := "InternalErrorException"
					awsErrMessage := "Something strange happened"
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				func(userInput *cognitoidentityprovider.AdminGetUserInput) (*cognitoidentityprovider.AdminGetUserOutput, error) {
					user := &cognitoidentityprovider.AdminGetUserOutput{
						UserAttributes: successfullyGetUser,
						UserStatus:     &status,
						Username:       &userId,
						Enabled:        aws.Bool(true),
					}
					return user, nil
				},
				func(userInput *cognitoidentityprovider.AdminEnableUserInput) (*cognitoidentityprovider.AdminEnableUserOutput, error) {
					return &cognitoidentityprovider.AdminEnableUserOutput{}, nil
				},
				func(userInput *cognitoidentityprovider.AdminDisableUserInput) (*cognitoidentityprovider.AdminDisableUserOutput, error) {
					return &cognitoidentityprovider.AdminDisableUserOutput{}, nil
				},
				forename,
				true,
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
				},
			},
			//reload user details failure
			{
				func(userInput *cognitoidentityprovider.AdminUpdateUserAttributesInput) (*cognitoidentityprovider.AdminUpdateUserAttributesOutput, error) {
					user := &cognitoidentityprovider.AdminUpdateUserAttributesOutput{}
					return user, nil
				},
				func(userInput *cognitoidentityprovider.AdminGetUserInput) (*cognitoidentityprovider.AdminGetUserOutput, error) {
					awsErrCode := "UserNotFoundException"
					awsErrMessage := "user could not be found"
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				func(userInput *cognitoidentityprovider.AdminEnableUserInput) (*cognitoidentityprovider.AdminEnableUserOutput, error) {
					return &cognitoidentityprovider.AdminEnableUserOutput{}, nil
				},
				func(userInput *cognitoidentityprovider.AdminDisableUserInput) (*cognitoidentityprovider.AdminDisableUserOutput, error) {
					return &cognitoidentityprovider.AdminDisableUserOutput{}, nil
				},
				forename,
				true,
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
				},
			},
		}

		for _, tt := range adminCreateUsersTests {
			m.AdminUpdateUserAttributesFunc = tt.updateUserFunction
			m.AdminGetUserFunc = tt.getUserFunction
			m.AdminEnableUserFunc = tt.enableUserFunction
			m.AdminDisableUserFunc = tt.disableUserFunction

			postBody := map[string]interface{}{"forename": tt.userForename, "lastname": lastname, "active": tt.userActive}
			body, err := json.Marshal(postBody)

			So(err, ShouldBeNil)

			r := httptest.NewRequest(http.MethodGet, userEndPoint, bytes.NewReader(body))

			successResponse, errorResponse := api.UpdateUserHandler(ctx, w, r)

			tt.assertions(successResponse, errorResponse)
		}
	})
}

func TestProcessUpdateCognitoError(t *testing.T) {
	ctx := context.Background()

	Convey("Processes UserNotFound to a 404 response", t, func() {
		awsErrCode := "UserNotFoundException"
		awsErrMessage := "user could not be found"
		awsOrigErr := errors.New(awsErrCode)
		awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
		errResponse := processUpdateCognitoError(ctx, awsErr, "Testing user not found")
		So(errResponse.Status, ShouldEqual, http.StatusNotFound)
		castErr := errResponse.Errors[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.UserNotFoundError)
		So(castErr.Description, ShouldEqual, "user could not be found")
	})

	Convey("Processes InternalError to a 500 response", t, func() {
		awsErrCode := "InternalErrorException"
		awsErrMessage := "something went wrong"
		awsOrigErr := errors.New(awsErrCode)
		awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
		errResponse := processUpdateCognitoError(ctx, awsErr, "Testing internal error")
		So(errResponse.Status, ShouldEqual, http.StatusInternalServerError)
		castErr := errResponse.Errors[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InternalError)
		So(castErr.Description, ShouldEqual, "something went wrong")
	})

	Convey("Processes InvalidField to a 400 response", t, func() {
		awsErrCode := "InvalidParameterException"
		awsErrMessage := "param invalid"
		awsOrigErr := errors.New(awsErrCode)
		awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
		errResponse := processUpdateCognitoError(ctx, awsErr, "Testing invalid param error")
		So(errResponse.Status, ShouldEqual, http.StatusBadRequest)
		castErr := errResponse.Errors[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidFieldError)
		So(castErr.Description, ShouldEqual, "param invalid")
	})
}

func TestChangePasswordHandler(t *testing.T) {

	var (
		ctx                                       = context.Background()
		email, password, session           string = "foo_bar123@ext.ons.gov.uk", "Password2", "auth-challenge-session"
		accessToken, idToken, refreshToken string = "aaaa.bbbb.cccc", "llll.mmmm.nnnn", "zzzz.yyyy.xxxx.wwww.vvvv"
		expireLength                       int64  = 500
	)

	api, w, m := apiSetup()

	Convey("RespondToAuthChallenge - check expected responses", t, func() {
		respondToAuthChallengeTests := []struct {
			respondToAuthChallengeFunction func(input *cognitoidentityprovider.RespondToAuthChallengeInput) (*cognitoidentityprovider.RespondToAuthChallengeOutput, error)
			httpResponse                   int
		}{
			{
				// Cognito successful password change
				func(input *cognitoidentityprovider.RespondToAuthChallengeInput) (*cognitoidentityprovider.RespondToAuthChallengeOutput, error) {
					return &cognitoidentityprovider.RespondToAuthChallengeOutput{
						AuthenticationResult: &cognitoidentityprovider.AuthenticationResultType{
							AccessToken:  &accessToken,
							ExpiresIn:    &expireLength,
							IdToken:      &idToken,
							RefreshToken: &refreshToken,
						},
					}, nil
				},
				http.StatusAccepted,
			},
			{
				// Cognito internal error
				func(input *cognitoidentityprovider.RespondToAuthChallengeInput) (*cognitoidentityprovider.RespondToAuthChallengeOutput, error) {
					awsErrCode := "InternalErrorException"
					awsErrMessage := "Something strange happened"
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				http.StatusInternalServerError,
			},
			{
				// Cognito invalid session
				func(input *cognitoidentityprovider.RespondToAuthChallengeInput) (*cognitoidentityprovider.RespondToAuthChallengeOutput, error) {
					awsErrCode := "CodeMismatchException"
					awsErrMessage := "session invalid"
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				http.StatusBadRequest,
			},
			{
				// Cognito invalid password
				func(input *cognitoidentityprovider.RespondToAuthChallengeInput) (*cognitoidentityprovider.RespondToAuthChallengeOutput, error) {
					awsErrCode := "InvalidPasswordException"
					awsErrMessage := "password invalid"
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				http.StatusBadRequest,
			},
			{
				// Cognito invalid user
				func(input *cognitoidentityprovider.RespondToAuthChallengeInput) (*cognitoidentityprovider.RespondToAuthChallengeOutput, error) {
					awsErrCode := "UserNotFoundException"
					awsErrMessage := "user not found"
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				http.StatusAccepted,
			},
		}

		for _, tt := range respondToAuthChallengeTests {
			m.RespondToAuthChallengeFunc = tt.respondToAuthChallengeFunction

			postBody := map[string]interface{}{"type": models.NewPasswordRequiredType, "email": email, "password": password, "session": session}
			body, _ := json.Marshal(postBody)
			r := httptest.NewRequest(http.MethodPut, changePasswordEndPoint, bytes.NewReader(body))

			successResponse, errorResponse := api.ChangePasswordHandler(ctx, w, r)

			// Check whether testing a success or error case
			if tt.httpResponse > 399 {
				So(successResponse, ShouldBeNil)
				So(errorResponse.Status, ShouldEqual, tt.httpResponse)
			} else {
				So(successResponse.Status, ShouldEqual, tt.httpResponse)
				So(errorResponse, ShouldBeNil)
			}
		}
	})
}

func TestConfirmForgotPasswordChangePasswordHandler(t *testing.T) {

	var (
		ctx                      = context.Background()
		email             string = "fred.bloggs@ons.gov.uk"
		password          string = "Password2@123456"
		verificationToken string = "999999"
	)

	api, w, m := apiSetup()
	Convey("ConfirmForgotPassword - check expected responses", t, func() {
		confirmForgotPasswordTests := []struct {
			confirmForgotPasswordFunction func(input *cognitoidentityprovider.ConfirmForgotPasswordInput) (*cognitoidentityprovider.ConfirmForgotPasswordOutput, error)
			httpResponse                  int
		}{
			{
				// Cognito successful password change
				func(input *cognitoidentityprovider.ConfirmForgotPasswordInput) (*cognitoidentityprovider.ConfirmForgotPasswordOutput, error) {
					tst := cognitoidentityprovider.ConfirmForgotPasswordOutput{}
					return &tst, nil
				},
				http.StatusAccepted,
			},
			{
				// Cognito internal error
				func(input *cognitoidentityprovider.ConfirmForgotPasswordInput) (*cognitoidentityprovider.ConfirmForgotPasswordOutput, error) {
					awsErrCode := "InternalErrorException"
					awsErrMessage := "Something strange happened"
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				http.StatusInternalServerError,
			},
			{
				// Cognito invalid session
				func(input *cognitoidentityprovider.ConfirmForgotPasswordInput) (*cognitoidentityprovider.ConfirmForgotPasswordOutput, error) {
					awsErrCode := "CodeMismatchException"
					awsErrMessage := "session invalid"
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				http.StatusBadRequest,
			},
			{
				// Cognito invalid password
				func(input *cognitoidentityprovider.ConfirmForgotPasswordInput) (*cognitoidentityprovider.ConfirmForgotPasswordOutput, error) {
					awsErrCode := "InvalidPasswordException"
					awsErrMessage := "password invalid"
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				http.StatusBadRequest,
			},
			{
				// Cognito invalid user
				func(input *cognitoidentityprovider.ConfirmForgotPasswordInput) (*cognitoidentityprovider.ConfirmForgotPasswordOutput, error) {
					awsErrCode := "UserNotFoundException"
					awsErrMessage := "user not found"
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				http.StatusAccepted,
			},
		}

		for _, tt := range confirmForgotPasswordTests {
			m.ConfirmForgotPasswordFunc = tt.confirmForgotPasswordFunction

			postBody := map[string]interface{}{"type": models.ForgottenPasswordType, "email": email, "password": password, "verification_token": verificationToken}
			body, _ := json.Marshal(postBody)
			r := httptest.NewRequest(http.MethodPut, changePasswordEndPoint, bytes.NewReader(body))

			successResponse, errorResponse := api.ChangePasswordHandler(ctx, w, r)

			// Check whether testing a success or error case
			if tt.httpResponse > 399 {
				So(successResponse, ShouldBeNil)
				So(errorResponse.Status, ShouldEqual, tt.httpResponse)
			} else {
				So(successResponse.Status, ShouldEqual, tt.httpResponse)
				So(errorResponse, ShouldBeNil)
			}
		}
	})

	Convey("ConfirmForgotPassword returns 500: error unmarshalling request body", t, func() {
		r := httptest.NewRequest(http.MethodPut, changePasswordEndPoint, bytes.NewReader(nil))

		successResponse, errorResponse := api.CreateUserHandler(ctx, w, r)

		So(successResponse, ShouldBeNil)
		castErr := errorResponse.Errors[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.JSONUnmarshalError)
		So(castErr.Description, ShouldEqual, models.ErrorUnmarshalFailedDescription)
	})

	Convey("Validation fails 400: validation of a required param throws validation errors", t, func() {
		validationTests := []struct {
			requestBody  map[string]interface{}
			errorCode    string
			httpResponse int
		}{
			// missing password change type
			{
				map[string]interface{}{"type": "", "email": email, "password": password, "verification_token": verificationToken},
				models.UnknownRequestTypeError,
				http.StatusBadRequest,
			},
			// missing a change request param
			{
				map[string]interface{}{"type": models.ForgottenPasswordType, "email": "", "password": password, "verification_token": verificationToken},
				models.InvalidEmailError,
				http.StatusBadRequest,
			},
		}

		for _, tt := range validationTests {
			body, _ := json.Marshal(tt.requestBody)
			r := httptest.NewRequest(http.MethodPut, changePasswordEndPoint, bytes.NewReader(body))

			successResponse, errorResponse := api.ChangePasswordHandler(ctx, w, r)

			So(successResponse, ShouldBeNil)
			So(errorResponse.Status, ShouldEqual, tt.httpResponse)
			So(len(errorResponse.Errors), ShouldEqual, 1)
			castErr := errorResponse.Errors[0].(*models.Error)
			So(castErr.Code, ShouldEqual, tt.errorCode)
		}
	})
}

func TestPasswordResetHandler(t *testing.T) {

	var (
		ctx          = context.Background()
		email string = "foo_bar123@ext.ons.gov.uk"
	)

	api, w, m := apiSetup()

	Convey("ForgotPassword - check expected responses", t, func() {
		respondToAuthChallengeTests := []struct {
			forgotPasswordFunction func(input *cognitoidentityprovider.ForgotPasswordInput) (*cognitoidentityprovider.ForgotPasswordOutput, error)
			httpResponse           int
		}{
			{
				// Cognito successful password change
				func(input *cognitoidentityprovider.ForgotPasswordInput) (*cognitoidentityprovider.ForgotPasswordOutput, error) {
					return &cognitoidentityprovider.ForgotPasswordOutput{
						CodeDeliveryDetails: &cognitoidentityprovider.CodeDeliveryDetailsType{},
					}, nil
				},
				http.StatusAccepted,
			},
			{
				// Cognito internal error
				func(input *cognitoidentityprovider.ForgotPasswordInput) (*cognitoidentityprovider.ForgotPasswordOutput, error) {
					awsErrCode := "InternalErrorException"
					awsErrMessage := "Something strange happened"
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				http.StatusInternalServerError,
			},
			{
				// Cognito too many requests
				func(input *cognitoidentityprovider.ForgotPasswordInput) (*cognitoidentityprovider.ForgotPasswordOutput, error) {
					awsErrCode := "TooManyRequestsException"
					awsErrMessage := "slow down"
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				http.StatusBadRequest,
			},
			{
				// Cognito invalid user
				func(input *cognitoidentityprovider.ForgotPasswordInput) (*cognitoidentityprovider.ForgotPasswordOutput, error) {
					awsErrCode := "UserNotFoundException"
					awsErrMessage := "user not found in user pool"
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				http.StatusAccepted,
			},
		}

		for _, tt := range respondToAuthChallengeTests {
			m.ForgotPasswordFunc = tt.forgotPasswordFunction

			postBody := map[string]interface{}{"email": email}
			body, _ := json.Marshal(postBody)
			r := httptest.NewRequest(http.MethodPost, requestResetEndPoint, bytes.NewReader(body))

			successResponse, errorResponse := api.PasswordResetHandler(ctx, w, r)

			// Check whether testing a success or error case
			if tt.httpResponse > 399 {
				So(successResponse, ShouldBeNil)
				So(errorResponse.Status, ShouldEqual, tt.httpResponse)
			} else {
				So(successResponse.Status, ShouldEqual, tt.httpResponse)
				So(errorResponse, ShouldBeNil)
			}
		}
	})

	Convey("ForgotPassword returns 500: error unmarshalling request body", t, func() {
		r := httptest.NewRequest(http.MethodPost, requestResetEndPoint, bytes.NewReader(nil))

		successResponse, errorResponse := api.PasswordResetHandler(ctx, w, r)

		So(successResponse, ShouldBeNil)
		castErr := errorResponse.Errors[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.JSONUnmarshalError)
		So(castErr.Description, ShouldEqual, models.ErrorUnmarshalFailedDescription)
	})

	Convey("Validation fails 400: validation of a required param throws validation errors", t, func() {
		validationTests := []struct {
			requestBody  map[string]interface{}
			errorCode    string
			httpResponse int
		}{
			// missing a change request param
			{
				map[string]interface{}{"email": ""},
				models.InvalidEmailError,
				http.StatusBadRequest,
			},
		}

		for _, tt := range validationTests {
			body, _ := json.Marshal(tt.requestBody)
			r := httptest.NewRequest(http.MethodPost, requestResetEndPoint, bytes.NewReader(body))

			successResponse, errorResponse := api.PasswordResetHandler(ctx, w, r)

			So(successResponse, ShouldBeNil)
			So(errorResponse.Status, ShouldEqual, tt.httpResponse)
			So(len(errorResponse.Errors), ShouldEqual, 1)
			castErr := errorResponse.Errors[0].(*models.Error)
			So(castErr.Code, ShouldEqual, tt.errorCode)
		}
	})
}

func TestListUserGroupsHandler(t *testing.T) {

	var (
		ctx       = context.Background()
		timestamp = time.Now()
		// nextToken  = "abc1234"
		groups = []*cognitoidentityprovider.GroupType{
			{
				CreationDate:     &timestamp,
				Description:      aws.String("A test group1"),
				GroupName:        aws.String("test-group1"),
				LastModifiedDate: &timestamp,
				Precedence:       aws.Int64(4),
				RoleArn:          aws.String(""),
				UserPoolId:       aws.String(""),
			},
			{
				CreationDate:     &timestamp,
				Description:      aws.String("A test group1"),
				GroupName:        aws.String("test-group1"),
				LastModifiedDate: &timestamp,
				Precedence:       aws.Int64(4),
				RoleArn:          aws.String(""),
				UserPoolId:       aws.String(""),
			},
		}
	)

	api, w, m := apiSetup()

	Convey("List groups for user -check expected responses", t, func() {
		listusergroups := []struct {
			getUserGroupsFunction func(input *cognitoidentityprovider.AdminListGroupsForUserInput) (*cognitoidentityprovider.AdminListGroupsForUserOutput, error)
			httpResponse          int
		}{
			{
				// 200 response from Cognito with empty NextToken
				func(input *cognitoidentityprovider.AdminListGroupsForUserInput) (*cognitoidentityprovider.AdminListGroupsForUserOutput, error) {
					return &cognitoidentityprovider.AdminListGroupsForUserOutput{
						Groups:    groups,
						NextToken: nil,
					}, nil
				},
				http.StatusOK,
			},
			{
				// 200 response from Cognito with empty NextToken
				func(input *cognitoidentityprovider.AdminListGroupsForUserInput) (*cognitoidentityprovider.AdminListGroupsForUserOutput, error) {
					return &cognitoidentityprovider.AdminListGroupsForUserOutput{
						Groups:    []*cognitoidentityprovider.GroupType{},
						NextToken: nil,
					}, nil
				},
				http.StatusOK,
			},
			{
				// 500 response from Cognito
				func(input *cognitoidentityprovider.AdminListGroupsForUserInput) (*cognitoidentityprovider.AdminListGroupsForUserOutput, error) {
					awsErrCode := "InternalErrorException"
					awsErrMessage := "Something strange happened"
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				http.StatusInternalServerError,
			},
			{
				//404 response from Cognito
				func(input *cognitoidentityprovider.AdminListGroupsForUserInput) (*cognitoidentityprovider.AdminListGroupsForUserOutput, error) {
					awsErrCode := "UserNotFoundException"
					awsErrMessage := "user could not be found"
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				http.StatusInternalServerError,
			},
		}

		for _, tt := range listusergroups {
			m.ListGroupsForUserFunc = tt.getUserGroupsFunction

			r := httptest.NewRequest(http.MethodGet, userListGroupsEndPoint, nil)

			urlVars := map[string]string{
				"id": "efgh5678",
			}
			r = mux.SetURLVars(r, urlVars)

			successResponse, errorResponse := api.ListUserGroupsHandler(ctx, w, r)

			// Check whether testing a success or error case
			if tt.httpResponse > 399 {
				So(successResponse, ShouldBeNil)
				So(errorResponse.Status, ShouldEqual, tt.httpResponse)
			} else {
				So(successResponse.Status, ShouldEqual, tt.httpResponse)
				So(errorResponse, ShouldBeNil)
			}
		}
	})
}

func TestGetGroupsforUser(t *testing.T) {

	var (
		userNotFoundDescription string = "User not found"
		userId                         = models.UserParams{
			ID: "abcd1234",
		}
		group_0 string = "test_group_0"
		group_1 string = "test_group_1"
	)

	listOfGroups := []*cognitoidentityprovider.GroupType{
		{
			GroupName: &group_0,
		},
	}

	api, _, m := apiSetup()
	Convey("error is returned when list groups for a user returns an error", t, func() {
		m.ListGroupsForUserFunc = func(input *cognitoidentityprovider.AdminListGroupsForUserInput) (*cognitoidentityprovider.AdminListGroupsForUserOutput, error) {
			var userNotFoundException cognitoidentityprovider.ResourceNotFoundException
			userNotFoundException.Message_ = &userNotFoundDescription
			return nil, &userNotFoundException
		}

		listGroupsforUserResponse, errorResponse := api.getGroupsForUser(nil, userId)

		So(listGroupsforUserResponse, ShouldBeNil)
		So(errorResponse.Error(), ShouldResemble, "ResourceNotFoundException: User not found")
	})

	Convey("When there is no next token cognito is called once and the list of groups in returned", t, func() {
		listOfGroupsForUser := []*cognitoidentityprovider.GroupType{
			{
				GroupName: &group_0,
			},
		}

		m.ListGroupsForUserFunc = func(input *cognitoidentityprovider.AdminListGroupsForUserInput) (*cognitoidentityprovider.AdminListGroupsForUserOutput, error) {

			listGroupsForUser := &cognitoidentityprovider.AdminListGroupsForUserOutput{
				Groups: []*cognitoidentityprovider.GroupType{
					{
						GroupName: &group_0,
					},
				},
			}
			return listGroupsForUser, nil
		}

		listOfUsersResponse, errorResponse := api.getGroupsForUser(nil, userId)

		So(listOfUsersResponse, ShouldResemble, listOfGroupsForUser)

		So(errorResponse, ShouldBeNil)

	})

	Convey("When there is a next token cognito is called more than once and the appended list of users in returned", t, func() {
		listOfGroupsForUser := []*cognitoidentityprovider.GroupType{
			{
				GroupName: &group_0,
			},
			{
				GroupName: &group_0,
			},
			{
				GroupName: &group_1,
			},
		}

		m.ListGroupsForUserFunc = func(input *cognitoidentityprovider.AdminListGroupsForUserInput) (*cognitoidentityprovider.AdminListGroupsForUserOutput, error) {
			nextToken := "nextToken"

			if input.NextToken != nil {
				listGroupsForUser := &cognitoidentityprovider.AdminListGroupsForUserOutput{
					NextToken: nil,
					Groups: []*cognitoidentityprovider.GroupType{
						{
							GroupName: &group_1,
						},
					},
				}
				return listGroupsForUser, nil
			} else {
				listGroupsForUser := &cognitoidentityprovider.AdminListGroupsForUserOutput{
					NextToken: &nextToken,
					Groups: []*cognitoidentityprovider.GroupType{
						{
							GroupName: &group_0,
						},
					},
				}
				return listGroupsForUser, nil
			}
		}

		listGroupsForUserResponse, errorResponse := api.getGroupsForUser(listOfGroups, userId)

		So(listGroupsForUserResponse, ShouldResemble, listOfGroupsForUser)
		So(errorResponse, ShouldBeNil)

	})

	Convey("When GetGroupsforUser in called with a list of groups the appended list of groups in returned", t, func() {

		listOfGroups := []*cognitoidentityprovider.GroupType{
			{
				GroupName: &group_0,
			},
		}

		returnedlistOfGroups := []*cognitoidentityprovider.GroupType{
			{
				GroupName: &group_0,
			},
			{
				GroupName: &group_0,
			},
		}

		m.ListGroupsForUserFunc = func(input *cognitoidentityprovider.AdminListGroupsForUserInput) (*cognitoidentityprovider.AdminListGroupsForUserOutput, error) {
			listGroupsForUser := &cognitoidentityprovider.AdminListGroupsForUserOutput{
				Groups: []*cognitoidentityprovider.GroupType{
					{
						GroupName: &group_0,
					},
				},
			}
			return listGroupsForUser, nil
		}

		listGroupsForUseResponse, errorResponse := api.getGroupsForUser(listOfGroups, userId)

		So(listGroupsForUseResponse, ShouldResemble, returnedlistOfGroups)
		So(errorResponse, ShouldBeNil)
	})
}
