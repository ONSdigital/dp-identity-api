package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ONSdigital/dp-identity-api/v2/cognito/mock"
	"github.com/ONSdigital/dp-identity-api/v2/models"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	. "github.com/smartystreets/goconvey/convey"
)

const signInEndPoint = "http://localhost:25600/v1/tokens"
const signOutEndPoint = "http://localhost:25600/v1/tokens/self"
const tokenRefreshEndPoint = "http://localhost:25600/v1/tokens/self" // #nosec

func TestAPI_TokensHandler(t *testing.T) {
	var (
		ctx                                      = context.Background()
		accessToken, idToken, refreshToken       = "aaaa.bbbb.cccc", "llll.mmmm.nnnn", "zzzz.yyyy.xxxx.wwww.vvvv"
		expireLength                       int64 = 500
	)

	api, w, m := apiSetup()

	// mock call to: AdminUserGlobalSignOut(input *cognitoidentityprovider.AdminUserGlobalSignOutInput) (*cognitoidentityprovider.AdminUserGlobalSignOutOutput, error)
	m.AdminUserGlobalSignOutFunc = func(_ *cognitoidentityprovider.AdminUserGlobalSignOutInput) (*cognitoidentityprovider.AdminUserGlobalSignOutOutput, error) {
		return &cognitoidentityprovider.AdminUserGlobalSignOutOutput{}, nil
	}
	// mock call to: InitiateAuth(input *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error)
	m.InitiateAuthFunc = func(_ *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error) {
		return &cognitoidentityprovider.InitiateAuthOutput{
			AuthenticationResult: &cognitoidentityprovider.AuthenticationResultType{
				AccessToken:  &accessToken,
				ExpiresIn:    &expireLength,
				IdToken:      &idToken,
				RefreshToken: &refreshToken,
			},
		}, nil
	}
	m.DescribeUserPoolClientFunc = func(_ *cognitoidentityprovider.DescribeUserPoolClientInput) (*cognitoidentityprovider.DescribeUserPoolClientOutput, error) {
		tokenValidDays := int64(1)
		refreshTokenUnits := cognitoidentityprovider.TimeUnitsTypeDays

		userPoolClient := &cognitoidentityprovider.DescribeUserPoolClientOutput{
			UserPoolClient: &cognitoidentityprovider.UserPoolClientType{
				RefreshTokenValidity: &tokenValidDays,
				TokenValidityUnits: &cognitoidentityprovider.TokenValidityUnitsType{
					RefreshToken: &refreshTokenUnits,
				},
			},
		}
		return userPoolClient, nil
	}

	Convey("Sign in success: no ErrorResponse, SuccessResponse Status 201", t, func() {
		body := map[string]interface{}{
			"email":    "email@ons.gov.uk",
			"password": "password",
		}
		jsonBody, err := json.Marshal(&body)
		So(err, ShouldBeNil)
		request := httptest.NewRequest(http.MethodPost, signInEndPoint, bytes.NewBuffer(jsonBody))

		successResponse, errorResponse := api.TokensHandler(ctx, w, request)

		So(errorResponse, ShouldBeNil)
		So(successResponse.Status, ShouldEqual, http.StatusCreated)
		var responseBody map[string]interface{}
		err = json.Unmarshal(successResponse.Body, &responseBody)
		So(err, ShouldBeNil)
		So(responseBody["expirationTime"], ShouldNotBeNil)
		So(responseBody["refreshTokenExpirationTime"], ShouldNotBeNil)
	})

	Convey("Sign In validation error: adds an error to the ErrorResponse and sets its Status to 400", t, func() {
		body := map[string]interface{}{
			"email":    "email@ons.gov.uk",
			"password": "",
		}
		jsonBody, err := json.Marshal(&body)
		So(err, ShouldBeNil)
		request := httptest.NewRequest(http.MethodPost, signInEndPoint, bytes.NewBuffer(jsonBody))

		successResponse, errorResponse := api.TokensHandler(ctx, w, request)

		So(successResponse, ShouldBeNil)
		So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
		So(len(errorResponse.Errors), ShouldEqual, 1)
		So(errorResponse.Errors[0].Error(), ShouldEqual, models.InvalidPasswordError)
	})

	Convey("Sign In Cognito internal error: adds an error to the ErrorResponse and sets its Status to 500", t, func() {
		awsOrigErr := errors.New(awsErrCode)
		awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
		// mock failed call to: InitiateAuth(input *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error)
		m.InitiateAuthFunc = func(_ *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error) {
			return nil, awsErr
		}

		body := map[string]interface{}{
			"email":    "email@ons.gov.uk",
			"password": "password",
		}
		jsonBody, err := json.Marshal(&body)
		So(err, ShouldBeNil)
		request := httptest.NewRequest(http.MethodPost, signInEndPoint, bytes.NewBuffer(jsonBody))

		successResponse, errorResponse := api.TokensHandler(ctx, w, request)

		So(successResponse, ShouldBeNil)
		So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
		So(len(errorResponse.Errors), ShouldEqual, 1)
		So(errorResponse.Errors[0].Error(), ShouldEqual, awsErr.Error())
	})

	Convey("Sign In Cognito request error: adds an error to the ErrorResponse and sets the Status correctly", t, func() {
		statusTests := []struct {
			awsErrCode       string
			awsErrMessage    string
			httpResponseCode int
		}{
			// http.StatusBadRequest - 400
			{
				"NotAuthorizedException",
				"User is not authorized",
				http.StatusBadRequest,
			},
			// http.StatusUnauthorized - 401
			{
				"NotAuthorizedException",
				"Incorrect username or password.",
				http.StatusUnauthorized,
			},
		}

		for _, tt := range statusTests {
			awsErr := awserr.New(tt.awsErrCode, tt.awsErrMessage, errors.New(tt.awsErrCode))
			// mock failed call to: InitiateAuth(input *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error)
			m.InitiateAuthFunc = func(_ *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error) {
				return nil, awsErr
			}

			body := map[string]interface{}{
				"email":    "email@ons.gov.uk",
				"password": "password",
			}
			jsonBody, err := json.Marshal(&body)
			So(err, ShouldBeNil)
			request := httptest.NewRequest(http.MethodPost, signInEndPoint, bytes.NewBuffer(jsonBody))

			successResponse, errorResponse := api.TokensHandler(ctx, w, request)

			request.Header.Get(WWWAuthenticateName)

			So(successResponse, ShouldBeNil)
			So(errorResponse.Status, ShouldEqual, tt.httpResponseCode)
			So(len(errorResponse.Errors), ShouldEqual, 1)
			So(errorResponse.Errors[0].Error(), ShouldEqual, awsErr.Error())
		}
	})

	// test Tokens handler's NEW_PASSWORD_REQUIRED challenge response
	Convey("Handle NEW_PASSWORD_REQUIRED challenge response", t, func() {
		var (
			newPasswordStatus, sessionID = "true", "AYABeBBsY5be-this-is-a-test-session-id-string-123456789iuerhcfdisieo-end"
		)

		// mock call to: InitiateAuth(input *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error)
		m.InitiateAuthFunc = func(_ *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error) {
			challengeName := "NEW_PASSWORD_REQUIRED"
			return &cognitoidentityprovider.InitiateAuthOutput{
				AuthenticationResult: nil,
				ChallengeName:        &challengeName,
				Session:              &sessionID,
			}, nil
		}

		body := map[string]interface{}{
			"email":    "email@ons.gov.uk",
			"password": "password",
		}
		jsonBody, err := json.Marshal(&body)
		So(err, ShouldBeNil)
		request := httptest.NewRequest(http.MethodPost, signInEndPoint, bytes.NewBuffer(jsonBody))

		successResponse, errorResponse := api.TokensHandler(ctx, w, request)

		So(errorResponse, ShouldBeNil)
		So(successResponse.Status, ShouldEqual, http.StatusAccepted)
		var responseBody map[string]interface{}
		err = json.Unmarshal(successResponse.Body, &responseBody)
		So(err, ShouldBeNil)
		So(responseBody["new_password_required"], ShouldEqual, newPasswordStatus)
		So(responseBody["session"], ShouldEqual, sessionID)
	})
}

