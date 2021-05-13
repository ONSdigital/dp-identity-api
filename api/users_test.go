package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-identity-api/apierrors"
	"github.com/ONSdigital/dp-identity-api/cognito/mock"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

const usersEndPoint = "http://localhost:25600/users"

func TestCreateUserHandler(t *testing.T) {

	var (
		routeMux                          = mux.NewRouter()
		ctx                               = context.Background()
		name, status, email, poolId, userException, clientId, clientSecret string = "Foo_Bar", "UNCONFIRMED", "foo_bar123@foobar.io.me", "us-west-11_bxushuds", "UsernameExistsException: User account already exists", "abc123", "bsjahsaj9djsiq"
	)

	m := &mock.MockCognitoIdentityProviderClient{}

	Convey("Admin create user - check expected responses", t, func() {
		adminCreateUsersTests := []struct {
			createUsersFunction func(userInput *cognitoidentityprovider.AdminCreateUserInput) (*cognitoidentityprovider.AdminCreateUserOutput, error)
			httpResponse int
		}{
			{
				// 201 response - user created
				func(userInput *cognitoidentityprovider.AdminCreateUserInput) (*cognitoidentityprovider.AdminCreateUserOutput, error){
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
				// 400 response - user already exists
				func(userInput *cognitoidentityprovider.AdminCreateUserInput) (*cognitoidentityprovider.AdminCreateUserOutput, error){
					var userExistsException cognitoidentityprovider.UsernameExistsException
					userExistsException.Message_ = &userException
					userExistsException.RespMetadata.StatusCode = http.StatusBadRequest
		
					return nil, &userExistsException
				},
				http.StatusBadRequest,
			},
			{
				// 500 response - internal error exception
				func(userInput *cognitoidentityprovider.AdminCreateUserInput) (*cognitoidentityprovider.AdminCreateUserOutput, error){
					var internaleErrorException cognitoidentityprovider.InternalErrorException
					internaleErrorException.Message_ = &userException
					internaleErrorException.RespMetadata.StatusCode = http.StatusInternalServerError
		
					return nil, &internaleErrorException
				},
				http.StatusInternalServerError,
			},
		}

		for _, tt := range adminCreateUsersTests {
			m.AdminCreateUserFunc = tt.createUsersFunction
			api := Setup(ctx, routeMux, m, poolId, clientId, clientSecret)

			postBody := map[string]interface{}{"username": name, "email": email}
	
			body, _ := json.Marshal(postBody)
	
			r := httptest.NewRequest(http.MethodPost, usersEndPoint, bytes.NewReader(body))
	
			w := httptest.NewRecorder()
	
			api.Router.ServeHTTP(w, r)
	
			So(w.Code, ShouldEqual, tt.httpResponse)
		}
	})

	Convey("Admin create user returns 500: error unmarshalling request body", t, func() {
		r := httptest.NewRequest(http.MethodPost, usersEndPoint, bytes.NewReader(nil))

		w := httptest.NewRecorder()

		api := Setup(ctx, routeMux, m, poolId, clientId, clientSecret)

		api.Router.ServeHTTP(w, r)

		errorBody, _ := ioutil.ReadAll(w.Body)
		var e models.ErrorStructure
		json.Unmarshal(errorBody, &e)

		So(w.Code, ShouldEqual, http.StatusInternalServerError)
		So(e.Errors[0].Message, ShouldEqual, "api endpoint POST user returned an error unmarshalling request body")
	})

	Convey("Validation fails 400: validating email and username throws validation errors", t, func() {
		userValidationTests := []struct {
			userDetails map[string]interface{}
			errorMessage []string
			httpResponse int
		}{
			// missing username
			{
				map[string]interface{}{"username": "", "email": email},
				[]string{
					apierrors.InvalidUserNameMessage,
				},		
				http.StatusBadRequest,
			},
			// missing email
			{
				map[string]interface{}{"username": name, "email": ""},
				[]string{
					apierrors.InvalidErrorMessage,
				},
				http.StatusBadRequest,
			},
			// missing both username and email
			{
				map[string]interface{}{"username": "", "email": ""},
				[]string{
					apierrors.InvalidUserNameMessage,
					apierrors.InvalidErrorMessage,
				},
				http.StatusBadRequest,
			},
		}
	
		for _, tt := range userValidationTests {
			body, _ := json.Marshal(tt.userDetails)
	
			r := httptest.NewRequest(http.MethodPost, usersEndPoint, bytes.NewReader(body))
	
			w := httptest.NewRecorder()
	
			api := Setup(ctx, routeMux, m, poolId, clientId, clientSecret)

			api.Router.ServeHTTP(w, r)
	
			errorBody, _ := ioutil.ReadAll(w.Body)
			var e models.ErrorStructure
			json.Unmarshal(errorBody, &e)
	
			So(w.Code, ShouldEqual, tt.httpResponse)
			So(len(e.Errors), ShouldEqual, len(tt.errorMessage))
			So(e.Errors[0].Message, ShouldEqual, tt.errorMessage[0])
			if len(e.Errors) > 1 {
				So(e.Errors[1].Message, ShouldEqual, tt.errorMessage[1])
			}
		}
	})
}