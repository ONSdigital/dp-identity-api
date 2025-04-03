package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/smithy-go"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-identity-api/v2/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	usersEndPoint                                 = "http://localhost:25600/v1/users"
	usersEndPointWithActiveFilterTrue             = "http://localhost:25600/v1/users?active=true"
	usersEndPointWithActiveFilterFalse            = "http://localhost:25600/v1/users?active=false"
	usersEndPointWithActiveFilterError            = "http://localhost:25600/v1/user?active=false"
	usersEndPointWithSortByEmail                  = "http://localhost:25600/v1/users?sort=email"
	usersEndPointWithSortByEmailAsc               = "http://localhost:25600/v1/users?sort=email:asc"
	usersEndPointWithSortByEmailDesc              = "http://localhost:25600/v1/users?sort=email:desc"
	usersEndPointWithSortBy2FieldsDesc            = "http://localhost:25600/v1/users?sort=forename:desc,lastname:desc"
	usersEndPointWithSortBy2KnownFieldsAndUnknown = "http://localhost:25600/v1/users?sort=forename:desc,lastname:desc,dog"
	userEndPoint                                  = "http://localhost:25600/v1/users/abcd1234"
	changePasswordEndPoint                        = "http://localhost:25600/v1/users/self/password" // #nosec
	requestResetEndPoint                          = "http://localhost:25600/v1/password-reset"
	userListGroupsEndPoint                        = "http://localhost:25600/v1/users/abcd1234/groups"
)