func TestAPI_SignOutHandler(t *testing.T) {
	var ctx = context.Background()

	api, w, m := apiSetup()

	// mock call to: GlobalSignOut(input *cognitoidentityprovider.GlobalSignOutInput) (*cognitoidentityprovider.GlobalSignOutOutput, error)
	m.GlobalSignOutFunc = func(_ *cognitoidentityprovider.GlobalSignOutInput) (*cognitoidentityprovider.GlobalSignOutOutput, error) {
		return &cognitoidentityprovider.GlobalSignOutOutput{}, nil
	}

	Convey("Global Sign Out success: no errors added to ErrorResponse Errors list", t, func() {
		request := httptest.NewRequest(http.MethodDelete, signOutEndPoint, http.NoBody)
		request.Header.Set(AccessTokenHeaderName, "Bearer zzzz-yyyy-xxxx")

		successResponse, errorResponse := api.SignOutHandler(ctx, w, request)

		So(errorResponse, ShouldBeNil)
		So(successResponse.Status, ShouldEqual, http.StatusNoContent)
		So(successResponse.Body, ShouldBeNil)
	})

	Convey("Global Sign Out validation error: adds an error to the ErrorResponse and sets its Status to 400", t, func() {
		request := httptest.NewRequest(http.MethodDelete, signOutEndPoint, http.NoBody)
		request.Header.Set(AccessTokenHeaderName, "")

		successResponse, errorResponse := api.SignOutHandler(ctx, w, request)

		So(successResponse, ShouldBeNil)
		So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
		So(len(errorResponse.Errors), ShouldEqual, 1)
		So(errorResponse.Errors[0].Error(), ShouldEqual, models.InvalidTokenError)
	})

	Convey("Global Sign Out Cognito internal error: adds an error to the ErrorResponse and sets its Status to 500", t, func() {
		awsOrigErr := errors.New(awsErrCode)
		awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
		// mock failed call to: GlobalSignOut(input *cognitoidentityprovider.GlobalSignOutInput) (*cognitoidentityprovider.GlobalSignOutOutput, error)
		m.GlobalSignOutFunc = func(_ *cognitoidentityprovider.GlobalSignOutInput) (*cognitoidentityprovider.GlobalSignOutOutput, error) {
			return nil, awsErr
		}

		request := httptest.NewRequest(http.MethodDelete, signOutEndPoint, http.NoBody)
		request.Header.Set(AccessTokenHeaderName, "Bearer zzzz-yyyy-xxxx")

		successResponse, errorResponse := api.SignOutHandler(ctx, w, request)

		So(successResponse, ShouldBeNil)
		So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
		So(len(errorResponse.Errors), ShouldEqual, 1)
		So(errorResponse.Errors[0].Error(), ShouldEqual, awsErr.Error())
	})

	Convey("Global Sign Out Cognito request error: adds an error to the ErrorResponse and sets its Status to 400", t, func() {
		awsErrCode := "NotAuthorizedException"
		awsErrMessage := "User is not authorized"
		awsOrigErr := errors.New(awsErrCode)
		awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
		// mock failed call to: GlobalSignOut(input *cognitoidentityprovider.GlobalSignOutInput) (*cognitoidentityprovider.GlobalSignOutOutput, error)
		m.GlobalSignOutFunc = func(_ *cognitoidentityprovider.GlobalSignOutInput) (*cognitoidentityprovider.GlobalSignOutOutput, error) {
			return nil, awsErr
		}

		request := httptest.NewRequest(http.MethodDelete, signOutEndPoint, http.NoBody)
		request.Header.Set(AccessTokenHeaderName, "Bearer zzzz-yyyy-xxxx")

		successResponse, errorResponse := api.SignOutHandler(ctx, w, request)

		So(successResponse, ShouldBeNil)
		So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
		So(len(errorResponse.Errors), ShouldEqual, 1)
		So(errorResponse.Errors[0].Error(), ShouldEqual, awsErr.Error())
	})
}

