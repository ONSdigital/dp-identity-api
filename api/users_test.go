package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	. "github.com/smartystreets/goconvey/convey"
)

const usersEndPoint = "http://localhost:25600/v1/users"
const userEndPoint = "http://localhost:25600/v1/users/abcd1234"
const changePasswordEndPoint = "http://localhost:25600/v1/users/self/password"
const requestResetEndPoint = "http://localhost:25600/v1/password-reset"

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
			//So(len(errorResponse.Errors), ShouldEqual, len(tt.errorCodes))
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
		adminCreateUsersTests := []struct {
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

		for _, tt := range adminCreateUsersTests {
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

	Convey("RespondToAuthChallenge returns 500: error unmarshalling request body", t, func() {
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
				map[string]interface{}{"type": "", "email": email, "password": password, "session": session},
				models.UnknownRequestTypeError,
				http.StatusBadRequest,
			},
			// missing a change request param
			{
				map[string]interface{}{"type": models.NewPasswordRequiredType, "email": "", "password": password, "session": session},
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
