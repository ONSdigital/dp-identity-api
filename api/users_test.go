package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-identity-api/cognito/mock"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

const usersEndPoint = "http://localhost:25600/users"

func TestCreateUserHandler(t *testing.T) {

	var (
		routeMux                                                                                                   = mux.NewRouter()
		ctx                                                                                                        = context.Background()
		name, surname, status, email, poolId, userException, clientId, clientSecret, authFlow, invalidEmail string = "bob", "bobbings", "UNCONFIRMED", "foo_bar123@ext.ons.gov.uk", "us-west-11_bxushuds", "UsernameExistsException: User account already exists", "abc123", "bsjahsaj9djsiq", "authflow", "foo_bar123@test.ons.gov.ie"
	)

	m := &mock.MockCognitoIdentityProviderClient{}

	api, _ := Setup(ctx, routeMux, m, poolId, clientId, clientSecret, authFlow)
	w := httptest.NewRecorder()

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

			postBody := map[string]interface{}{"forename": name, "surname": surname, "email": email}
			body, _ := json.Marshal(postBody)
			r := httptest.NewRequest(http.MethodPost, usersEndPoint, bytes.NewReader(body))

			successResponse, errorResponse := api.CreateUserHandler(w, r, ctx)

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

		successResponse, errorResponse := api.CreateUserHandler(w, r, ctx)

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
				map[string]interface{}{"forename": name, "surname": surname, "email": ""},
				[]string{
					models.InvalidEmailError,
				},
				http.StatusBadRequest,
			},
			// missing both forename and surname
			{
				map[string]interface{}{"forename": "", "surname": "", "email": email},
				[]string{
					models.InvalidForenameError,
					models.InvalidSurnameError,
				},
				http.StatusBadRequest,
			},
			// missing surname
			{
				map[string]interface{}{"forename": name, "surname": "", "email": email},
				[]string{
					models.InvalidSurnameError,
				},
				http.StatusBadRequest,
			},
			// missing forename
			{
				map[string]interface{}{"forename": "", "surname": surname, "email": email},
				[]string{
					models.InvalidForenameError,
				},
				http.StatusBadRequest,
			},
			// missing forename, surname and email
			{
				map[string]interface{}{"forename": "", "surname": "", "email": ""},
				[]string{
					models.InvalidForenameError,
					models.InvalidSurnameError,
					models.InvalidEmailError,
				},
				http.StatusBadRequest,
			},
			// invalid email
			{
				map[string]interface{}{"forename": name, "surname": surname, "email": invalidEmail},
				[]string{
					models.InvalidEmailError,
				},
				http.StatusBadRequest,
			},
		}

		for _, tt := range userValidationTests {
			body, _ := json.Marshal(tt.userDetails)
			r := httptest.NewRequest(http.MethodPost, usersEndPoint, bytes.NewReader(body))

			successResponse, errorResponse := api.CreateUserHandler(w, r, ctx)

			So(successResponse, ShouldBeNil)
			So(errorResponse.Status, ShouldEqual, tt.httpResponse)
			So(len(errorResponse.Errors), ShouldEqual, len(tt.errorCodes))
			castErr := errorResponse.Errors[0].(*models.Error)
			So(castErr.Code, ShouldEqual, tt.errorCodes[0])
			if len(errorResponse.Errors) > 1 {
				castErr = errorResponse.Errors[1].(*models.Error)
				So(castErr.Code, ShouldEqual, tt.errorCodes[1])
			}
		}
	})
}