func TestAPI_RefreshHandler(t *testing.T) {
	var (
		ctx                              = context.Background()
		accessToken, returnIDToken       = "aaaa.bbbb.cccc", "llll.mmmm.nnnn"
		expireLength               int64 = 500
	)

	api, w, m := apiSetup()

	// mock call to: InitiateAuth(input *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error)
	m.InitiateAuthFunc = func(_ *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error) {
		return &cognitoidentityprovider.InitiateAuthOutput{
			AuthenticationResult: &cognitoidentityprovider.AuthenticationResultType{
				AccessToken: &accessToken,
				ExpiresIn:   &expireLength,
				IdToken:     &returnIDToken,
			},
		}, nil
	}

	Convey("Token refresh success: no errors added to ErrorResponse Errors list", t, func() {
		request := httptest.NewRequest(http.MethodPut, tokenRefreshEndPoint, http.NoBody)
		idToken := mock.GenerateMockIDToken("test@ons.gov.uk")
		So(idToken, ShouldNotEqual, "")
		request.Header.Set(IDTokenHeaderName, idToken)
		request.Header.Set(RefreshTokenHeaderName, "aaaa.bbbb.cccc.dddd.eeee")

		successResponse, errorResponse := api.RefreshHandler(ctx, w, request)

		So(errorResponse, ShouldBeNil)
		So(successResponse.Status, ShouldEqual, http.StatusCreated)
		var responseBody map[string]interface{}
		err := json.Unmarshal(successResponse.Body, &responseBody)
		So(err, ShouldBeNil)
		So(responseBody["expirationTime"], ShouldNotBeNil)
	})

	Convey("Token refresh validation error: adds an error to the ErrorResponse and sets its Status to 400", t, func() {
		request := httptest.NewRequest(http.MethodPut, tokenRefreshEndPoint, http.NoBody)
		request.Header.Set(IDTokenHeaderName, "")
		request.Header.Set(RefreshTokenHeaderName, "aaaa.bbbb.cccc.dddd.eeee")

		successResponse, errorResponse := api.RefreshHandler(ctx, w, request)

		So(successResponse, ShouldBeNil)
		So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
		So(len(errorResponse.Errors), ShouldEqual, 1)
		So(errorResponse.Errors[0].Error(), ShouldEqual, models.InvalidTokenError)
	})

	Convey("Token refresh Cognito internal error: adds an error to the ErrorResponse and sets its Status to 500", t, func() {
		awsOrigErr := errors.New(awsErrCode)
		awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
		// mock failed call to: InitiateAuth(input *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error)
		m.InitiateAuthFunc = func(_ *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error) {
			return nil, awsErr
		}

		request := httptest.NewRequest(http.MethodDelete, signOutEndPoint, http.NoBody)
		idToken := mock.GenerateMockIDToken("test@ons.gov.uk")
		So(idToken, ShouldNotEqual, "")
		request.Header.Set(IDTokenHeaderName, idToken)
		request.Header.Set(RefreshTokenHeaderName, "aaaa.bbbb.cccc.dddd.eeee")

		successResponse, errorResponse := api.RefreshHandler(ctx, w, request)

		So(successResponse, ShouldBeNil)
		So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
		So(len(errorResponse.Errors), ShouldEqual, 1)
		So(errorResponse.Errors[0].Error(), ShouldEqual, awsErr.Error())
	})

	Convey("Token refresh Cognito request error: adds an error to the ErrorResponse and sets its Status to 403", t, func() {
		awsErrCode := "NotAuthorizedException"
		awsErrMessage := "User is not authorized"
		awsOrigErr := errors.New(awsErrCode)
		awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
		// mock failed call to: InitiateAuth(input *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error)
		m.InitiateAuthFunc = func(_ *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error) {
			return nil, awsErr
		}

		request := httptest.NewRequest(http.MethodDelete, signOutEndPoint, http.NoBody)
		idToken := mock.GenerateMockIDToken("test@ons.gov.uk")
		So(idToken, ShouldNotEqual, "")
		request.Header.Set(IDTokenHeaderName, idToken)
		request.Header.Set(RefreshTokenHeaderName, "aaaa.bbbb.cccc.dddd.eeee")

		successResponse, errorResponse := api.RefreshHandler(ctx, w, request)

		So(successResponse, ShouldBeNil)
		So(errorResponse.Status, ShouldEqual, http.StatusForbidden)
		So(len(errorResponse.Errors), ShouldEqual, 1)
		So(errorResponse.Errors[0].Error(), ShouldEqual, awsErr.Error())
	})
}