func TestCreateUserHandler(t *testing.T) {
	var (
		ctx                                = context.Background()
		name, surname, email, invalidEmail = "bob", "bobbings", "foo_bar123@ext.ons.gov.uk", "foo_bar123@test.ons.gov.ie"
		userException                      = "UsernameExistsException: User account already exists"
		userStatusType                     = types.UserStatusTypeUnconfirmed
	)

	api, w, m := apiSetup()

	Convey("Admin create user - check expected responses", t, func() {
		adminCreateUsersTests := []struct {
			listUsersFunction   func(_ context.Context, userInput *cognitoidentityprovider.ListUsersInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUsersOutput, error)
			createUsersFunction func(_ context.Context, userInput *cognitoidentityprovider.AdminCreateUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminCreateUserOutput, error)
			httpResponse        int
		}{
			{
				// 200 response - no duplicate emails found
				func(_ context.Context, _ *cognitoidentityprovider.ListUsersInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUsersOutput, error) {
					users := &models.ListUsersOutput{
						ListUsersOutput: &cognitoidentityprovider.ListUsersOutput{
							Users: []types.UserType{},
						},
					}
					return users.ListUsersOutput, nil
				},
				// 201 response - user created
				func(_ context.Context, _ *cognitoidentityprovider.AdminCreateUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminCreateUserOutput, error) {
					user := models.CreateUserOutput{
						UserOutput: &cognitoidentityprovider.AdminCreateUserOutput{
							User: &types.UserType{
								Username:   &name,
								UserStatus: userStatusType,
							},
						},
					}
					return user.UserOutput, nil
				},
				http.StatusCreated,
			},
			{
				// 200 response - no duplicate emails found
				func(_ context.Context, _ *cognitoidentityprovider.ListUsersInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUsersOutput, error) {
					users := &models.ListUsersOutput{
						ListUsersOutput: &cognitoidentityprovider.ListUsersOutput{
							Users: []types.UserType{},
						},
					}
					return users.ListUsersOutput, nil
				},
				// 400 response - user already exists
				func(_ context.Context, _ *cognitoidentityprovider.AdminCreateUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminCreateUserOutput, error) {
					var userExistsException types.UsernameExistsException
					userExistsException.Message = &userException
					// userExistsException.RespMetadata.StatusCode = http.StatusBadRequest	// TODO find out how to replace this for aws-sdk-go-v2

					return &cognitoidentityprovider.AdminCreateUserOutput{}, &userExistsException
				},
				http.StatusBadRequest,
			},
			{
				// 400 response - duplicate email found
				func(_ context.Context, _ *cognitoidentityprovider.ListUsersInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUsersOutput, error) {
					users := &models.ListUsersOutput{
						ListUsersOutput: &cognitoidentityprovider.ListUsersOutput{
							Users: []types.UserType{
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
				func(_ context.Context, _ *cognitoidentityprovider.ListUsersInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUsersOutput, error) {
					users := &models.ListUsersOutput{
						ListUsersOutput: &cognitoidentityprovider.ListUsersOutput{
							Users: []types.UserType{},
						},
					}
					return users.ListUsersOutput, nil
				},
				// 500 response - internal error exception
				func(_ context.Context, _ *cognitoidentityprovider.AdminCreateUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminCreateUserOutput, error) {
					var internalErrorException types.InternalErrorException
					internalErrorException.Message = &userException
					//internalErrorException.RespMetadata.StatusCode = http.StatusInternalServerError	// TODO find out how to replace this for aws-sdk-go-v2

					return &cognitoidentityprovider.AdminCreateUserOutput{}, &internalErrorException
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
			listUsersFunction func(_ context.Context, userInput *cognitoidentityprovider.ListUsersInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUsersOutput, error)
			httpResponse      int
		}{
			{
				// 200 response from Cognito
				func(_ context.Context, _ *cognitoidentityprovider.ListUsersInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUsersOutput, error) {
					users := &cognitoidentityprovider.ListUsersOutput{
						Users: []types.UserType{},
					}
					return users, nil
				},
				http.StatusOK,
			},
			{
				// 500 response from Cognito
				func(_ context.Context, _ *cognitoidentityprovider.ListUsersInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUsersOutput, error) {
					//awsOrigErr := errors.New(awsErrCode)
					//awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					awsOrigErr := smithy.ErrorFault(1) //server error
					awsErr := &smithy.GenericAPIError{
						Code:    awsErrCode,
						Message: awsErrMessage,
						Fault:   awsOrigErr,
					}
					return nil, awsErr
				},
				http.StatusInternalServerError,
			},
		}

		for _, tt := range adminCreateUsersTests {
			m.ListUsersFunc = tt.listUsersFunction

			r := httptest.NewRequest(http.MethodGet, usersEndPoint, http.NoBody)
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

func TestListUserHandlerWithFilter(t *testing.T) {
	var ctx = context.Background()
	api, w, m := apiSetup()

	Convey("List user - check expected responses", t, func() {
		listUsersTest := []struct {
			description       string
			endpoint          *http.Request
			listUsersFunction func(_ context.Context, userInput *cognitoidentityprovider.ListUsersInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUsersOutput, error)
			assertions        func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse)
		}{
			{
				"200 response from Cognito active filter true",
				httptest.NewRequest(http.MethodGet, usersEndPointWithActiveFilterTrue, http.NoBody),
				func(_ context.Context, _ *cognitoidentityprovider.ListUsersInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUsersOutput, error) {
					users := &cognitoidentityprovider.ListUsersOutput{
						Users: []types.UserType{},
					}
					return users, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(errorResponse, ShouldBeNil)
					So(successResponse, ShouldNotBeNil)
				},
			},
			{
				"200 response from Cognito active filter false",
				httptest.NewRequest(http.MethodGet, usersEndPointWithActiveFilterFalse, http.NoBody),
				func(_ context.Context, _ *cognitoidentityprovider.ListUsersInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUsersOutput, error) {
					users := &cognitoidentityprovider.ListUsersOutput{
						Users: []types.UserType{},
					}
					return users, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(errorResponse, ShouldBeNil)
					So(successResponse, ShouldNotBeNil)
					So(successResponse.Status, ShouldEqual, 200)
				},
			},
			{
				"200 response from Cognito no active filter",
				httptest.NewRequest(http.MethodGet, usersEndPoint, http.NoBody),
				func(_ context.Context, _ *cognitoidentityprovider.ListUsersInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUsersOutput, error) {
					users := &cognitoidentityprovider.ListUsersOutput{
						Users: []types.UserType{},
					}
					return users, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(errorResponse, ShouldBeNil)
					So(successResponse, ShouldNotBeNil)
					So(successResponse.Status, ShouldEqual, 200)
				},
			},

			{
				"400 response from Cognito active filter incorrect",
				httptest.NewRequest(http.MethodGet, usersEndPointWithActiveFilterError, http.NoBody),
				func(_ context.Context, _ *cognitoidentityprovider.ListUsersInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUsersOutput, error) {
					users := &cognitoidentityprovider.ListUsersOutput{
						Users: []types.UserType{},
					}
					return users, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(errorResponse, ShouldNotBeNil)
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldEqual, 400)
					So(errorResponse.Errors[0].Error(), ShouldResemble, "InvalidFilterQuery")
				},
			},
		}

		for _, tt := range listUsersTest {
			Convey(tt.description, func() {
				m.ListUsersFunc = tt.listUsersFunction
				r := tt.endpoint
				successResponse, errorResponse := api.ListUsersHandler(ctx, w, r)
				tt.assertions(successResponse, errorResponse)
			},
			)
		}
	})
}

func TestListUserHandlerWithSort(t *testing.T) {
	var (
		ctx = context.Background()
	)

	api, w, m := apiSetup()

	Convey("List user - check expected responses", t, func() {
		listUsersTest := []struct {
			description       string
			endpoint          *http.Request
			listUsersFunction func(_ context.Context, userInput *cognitoidentityprovider.ListUsersInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUsersOutput, error)
			assertions        func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse)
		}{
			{
				description: "200 response from Cognito sort asc by Email ",
				endpoint:    httptest.NewRequest(http.MethodGet, usersEndPointWithSortByEmail, http.NoBody),
				listUsersFunction: func(_ context.Context, _ *cognitoidentityprovider.ListUsersInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUsersOutput, error) {
					var cognitoUsersList []types.UserType
					cognitoUsersList = listUserOutput("Adam", "Zee", "email1@ons.gov.uk", "user-1", cognitoUsersList)
					cognitoUsersList = listUserOutput("Bob", "Yellow", "email2@ons.gov.uk", "user-2", cognitoUsersList)
					cognitoUsersList = listUserOutput("Colin", "White", "email3@ons.gov.uk", "user-3", cognitoUsersList)

					users := &cognitoidentityprovider.ListUsersOutput{
						Users: cognitoUsersList,
					}
					return users, nil
				},
				assertions: func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(errorResponse, ShouldBeNil)
					So(successResponse, ShouldNotBeNil)
					So(successResponse.Status, ShouldEqual, 200)
				},
			},
			{
				description: "200 response from Cognito sort EmailAsc  ",
				endpoint:    httptest.NewRequest(http.MethodGet, usersEndPointWithSortByEmailAsc, http.NoBody),
				listUsersFunction: func(_ context.Context, _ *cognitoidentityprovider.ListUsersInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUsersOutput, error) {
					var cognitoUsersList []types.UserType
					cognitoUsersList = listUserOutput("Adam", "Zee", "email1@ons.gov.uk", "user-1", cognitoUsersList)
					cognitoUsersList = listUserOutput("Bob", "Yellow", "email2@ons.gov.uk", "user-2", cognitoUsersList)
					cognitoUsersList = listUserOutput("Colin", "White", "email3@ons.gov.uk", "user-3", cognitoUsersList)

					users := &cognitoidentityprovider.ListUsersOutput{
						Users: cognitoUsersList,
					}
					return users, nil
				},
				assertions: func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(errorResponse, ShouldBeNil)
					So(successResponse, ShouldNotBeNil)
					So(successResponse.Status, ShouldEqual, 200)
				},
			},
			{
				description: "200 response from Cognito sort EmailDesc  ",
				endpoint:    httptest.NewRequest(http.MethodGet, usersEndPointWithSortByEmailDesc, http.NoBody),
				listUsersFunction: func(_ context.Context, _ *cognitoidentityprovider.ListUsersInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUsersOutput, error) {
					var cognitoUsersList []types.UserType
					cognitoUsersList = listUserOutput("Adam", "Zee", "email1@ons.gov.uk", "user-1", cognitoUsersList)
					cognitoUsersList = listUserOutput("Bob", "Yellow", "email2@ons.gov.uk", "user-2", cognitoUsersList)
					cognitoUsersList = listUserOutput("Colin", "White", "email3@ons.gov.uk", "user-3", cognitoUsersList)

					users := &cognitoidentityprovider.ListUsersOutput{
						Users: cognitoUsersList,
					}
					return users, nil
				},
				assertions: func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(errorResponse, ShouldBeNil)
					So(successResponse, ShouldNotBeNil)
					So(successResponse.Status, ShouldEqual, 200)
				},
			},
			{
				description: "200 response from Cognito sort forename:desc, lastname:desc  ",
				endpoint:    httptest.NewRequest(http.MethodGet, usersEndPointWithSortBy2FieldsDesc, http.NoBody),
				listUsersFunction: func(_ context.Context, _ *cognitoidentityprovider.ListUsersInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUsersOutput, error) {
					var cognitoUsersList []types.UserType
					cognitoUsersList = listUserOutput("Adam", "Zee", "email1@ons.gov.uk", "user-1", cognitoUsersList)
					cognitoUsersList = listUserOutput("Bob", "Yellow", "email2@ons.gov.uk", "user-2", cognitoUsersList)
					cognitoUsersList = listUserOutput("Colin", "White", "email3@ons.gov.uk", "user-3", cognitoUsersList)

					users := &cognitoidentityprovider.ListUsersOutput{
						Users: cognitoUsersList,
					}
					return users, nil
				},
				assertions: func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(errorResponse, ShouldBeNil)
					So(successResponse, ShouldNotBeNil)
					So(successResponse.Status, ShouldEqual, 200)
				},
			},
			{
				description: "200 response from Cognito sort forename:desc, lastname:desc, dog  ",
				endpoint:    httptest.NewRequest(http.MethodGet, usersEndPointWithSortBy2KnownFieldsAndUnknown, http.NoBody),
				listUsersFunction: func(_ context.Context, _ *cognitoidentityprovider.ListUsersInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUsersOutput, error) {
					var cognitoUsersList []types.UserType
					cognitoUsersList = listUserOutput("Adam", "Zee", "email1@ons.gov.uk", "user-1", cognitoUsersList)
					cognitoUsersList = listUserOutput("Bob", "Yellow", "email2@ons.gov.uk", "user-2", cognitoUsersList)
					cognitoUsersList = listUserOutput("Colin", "White", "email3@ons.gov.uk", "user-3", cognitoUsersList)

					users := &cognitoidentityprovider.ListUsersOutput{
						Users: cognitoUsersList,
					}
					return users, nil
				},
				assertions: func(_ *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(errorResponse, ShouldNotBeNil)
					So(errorResponse.Errors[0].Error(), ShouldResemble, " request query sort parameter not found dog")
					fmt.Println(errorResponse.Errors[0])

					So(errorResponse.Status, ShouldEqual, 400)
				},
			},
		}
		for _, tt := range listUsersTest {
			Convey(tt.description, func() {
				m.ListUsersFunc = tt.listUsersFunction
				r := tt.endpoint
				successResponse, errorResponse := api.ListUsersHandler(ctx, w, r)
				tt.assertions(successResponse, errorResponse)
			},
			)
		}
	})
}

func listUserOutput(forename, surname, email, id string, cognitoUsersList []types.UserType) []types.UserType {
	var status = types.UserStatusTypeConfirmed
	cognitoUser := types.UserType{
		Attributes: []types.AttributeType{
			{
				Name:  aws.String("given_name"),
				Value: &forename,
			},
			{
				Name:  aws.String("family_name"),
				Value: &surname,
			},
			{
				Name:  aws.String("email"),
				Value: &email,
			},
		},
		UserStatus: status,
		Username:   &id,
		Enabled:    true,
	}

	cognitoUsersList = append(cognitoUsersList, cognitoUser)
	return cognitoUsersList
}

func TestGetUserHandler(t *testing.T) {
	var (
		ctx                                      = context.Background()
		forename, lastname, email, userID        = "bob", "bobbings", "foo_bar123@ext.ons.gov.uk", "abcd1234"
		givenNameAttr, familyNameAttr, emailAttr = "given_name", "family_name", "email"
		status                                   = types.UserStatusTypeUnconfirmed
	)

	api, w, m := apiSetup()

	Convey("Get user - check expected responses", t, func() {
		adminGetUsersTests := []struct {
			getUserFunction func(_ context.Context, userInput *cognitoidentityprovider.AdminGetUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminGetUserOutput, error)
			httpResponse    int
		}{
			{
				// 200 response from Cognito
				func(_ context.Context, _ *cognitoidentityprovider.AdminGetUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminGetUserOutput, error) {
					user := &cognitoidentityprovider.AdminGetUserOutput{
						UserAttributes: []types.AttributeType{
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
						UserStatus: status,
						Username:   &userID,
						Enabled:    true,
					}
					return user, nil
				},
				http.StatusOK,
			},
			{
				// 500 response from Cognito
				func(_ context.Context, _ *cognitoidentityprovider.AdminGetUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminGetUserOutput, error) {
					//awsOrigErr := errors.New(awsErrCode)
					//awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					awsOrigErr := smithy.ErrorFault(1) //server error
					awsErr := &smithy.GenericAPIError{
						Code:    awsErrCode,
						Message: awsErrMessage,
						Fault:   awsOrigErr,
					}
					return nil, awsErr
				},
				http.StatusInternalServerError,
			},
			{
				// 404 response from Cognito
				func(_ context.Context, _ *cognitoidentityprovider.AdminGetUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminGetUserOutput, error) {
					//awsOrigErr := errors.New(awsUNFErrCode)
					//awsErr := awserr.New(awsUNFErrCode, awsUNFErrMessage, awsOrigErr)
					awsOrigErr := smithy.ErrorFault(2) //client error
					awsErr := &smithy.GenericAPIError{
						Code:    awsUNFErrCode,
						Message: awsUNFErrMessage,
						Fault:   awsOrigErr,
					}
					return nil, awsErr
				},
				http.StatusNotFound,
			},
		}

		for _, tt := range adminGetUsersTests {
			m.AdminGetUserFunc = tt.getUserFunction

			r := httptest.NewRequest(http.MethodGet, userEndPoint, http.NoBody)

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
		ctx                                      = context.Background()
		forename, lastname, email, userID        = "bob", "bobbings", "foo_bar123@ext.ons.gov.uk", "abcd1234"
		givenNameAttr, familyNameAttr, emailAttr = "given_name", "family_name", "email"
		status                                   = types.UserStatusTypeConfirmed
	)

	api, w, m := apiSetup()

	successfullyGetUser := []types.AttributeType{
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
			updateUserFunction  func(_ context.Context, userInput *cognitoidentityprovider.AdminUpdateUserAttributesInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminUpdateUserAttributesOutput, error)
			getUserFunction     func(_ context.Context, userInput *cognitoidentityprovider.AdminGetUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminGetUserOutput, error)
			enableUserFunction  func(_ context.Context, userInput *cognitoidentityprovider.AdminEnableUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminEnableUserOutput, error)
			disableUserFunction func(_ context.Context, userInput *cognitoidentityprovider.AdminDisableUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminDisableUserOutput, error)
			userForename        string
			userActive          bool
			assertions          func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse)
		}{
			// 200 response from Cognito
			{
				func(_ context.Context, _ *cognitoidentityprovider.AdminUpdateUserAttributesInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminUpdateUserAttributesOutput, error) {
					user := &cognitoidentityprovider.AdminUpdateUserAttributesOutput{}
					return user, nil
				},
				func(_ context.Context, _ *cognitoidentityprovider.AdminGetUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminGetUserOutput, error) {
					user := &cognitoidentityprovider.AdminGetUserOutput{
						UserAttributes: successfullyGetUser,
						UserStatus:     status,
						Username:       &userID,
						Enabled:        true,
					}
					return user, nil
				},
				func(_ context.Context, _ *cognitoidentityprovider.AdminEnableUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminEnableUserOutput, error) {
					return &cognitoidentityprovider.AdminEnableUserOutput{}, nil
				},
				func(_ context.Context, _ *cognitoidentityprovider.AdminDisableUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminDisableUserOutput, error) {
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
			// local validation failure
			{
				func(_ context.Context, _ *cognitoidentityprovider.AdminUpdateUserAttributesInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminUpdateUserAttributesOutput, error) {
					user := &cognitoidentityprovider.AdminUpdateUserAttributesOutput{}
					return user, nil
				},
				func(_ context.Context, _ *cognitoidentityprovider.AdminGetUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminGetUserOutput, error) {
					user := &cognitoidentityprovider.AdminGetUserOutput{
						UserAttributes: successfullyGetUser,
						UserStatus:     status,
						Username:       &userID,
						Enabled:        true,
					}
					return user, nil
				},
				func(_ context.Context, _ *cognitoidentityprovider.AdminEnableUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminEnableUserOutput, error) {
					return &cognitoidentityprovider.AdminEnableUserOutput{}, nil
				},
				func(_ context.Context, _ *cognitoidentityprovider.AdminDisableUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminDisableUserOutput, error) {
					return &cognitoidentityprovider.AdminDisableUserOutput{}, nil
				},
				"",
				true,
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
				},
			},
			// 404 response from Cognito enable user
			{
				func(_ context.Context, _ *cognitoidentityprovider.AdminUpdateUserAttributesInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminUpdateUserAttributesOutput, error) {
					user := &cognitoidentityprovider.AdminUpdateUserAttributesOutput{}
					return user, nil
				},
				func(_ context.Context, _ *cognitoidentityprovider.AdminGetUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminGetUserOutput, error) {
					user := &cognitoidentityprovider.AdminGetUserOutput{
						UserAttributes: successfullyGetUser,
						UserStatus:     status,
						Username:       &userID,
						Enabled:        true,
					}
					return user, nil
				},
				func(_ context.Context, _ *cognitoidentityprovider.AdminEnableUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminEnableUserOutput, error) {
					//awsOrigErr := errors.New(awsUNFErrCode)
					//awsErr := awserr.New(awsUNFErrCode, awsUNFErrMessage, awsOrigErr)
					awsOrigErr := smithy.ErrorFault(2) //client error
					awsErr := &smithy.GenericAPIError{
						Code:    awsUNFErrCode,
						Message: awsUNFErrMessage,
						Fault:   awsOrigErr,
					}
					return nil, awsErr
				},
				func(_ context.Context, _ *cognitoidentityprovider.AdminDisableUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminDisableUserOutput, error) {
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
				func(_ context.Context, _ *cognitoidentityprovider.AdminUpdateUserAttributesInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminUpdateUserAttributesOutput, error) {
					user := &cognitoidentityprovider.AdminUpdateUserAttributesOutput{}
					return user, nil
				},
				func(_ context.Context, _ *cognitoidentityprovider.AdminGetUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminGetUserOutput, error) {
					user := &cognitoidentityprovider.AdminGetUserOutput{
						UserAttributes: successfullyGetUser,
						UserStatus:     status,
						Username:       &userID,
						Enabled:        true,
					}
					return user, nil
				},
				func(_ context.Context, _ *cognitoidentityprovider.AdminEnableUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminEnableUserOutput, error) {
					//awsOrigErr := errors.New(awsErrCode)
					//awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					awsOrigErr := smithy.ErrorFault(1) // server error
					awsErr := &smithy.GenericAPIError{
						Code:    awsErrCode,
						Message: awsErrMessage,
						Fault:   awsOrigErr,
					}
					return nil, awsErr
				},
				func(_ context.Context, _ *cognitoidentityprovider.AdminDisableUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminDisableUserOutput, error) {
					return &cognitoidentityprovider.AdminDisableUserOutput{}, nil
				},
				forename,
				true,
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
				},
			},
			// 404 response from Cognito disable user
			{
				func(_ context.Context, _ *cognitoidentityprovider.AdminUpdateUserAttributesInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminUpdateUserAttributesOutput, error) {
					user := &cognitoidentityprovider.AdminUpdateUserAttributesOutput{}
					return user, nil
				},
				func(_ context.Context, _ *cognitoidentityprovider.AdminGetUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminGetUserOutput, error) {
					user := &cognitoidentityprovider.AdminGetUserOutput{
						UserAttributes: successfullyGetUser,
						UserStatus:     status,
						Username:       &userID,
						Enabled:        true,
					}
					return user, nil
				},
				func(_ context.Context, _ *cognitoidentityprovider.AdminEnableUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminEnableUserOutput, error) {
					return &cognitoidentityprovider.AdminEnableUserOutput{}, nil
				},
				func(_ context.Context, _ *cognitoidentityprovider.AdminDisableUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminDisableUserOutput, error) {
					//awsOrigErr := errors.New(awsUNFErrCode)
					//awsErr := awserr.New(awsUNFErrCode, awsUNFErrMessage, awsOrigErr)
					awsOrigErr := smithy.ErrorFault(2) // client error
					awsErr := &smithy.GenericAPIError{
						Code:    awsUNFErrCode,
						Message: awsUNFErrMessage,
						Fault:   awsOrigErr,
					}
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
				func(_ context.Context, _ *cognitoidentityprovider.AdminUpdateUserAttributesInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminUpdateUserAttributesOutput, error) {
					user := &cognitoidentityprovider.AdminUpdateUserAttributesOutput{}
					return user, nil
				},
				func(_ context.Context, _ *cognitoidentityprovider.AdminGetUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminGetUserOutput, error) {
					user := &cognitoidentityprovider.AdminGetUserOutput{
						UserAttributes: successfullyGetUser,
						UserStatus:     status,
						Username:       &userID,
						Enabled:        true,
					}
					return user, nil
				},
				func(_ context.Context, _ *cognitoidentityprovider.AdminEnableUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminEnableUserOutput, error) {
					return &cognitoidentityprovider.AdminEnableUserOutput{}, nil
				},
				func(_ context.Context, _ *cognitoidentityprovider.AdminDisableUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminDisableUserOutput, error) {
					//awsOrigErr := errors.New(awsErrCode)
					//awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					awsOrigErr := smithy.ErrorFault(1) // server error
					awsErr := &smithy.GenericAPIError{
						Code:    awsErrCode,
						Message: awsErrMessage,
						Fault:   awsOrigErr,
					}
					return nil, awsErr
				},
				forename,
				false,
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
				},
			},
			// 404 response from Cognito user update
			{
				func(_ context.Context, _ *cognitoidentityprovider.AdminUpdateUserAttributesInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminUpdateUserAttributesOutput, error) {
					//awsOrigErr := errors.New(awsUNFErrCode)
					//awsErr := awserr.New(awsUNFErrCode, awsUNFErrMessage, awsOrigErr)
					awsOrigErr := smithy.ErrorFault(2) // client error
					awsErr := &smithy.GenericAPIError{
						Code:    awsUNFErrCode,
						Message: awsUNFErrMessage,
						Fault:   awsOrigErr,
					}
					return nil, awsErr
				},
				func(_ context.Context, _ *cognitoidentityprovider.AdminGetUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminGetUserOutput, error) {
					user := &cognitoidentityprovider.AdminGetUserOutput{
						UserAttributes: successfullyGetUser,
						UserStatus:     status,
						Username:       &userID,
						Enabled:        true,
					}
					return user, nil
				},
				func(_ context.Context, _ *cognitoidentityprovider.AdminEnableUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminEnableUserOutput, error) {
					return &cognitoidentityprovider.AdminEnableUserOutput{}, nil
				},
				func(_ context.Context, _ *cognitoidentityprovider.AdminDisableUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminDisableUserOutput, error) {
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
				func(_ context.Context, _ *cognitoidentityprovider.AdminUpdateUserAttributesInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminUpdateUserAttributesOutput, error) {
					//awsOrigErr := errors.New(awsErrCode)
					//awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					awsOrigErr := smithy.ErrorFault(1) // server error
					awsErr := &smithy.GenericAPIError{
						Code:    awsErrCode,
						Message: awsErrMessage,
						Fault:   awsOrigErr,
					}
					return nil, awsErr
				},
				func(_ context.Context, _ *cognitoidentityprovider.AdminGetUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminGetUserOutput, error) {
					user := &cognitoidentityprovider.AdminGetUserOutput{
						UserAttributes: successfullyGetUser,
						UserStatus:     status,
						Username:       &userID,
						Enabled:        true,
					}
					return user, nil
				},
				func(_ context.Context, _ *cognitoidentityprovider.AdminEnableUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminEnableUserOutput, error) {
					return &cognitoidentityprovider.AdminEnableUserOutput{}, nil
				},
				func(_ context.Context, _ *cognitoidentityprovider.AdminDisableUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminDisableUserOutput, error) {
					return &cognitoidentityprovider.AdminDisableUserOutput{}, nil
				},
				forename,
				true,
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
				},
			},
			// reload user details failure
			{
				func(_ context.Context, _ *cognitoidentityprovider.AdminUpdateUserAttributesInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminUpdateUserAttributesOutput, error) {
					user := &cognitoidentityprovider.AdminUpdateUserAttributesOutput{}
					return user, nil
				},
				func(_ context.Context, _ *cognitoidentityprovider.AdminGetUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminGetUserOutput, error) {
					//awsOrigErr := errors.New(awsUNFErrCode)
					//awsErr := awserr.New(awsUNFErrCode, awsUNFErrMessage, awsOrigErr)
					awsOrigErr := smithy.ErrorFault(2) // client error
					awsErr := &smithy.GenericAPIError{
						Code:    awsUNFErrCode,
						Message: awsUNFErrMessage,
						Fault:   awsOrigErr,
					}
					return nil, awsErr
				},
				func(_ context.Context, _ *cognitoidentityprovider.AdminEnableUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminEnableUserOutput, error) {
					return &cognitoidentityprovider.AdminEnableUserOutput{}, nil
				},
				func(_ context.Context, _ *cognitoidentityprovider.AdminDisableUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminDisableUserOutput, error) {
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
		//awsOrigErr := errors.New(awsUNFErrCode)
		//awsErr := awserr.New(awsUNFErrCode, awsUNFErrMessage, awsOrigErr)
		awsOrigErr := smithy.ErrorFault(2) // client error
		awsErr := &smithy.GenericAPIError{
			Code:    awsUNFErrCode,
			Message: awsUNFErrMessage,
			Fault:   awsOrigErr,
		}
		errResponse := processUpdateCognitoError(ctx, awsErr, "Testing user not found")
		So(errResponse.Status, ShouldEqual, http.StatusNotFound)
		castErr := errResponse.Errors[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.UserNotFoundError)
		So(castErr.Description, ShouldEqual, "user could not be found")
	})

	Convey("Processes InternalError to a 500 response", t, func() {
		awsErrCode := "InternalErrorException"
		awsErrMessage := "something went wrong"
		//awsOrigErr := errors.New(awsErrCode)
		//awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
		awsOrigErr := smithy.ErrorFault(1) // server error
		awsErr := &smithy.GenericAPIError{
			Code:    awsErrCode,
			Message: awsErrMessage,
			Fault:   awsOrigErr,
		}
		errResponse := processUpdateCognitoError(ctx, awsErr, "Testing internal error")
		So(errResponse.Status, ShouldEqual, http.StatusInternalServerError)
		castErr := errResponse.Errors[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InternalError)
		So(castErr.Description, ShouldEqual, "something went wrong")
	})

	Convey("Processes InvalidField to a 400 response", t, func() {
		awsErrCode := "InvalidParameterException"
		awsErrMessage := "param invalid"
		//awsOrigErr := errors.New(awsErrCode)
		//awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
		awsOrigErr := smithy.ErrorFault(2) // client error
		awsErr := &smithy.GenericAPIError{
			Code:    awsErrCode,
			Message: awsErrMessage,
			Fault:   awsOrigErr,
		}
		errResponse := processUpdateCognitoError(ctx, awsErr, "Testing invalid param error")
		So(errResponse.Status, ShouldEqual, http.StatusBadRequest)
		castErr := errResponse.Errors[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidFieldError)
		So(castErr.Description, ShouldEqual, "param invalid")
	})
}

func TestChangePasswordHandler(t *testing.T) {
	var (
		ctx                                      = context.Background()
		email, password, session                 = "foo_bar123@ext.ons.gov.uk", "Password2", "auth-challenge-session"
		accessToken, idToken, refreshToken       = "aaaa.bbbb.cccc", "llll.mmmm.nnnn", "zzzz.yyyy.xxxx.wwww.vvvv"
		expireLength                       int32 = 500
	)

	api, w, m := apiSetup()

	m.DescribeUserPoolClientFunc = func(_ context.Context, _ *cognitoidentityprovider.DescribeUserPoolClientInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.DescribeUserPoolClientOutput, error) {
		tokenValidDays := int32(1)
		refreshTokenUnits := types.TimeUnitsTypeDays

		userPoolClient := &cognitoidentityprovider.DescribeUserPoolClientOutput{
			UserPoolClient: &types.UserPoolClientType{
				RefreshTokenValidity: tokenValidDays,
				TokenValidityUnits: &types.TokenValidityUnitsType{
					RefreshToken: refreshTokenUnits,
				},
			},
		}
		return userPoolClient, nil
	}

	Convey("RespondToAuthChallenge - check expected responses", t, func() {
		respondToAuthChallengeTests := []struct {
			respondToAuthChallengeFunction func(_ context.Context, input *cognitoidentityprovider.RespondToAuthChallengeInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.RespondToAuthChallengeOutput, error)
			assertions                     func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse)
		}{
			{
				// Cognito successful password change
				func(_ context.Context, _ *cognitoidentityprovider.RespondToAuthChallengeInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.RespondToAuthChallengeOutput, error) {
					return &cognitoidentityprovider.RespondToAuthChallengeOutput{
						AuthenticationResult: &types.AuthenticationResultType{
							AccessToken:  &accessToken,
							ExpiresIn:    expireLength,
							IdToken:      &idToken,
							RefreshToken: &refreshToken,
						},
					}, nil
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					var responseBody map[string]interface{}
					_ = json.Unmarshal(successResponse.Body, &responseBody)
					So(successResponse, ShouldNotBeNil)
					So(successResponse.Status, ShouldEqual, http.StatusAccepted)
					So(errorResponse, ShouldBeNil)
					So(responseBody["expirationTime"], ShouldNotBeNil)
					So(responseBody["refreshTokenExpirationTime"], ShouldNotBeNil)
				},
			},
			{
				// Cognito internal error
				func(_ context.Context, _ *cognitoidentityprovider.RespondToAuthChallengeInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.RespondToAuthChallengeOutput, error) {
					//awsOrigErr := errors.New(awsErrCode)
					//awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					awsOrigErr := smithy.ErrorFault(1) // server error
					awsErr := &smithy.GenericAPIError{
						Code:    awsErrCode,
						Message: awsErrMessage,
						Fault:   awsOrigErr,
					}
					return nil, awsErr
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusInternalServerError)
				},
			},
			{
				// Cognito invalid session
				func(_ context.Context, _ *cognitoidentityprovider.RespondToAuthChallengeInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.RespondToAuthChallengeOutput, error) {
					awsErrCode := "CodeMismatchException"
					awsErrMessage := "session invalid"
					//awsOrigErr := errors.New(awsErrCode)
					//awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					awsOrigErr := smithy.ErrorFault(1) // server error
					awsErr := &smithy.GenericAPIError{
						Code:    awsErrCode,
						Message: awsErrMessage,
						Fault:   awsOrigErr,
					}
					return nil, awsErr
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
				},
			},
			{
				// Cognito invalid password
				func(_ context.Context, _ *cognitoidentityprovider.RespondToAuthChallengeInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.RespondToAuthChallengeOutput, error) {
					awsErrCode := "InvalidPasswordException"
					awsErrMessage := "password invalid"
					//awsOrigErr := errors.New(awsErrCode)
					//awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					awsOrigErr := smithy.ErrorFault(2) // client error
					awsErr := &smithy.GenericAPIError{
						Code:    awsErrCode,
						Message: awsErrMessage,
						Fault:   awsOrigErr,
					}
					return nil, awsErr
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldBeNil)
					So(errorResponse.Status, ShouldEqual, http.StatusBadRequest)
				},
			},
			{
				// Cognito invalid user
				func(_ context.Context, _ *cognitoidentityprovider.RespondToAuthChallengeInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.RespondToAuthChallengeOutput, error) {
					awsErrCode := "UserNotFoundException"
					awsErrMessage := "user not found"
					//awsOrigErr := errors.New(awsErrCode)
					//awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					awsOrigErr := smithy.ErrorFault(2) // client error
					awsErr := &smithy.GenericAPIError{
						Code:    awsErrCode,
						Message: awsErrMessage,
						Fault:   awsOrigErr,
					}
					return nil, awsErr
				},
				func(successResponse *models.SuccessResponse, errorResponse *models.ErrorResponse) {
					So(successResponse, ShouldNotBeNil)
					So(successResponse.Status, ShouldEqual, http.StatusAccepted)
					So(errorResponse, ShouldBeNil)
				},
			},
		}
		for _, tt := range respondToAuthChallengeTests {
			m.RespondToAuthChallengeFunc = tt.respondToAuthChallengeFunction
			postBody := map[string]interface{}{"type": models.NewPasswordRequiredType, "email": email, "password": password, "session": session}
			body, _ := json.Marshal(postBody)
			r := httptest.NewRequest(http.MethodPut, changePasswordEndPoint, bytes.NewReader(body))
			successResponse, errorResponse := api.ChangePasswordHandler(ctx, w, r)
			tt.assertions(successResponse, errorResponse)
		}
	})
}

func TestConfirmForgotPasswordChangePasswordHandler(t *testing.T) {
	var (
		ctx               = context.Background()
		email             = "fred.bloggs@ons.gov.uk"
		password          = "Password2@123456"
		verificationToken = "999999"
	)

	api, w, m := apiSetup()
	Convey("ConfirmForgotPassword - check expected responses", t, func() {
		confirmForgotPasswordTests := []struct {
			confirmForgotPasswordFunction func(_ context.Context, _ *cognitoidentityprovider.ConfirmForgotPasswordInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ConfirmForgotPasswordOutput, error)
			httpResponse                  int
		}{
			// Cognito successful password change
			{
				func(_ context.Context, _ *cognitoidentityprovider.ConfirmForgotPasswordInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ConfirmForgotPasswordOutput, error) {
					tst := cognitoidentityprovider.ConfirmForgotPasswordOutput{}
					return &tst, nil
				},
				http.StatusAccepted,
			},
			// Cognito internal error
			{
				func(_ context.Context, _ *cognitoidentityprovider.ConfirmForgotPasswordInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ConfirmForgotPasswordOutput, error) {
					//awsOrigErr := errors.New(awsErrCode)
					//awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					awsOrigErr := smithy.ErrorFault(1) // server error
					awsErr := &smithy.GenericAPIError{
						Code:    awsErrCode,
						Message: awsErrMessage,
						Fault:   awsOrigErr,
					}
					return nil, awsErr
				},
				http.StatusInternalServerError,
			},
			// Cognito invalid token
			{
				func(_ context.Context, _ *cognitoidentityprovider.ConfirmForgotPasswordInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ConfirmForgotPasswordOutput, error) {
					awsErrCode := "CodeMismatchException"
					awsErrMessage := "session invalid"
					//awsOrigErr := errors.New(awsErrCode)
					//awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					awsOrigErr := smithy.ErrorFault(0) // unknown error
					awsErr := &smithy.GenericAPIError{
						Code:    awsErrCode,
						Message: awsErrMessage,
						Fault:   awsOrigErr,
					}
					return nil, awsErr
				},
				http.StatusBadRequest,
			},
			// Cognito expired token
			{
				func(_ context.Context, _ *cognitoidentityprovider.ConfirmForgotPasswordInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ConfirmForgotPasswordOutput, error) {
					awsErrCode := "ExpiredCodeException"
					awsErrMessage := "token expired"
					//awsOrigErr := errors.New(awsErrCode)
					//awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					awsOrigErr := smithy.ErrorFault(2) // client error
					awsErr := &smithy.GenericAPIError{
						Code:    awsErrCode,
						Message: awsErrMessage,
						Fault:   awsOrigErr,
					}
					return nil, awsErr
				},
				http.StatusBadRequest,
			},
			// Cognito invalid password
			{
				func(_ context.Context, _ *cognitoidentityprovider.ConfirmForgotPasswordInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ConfirmForgotPasswordOutput, error) {
					awsErrCode := "InvalidPasswordException"
					awsErrMessage := "password invalid"
					//awsOrigErr := errors.New(awsErrCode)
					//awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					awsOrigErr := smithy.ErrorFault(2) // client error
					awsErr := &smithy.GenericAPIError{
						Code:    awsErrCode,
						Message: awsErrMessage,
						Fault:   awsOrigErr,
					}
					return nil, awsErr
				},
				http.StatusBadRequest,
			},
			// Cognito invalid user
			{
				func(_ context.Context, _ *cognitoidentityprovider.ConfirmForgotPasswordInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ConfirmForgotPasswordOutput, error) {
					awsErrCode := "UserNotFoundException"
					awsErrMessage := "user not found"
					//awsOrigErr := errors.New(awsErrCode)
					//awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					awsOrigErr := smithy.ErrorFault(2) // client error
					awsErr := &smithy.GenericAPIError{
						Code:    awsErrCode,
						Message: awsErrMessage,
						Fault:   awsOrigErr,
					}
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
				models.InvalidUserIDError,
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
		ctx   = context.Background()
		email = "foo_bar123@ext.ons.gov.uk"
	)

	api, w, m := apiSetup()

	Convey("ForgotPassword - check expected responses", t, func() {
		respondToAuthChallengeTests := []struct {
			forgotPasswordFunction func(_ context.Context, input *cognitoidentityprovider.ForgotPasswordInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ForgotPasswordOutput, error)
			httpResponse           int
		}{
			{
				// Cognito successful password change
				func(_ context.Context, _ *cognitoidentityprovider.ForgotPasswordInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ForgotPasswordOutput, error) {
					return &cognitoidentityprovider.ForgotPasswordOutput{
						CodeDeliveryDetails: &types.CodeDeliveryDetailsType{},
					}, nil
				},
				http.StatusAccepted,
			},
			{
				// Cognito internal error
				func(_ context.Context, _ *cognitoidentityprovider.ForgotPasswordInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ForgotPasswordOutput, error) {
					//awsOrigErr := errors.New(awsErrCode)
					//awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					awsOrigErr := smithy.ErrorFault(1) // server error
					awsErr := &smithy.GenericAPIError{
						Code:    awsErrCode,
						Message: awsErrMessage,
						Fault:   awsOrigErr,
					}
					return nil, awsErr
				},
				http.StatusInternalServerError,
			},
			{
				// Cognito too many requests
				func(_ context.Context, _ *cognitoidentityprovider.ForgotPasswordInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ForgotPasswordOutput, error) {
					awsErrCode := "TooManyRequestsException"
					awsErrMessage := "slow down"
					//awsOrigErr := errors.New(awsErrCode)
					//awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					awsOrigErr := smithy.ErrorFault(2) // client error
					awsErr := &smithy.GenericAPIError{
						Code:    awsErrCode,
						Message: awsErrMessage,
						Fault:   awsOrigErr,
					}
					return nil, awsErr
				},
				http.StatusBadRequest,
			},
			{
				// Cognito invalid user
				func(_ context.Context, _ *cognitoidentityprovider.ForgotPasswordInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ForgotPasswordOutput, error) {
					awsErrCode := "UserNotFoundException"
					awsErrMessage := "user not found in user pool"
					//awsOrigErr := errors.New(awsErrCode)
					//awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					awsOrigErr := smithy.ErrorFault(2) // client error
					awsErr := &smithy.GenericAPIError{
						Code:    awsErrCode,
						Message: awsErrMessage,
						Fault:   awsOrigErr,
					}
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
		groups    = []types.GroupType{
			{
				CreationDate:     &timestamp,
				Description:      aws.String("A test group1"),
				GroupName:        aws.String("test-group1"),
				LastModifiedDate: &timestamp,
				Precedence:       aws.Int32(4),
				RoleArn:          aws.String(""),
				UserPoolId:       aws.String(""),
			},
			{
				CreationDate:     &timestamp,
				Description:      aws.String("A test group1"),
				GroupName:        aws.String("test-group1"),
				LastModifiedDate: &timestamp,
				Precedence:       aws.Int32(4),
				RoleArn:          aws.String(""),
				UserPoolId:       aws.String(""),
			},
		}
	)

	api, w, m := apiSetup()

	Convey("List groups for user -check expected responses", t, func() {
		listusergroups := []struct {
			getUserGroupsFunction func(_ context.Context, input *cognitoidentityprovider.AdminListGroupsForUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminListGroupsForUserOutput, error)
			httpResponse          int
		}{
			{
				// 200 response from Cognito with empty NextToken
				func(_ context.Context, _ *cognitoidentityprovider.AdminListGroupsForUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminListGroupsForUserOutput, error) {
					return &cognitoidentityprovider.AdminListGroupsForUserOutput{
						Groups:    groups,
						NextToken: nil,
					}, nil
				},
				http.StatusOK,
			},
			{
				// 200 response from Cognito with empty NextToken
				func(_ context.Context, _ *cognitoidentityprovider.AdminListGroupsForUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminListGroupsForUserOutput, error) {
					return &cognitoidentityprovider.AdminListGroupsForUserOutput{
						Groups:    []types.GroupType{},
						NextToken: nil,
					}, nil
				},
				http.StatusOK,
			},
			{
				// 500 response from Cognito
				func(_ context.Context, _ *cognitoidentityprovider.AdminListGroupsForUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminListGroupsForUserOutput, error) {
					//awsOrigErr := errors.New(awsErrCode)
					//awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
					awsOrigErr := smithy.ErrorFault(1) // server error
					awsErr := &smithy.GenericAPIError{
						Code:    awsErrCode,
						Message: awsErrMessage,
						Fault:   awsOrigErr,
					}
					return nil, awsErr
				},
				http.StatusInternalServerError,
			},
			{
				// 404 response from Cognito
				func(_ context.Context, _ *cognitoidentityprovider.AdminListGroupsForUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminListGroupsForUserOutput, error) {
					//awsOrigErr := errors.New(awsUNFErrCode)
					//awsErr := awserr.New(awsUNFErrCode, awsUNFErrMessage, awsOrigErr)
					awsOrigErr := smithy.ErrorFault(2) // client error
					awsErr := &smithy.GenericAPIError{
						Code:    awsUNFErrCode,
						Message: awsUNFErrMessage,
						Fault:   awsOrigErr,
					}
					return nil, awsErr
				},
				http.StatusInternalServerError,
			},
		}

		for _, tt := range listusergroups {
			m.ListGroupsForUserFunc = tt.getUserGroupsFunction

			r := httptest.NewRequest(http.MethodGet, userListGroupsEndPoint, http.NoBody)

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
		userNotFoundDescription = "User not found"
		userID                  = models.UserParams{
			ID: "abcd1234",
		}
		group0 = "test_group_0"
		group1 = "test_group_1"
	)

	listOfGroups := []types.GroupType{
		{
			GroupName: &group0,
		},
	}

	api, _, m := apiSetup()
	Convey("error is returned when list groups for a user returns an error", t, func() {
		m.ListGroupsForUserFunc = func(_ context.Context, _ *cognitoidentityprovider.AdminListGroupsForUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminListGroupsForUserOutput, error) {
			var userNotFoundException types.ResourceNotFoundException
			userNotFoundException.Message = &userNotFoundDescription
			return nil, &userNotFoundException
		}

		listGroupsforUserResponse, errorResponse := api.getGroupsForUser(ctx, nil, userID)

		So(listGroupsforUserResponse, ShouldBeNil)
		So(errorResponse.Error(), ShouldResemble, "ResourceNotFoundException: User not found")
	})

	Convey("When there is no next token cognito is called once and the list of groups in returned", t, func() {
		listOfGroupsForUser := []types.GroupType{
			{
				GroupName: &group0,
			},
		}

		m.ListGroupsForUserFunc = func(_ context.Context, _ *cognitoidentityprovider.AdminListGroupsForUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminListGroupsForUserOutput, error) {
			listGroupsForUser := &cognitoidentityprovider.AdminListGroupsForUserOutput{
				Groups: []types.GroupType{
					{
						GroupName: &group0,
					},
				},
			}
			return listGroupsForUser, nil
		}

		listOfUsersResponse, errorResponse := api.getGroupsForUser(ctx, nil, userID)

		So(listOfUsersResponse, ShouldResemble, listOfGroupsForUser)

		So(errorResponse, ShouldBeNil)
	})

	Convey("When there is a next token cognito is called more than once and the appended list of users in returned", t, func() {
		listOfGroupsForUser := []types.GroupType{
			{
				GroupName: &group0,
			},
			{
				GroupName: &group0,
			},
			{
				GroupName: &group1,
			},
		}

		m.ListGroupsForUserFunc = func(_ context.Context, input *cognitoidentityprovider.AdminListGroupsForUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminListGroupsForUserOutput, error) {
			nextToken := "nextToken"

			if input.NextToken != nil {
				listGroupsForUser := &cognitoidentityprovider.AdminListGroupsForUserOutput{
					NextToken: nil,
					Groups: []types.GroupType{
						{
							GroupName: &group1,
						},
					},
				}
				return listGroupsForUser, nil
			}
			listGroupsForUser := &cognitoidentityprovider.AdminListGroupsForUserOutput{
				NextToken: &nextToken,
				Groups: []types.GroupType{
					{
						GroupName: &group0,
					},
				},
			}
			return listGroupsForUser, nil
		}

		listGroupsForUserResponse, errorResponse := api.getGroupsForUser(ctx, listOfGroups, userID)

		So(listGroupsForUserResponse, ShouldResemble, listOfGroupsForUser)
		So(errorResponse, ShouldBeNil)
	})

	Convey("When GetGroupsforUser in called with a list of groups the appended list of groups in returned", t, func() {
		listOfGroups := []types.GroupType{
			{
				GroupName: &group0,
			},
		}

		returnedlistOfGroups := []types.GroupType{
			{
				GroupName: &group0,
			},
			{
				GroupName: &group0,
			},
		}

		m.ListGroupsForUserFunc = func(_ context.Context, _ *cognitoidentityprovider.AdminListGroupsForUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminListGroupsForUserOutput, error) {
			listGroupsForUser := &cognitoidentityprovider.AdminListGroupsForUserOutput{
				Groups: []types.GroupType{
					{
						GroupName: &group0,
					},
				},
			}
			return listGroupsForUser, nil
		}

		listGroupsForUseResponse, errorResponse := api.getGroupsForUser(ctx, listOfGroups, userID)

		So(listGroupsForUseResponse, ShouldResemble, returnedlistOfGroups)
		So(errorResponse, ShouldBeNil)
	})
}

func TestIsValidFilter(t *testing.T) {
	api, _, _ := apiSetup()

	Convey("Validate Filter - check expected responses", t, func() {
		validateFilterTest := []struct {
			description string
			path        string
			query       string
			assertions  func(successResponse string, errorResponse error)
		}{
			{
				"active true",
				"/v1/users",
				"active=true",
				func(successResponse string, errorResponse error) {
					So(errorResponse, ShouldBeNil)
					So(successResponse, ShouldNotBeNil)
					So(successResponse, ShouldResemble, "status=\"Enabled\"")
				},
			},
			{
				"active false",
				"/v1/users",
				"active=false",
				func(successResponse string, errorResponse error) {
					So(errorResponse, ShouldBeNil)
					So(successResponse, ShouldNotBeNil)
					So(successResponse, ShouldResemble, "status=\"Disabled\"")
				},
			},
			{
				"active another string",
				"v1/user",
				"active=string",
				func(successResponse string, errorResponse error) {
					So(errorResponse, ShouldNotBeNil)
					So(successResponse, ShouldBeEmpty)
					So(errorResponse.Error(), ShouldResemble, "InvalidFilterQuery")
					castErr := errorResponse.(*models.Error)
					So(castErr.Code, ShouldEqual, models.InvalidFilterQuery)
					So(castErr.Description, ShouldEqual, models.InvalidFilterQueryDescription)
				},
			},
			{
				"active another path",
				"v1/group",
				"active=true",
				func(successResponse string, errorResponse error) {
					So(errorResponse, ShouldNotBeNil)
					So(successResponse, ShouldBeEmpty)
					castErr := errorResponse.(*models.Error)
					So(castErr.Code, ShouldEqual, models.InvalidFilterQuery)
					So(castErr.Description, ShouldEqual, models.InvalidFilterQueryDescription)
				},
			},
		}

		for _, tt := range validateFilterTest {
			Convey(tt.description, func() {
				successResponse, errorResponse := api.GetFilterStringAndValidate(tt.path, tt.query)
				tt.assertions(successResponse, errorResponse)
			},
			)
		}
	})
}