func TestSignOutAllUsersHandlerAccessForProcessing(t *testing.T) {
	var ctx = context.Background()

	api, w, m := apiSetup()

	Convey("Testing users global signout - handler", t, func() {
		signOutAllUsersTests := []struct {
			listUsersFunction              func(userInput *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error)
			adminUserGlobalSignOutFunction func(signOutInput *cognitoidentityprovider.AdminUserGlobalSignOutInput) (*cognitoidentityprovider.AdminUserGlobalSignOutOutput, error)
			httpResponse                   int
		}{
			{
				// 200 response from Cognito - 202 from identity api
				func(_ *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error) {
					users := mock.BulkGenerateUsers(3, nil)
					users.PaginationToken = nil
					return users, nil
				},
				func(_ *cognitoidentityprovider.AdminUserGlobalSignOutInput) (*cognitoidentityprovider.AdminUserGlobalSignOutOutput, error) {
					return &cognitoidentityprovider.AdminUserGlobalSignOutOutput{}, nil
				},
				http.StatusAccepted,
			},
		}
		for _, tt := range signOutAllUsersTests {
			m.ListUsersFunc = tt.listUsersFunction
			m.AdminUserGlobalSignOutFunc = tt.adminUserGlobalSignOutFunction
			r := httptest.NewRequest(http.MethodPost, usersEndPoint, http.NoBody)

			successResponse, errorResponse := api.SignOutAllUsersHandler(ctx, w, r)
			So(successResponse.Status, ShouldEqual, tt.httpResponse)
			So(errorResponse, ShouldBeNil)
		}
	})
}

func TestSignOutAllUsersHandlerInternalServerError(t *testing.T) {
	var ctx = context.Background()

	api, w, m := apiSetup()

	Convey("Testing users global signout - handler", t, func() {
		signOutAllUsersTests := []struct {
			listUsersFunction              func(userInput *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error)
			adminUserGlobalSignOutFunction func(signOutInput *cognitoidentityprovider.AdminUserGlobalSignOutInput) (*cognitoidentityprovider.AdminUserGlobalSignOutOutput, error)
			httpResponse                   int
		}{
			{
				// 500 response from Cognito's ListUsers API endpoint
				func(_ *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error) {
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				func(_ *cognitoidentityprovider.AdminUserGlobalSignOutInput) (*cognitoidentityprovider.AdminUserGlobalSignOutOutput, error) {
					return &cognitoidentityprovider.AdminUserGlobalSignOutOutput{}, nil
				},
				http.StatusInternalServerError,
			},
		}
		for _, tt := range signOutAllUsersTests {
			m.ListUsersFunc = tt.listUsersFunction
			m.AdminUserGlobalSignOutFunc = tt.adminUserGlobalSignOutFunction
			r := httptest.NewRequest(http.MethodGet, usersEndPoint, http.NoBody)

			successResponse, errorResponse := api.SignOutAllUsersHandler(ctx, w, r)
			So(successResponse, ShouldBeNil)
			So(errorResponse.Status, ShouldEqual, tt.httpResponse)
		}
	})
}

func TestSignOutAllUsersGoRoutine(t *testing.T) {
	var ctx = context.Background()

	api, _, m := apiSetup()

	// a list of known UUIDs for testing
	userNamesList := []string{
		"41af9e4e-3bb8-46a2-ba33-19acc6698d5f",
		"a03dfc5e-39b7-4229-a87c-a2ee91bc6870",
		"4affc660-3c4b-4111-85bb-c83e76f7f81d",
		"0a7a64b7-e61b-4a37-b5fc-33df36f7dfd7",
	}

	Convey("Testing users global signout - go routine", t, func() {
		signOutAllUsersTests := []struct {
			listUsersFunction              func(userInput *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error)
			adminUserGlobalSignOutFunction func(signOutInput *cognitoidentityprovider.AdminUserGlobalSignOutInput) (*cognitoidentityprovider.AdminUserGlobalSignOutOutput, error)
			globalUserSignOutMod           models.GlobalSignOut
			numberOfUsers                  int
			expectedResults                int
			httpResponse                   int
		}{
			{
				// 200 response from Cognito - 202 from identity api
				func(_ *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error) {
					users := mock.BulkGenerateUsers(3, nil)
					users.PaginationToken = nil
					return users, nil
				},
				func(_ *cognitoidentityprovider.AdminUserGlobalSignOutInput) (*cognitoidentityprovider.AdminUserGlobalSignOutOutput, error) {
					return &cognitoidentityprovider.AdminUserGlobalSignOutOutput{}, nil
				},
				models.GlobalSignOut{
					ResultsChannel: make(chan string, 4),
					BackoffSchedule: []time.Duration{
						1 * time.Second,
						2 * time.Second,
						3 * time.Second,
					},
					RetryAllowed: true,
				},
				4,
				4,
				http.StatusAccepted,
			},
			{
				// 500 response from Cognito
				func(_ *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error) {
					users := mock.BulkGenerateUsers(3, userNamesList)
					users.PaginationToken = nil
					return users, nil
				},
				func(signOutInput *cognitoidentityprovider.AdminUserGlobalSignOutInput) (*cognitoidentityprovider.AdminUserGlobalSignOutOutput, error) {
					if *signOutInput.Username == userNamesList[3] {
						awsOrigErr := errors.New(awsErrCode)
						awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
						return nil, awsErr
					}
					return &cognitoidentityprovider.AdminUserGlobalSignOutOutput{}, nil
				},
				models.GlobalSignOut{
					ResultsChannel: make(chan string, 4),
					BackoffSchedule: []time.Duration{
						1 * time.Second,
						2 * time.Second,
						3 * time.Second,
					},
					RetryAllowed: true,
				},
				4,
				3,
				http.StatusAccepted,
			},
			{
				// 429 response from Cognito
				func(_ *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error) {
					users := mock.BulkGenerateUsers(10, userNamesList)
					users.PaginationToken = nil
					return users, nil
				},
				func(signOutInput *cognitoidentityprovider.AdminUserGlobalSignOutInput) (*cognitoidentityprovider.AdminUserGlobalSignOutOutput, error) {
					if *signOutInput.Username == userNamesList[3] {
						awsErrCode := "TooManyRequestsException"
						awsErrMessage := "Too many requets received"
						awsOrigErr := errors.New(awsErrCode)
						awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
						return nil, awsErr
					}
					return &cognitoidentityprovider.AdminUserGlobalSignOutOutput{}, nil
				},
				models.GlobalSignOut{
					ResultsChannel: make(chan string, 10),
					BackoffSchedule: []time.Duration{
						1 * time.Second,
						2 * time.Second,
						3 * time.Second,
					},
					RetryAllowed: true,
				},
				10,
				9,
				http.StatusAccepted,
			},
		}
		for _, tt := range signOutAllUsersTests {
			m.ListUsersFunc = tt.listUsersFunction
			m.AdminUserGlobalSignOutFunc = tt.adminUserGlobalSignOutFunction

			// test concurrent go routine
			usersList := models.UsersList{}
			generatedUsers := mock.BulkGenerateUsers(tt.numberOfUsers, userNamesList)
			usersList.MapCognitoUsers(&generatedUsers.Users)

			api.SignOutUsersWorker(ctx, &tt.globalUserSignOutMod, &usersList.Users)

			// we should receive the expected number of processed usernames on the ResultsChannel
			So(len(tt.globalUserSignOutMod.ResultsChannel), ShouldEqual, tt.expectedResults)
		}
	})
}

func TestSignOutAllUsersGetAllUsersList(t *testing.T) {
	var ctx = context.Background()

	api, _, m := apiSetup()

	Convey("Testing users global signout - go routine", t, func() {
		var (
			paginationToken, usersFilterString = "abc-123-xyz-345-xxx", "status=\"Enabled\""
			backoff                            = []time.Duration{
				1 * time.Second,
				2 * time.Second,
				3 * time.Second,
			}
			errCode int
		)
		getAllUsersTests := []struct {
			listUsersFunction func(userInput *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error)
			BackoffSchedule   []time.Duration
			expectedUserNumb  int
			httpResponse      []int
		}{
			{
				func(userInput *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error) {
					if userInput.PaginationToken != nil {
						users := mock.BulkGenerateUsers(3, nil)
						users.PaginationToken = nil
						return users, nil
					}
					users := mock.BulkGenerateUsers(14, nil)
					users.PaginationToken = &paginationToken
					return users, nil
				},
				backoff,
				17,
				[]int{
					http.StatusOK,
				},
			},
			{
				func(_ *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error) {
					awsOrigErr := errors.New(awsErrCode)
					awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					return nil, awsErr
				},
				backoff,
				0,
				[]int{
					http.StatusInternalServerError,
				},
			},
			{
				func(userInput *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error) {
					if userInput.PaginationToken != nil {
						awsOrigErr := errors.New(awsErrCode)
						awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
						// set error code index reference to 1 - expecting a http.StatusInternalServerError here
						errCode = 1
						return nil, awsErr
					}
					users := mock.BulkGenerateUsers(3, nil)
					return users, nil
				},
				backoff,
				3,
				[]int{
					http.StatusOK,
					http.StatusInternalServerError,
				},
			},
		}
		for _, tt := range getAllUsersTests {
			m.ListUsersFunc = tt.listUsersFunction
			usersList, awsErr := api.ListUsersWorker(ctx, &usersFilterString, tt.BackoffSchedule)

			// we should receive the expected number of processed usernames on the ResultsChannel
			if tt.httpResponse[errCode] >= http.StatusBadRequest {
				So(usersList, ShouldBeNil)
				So(awsErr.Status, ShouldEqual, tt.httpResponse[errCode])
			} else {
				So(usersList, ShouldNotBeNil)
				So(awsErr, ShouldBeNil)
				So(len(*usersList), ShouldEqual, tt.expectedUserNumb)
			}
		}
	})
}
