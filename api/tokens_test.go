package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ONSdigital/dp-identity-api/apierrorsdeprecated"
	"github.com/ONSdigital/dp-identity-api/cognito/mock"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gorilla/mux"

	. "github.com/smartystreets/goconvey/convey"
)

const signOutEndPoint = "http://localhost:25600/tokens/self"

func TestPasswordHasBeenProvided(t *testing.T) {

	Convey("A password has been provided", t, func() {

		body := AuthParams{
			Password: "password",
		}

		passwordResponse := passwordValidation(body)
		So(passwordResponse, ShouldBeTrue)
	})

	Convey("There isn't a password field in the body and the password isn't validated", t, func() {

		body := AuthParams{}

		passwordResponse := passwordValidation(body)
		So(passwordResponse, ShouldBeFalse)
	})

	Convey("There isn't a password value in the body and the password isn't validated", t, func() {

		body := AuthParams{
			Password: "",
		}

		passwordResponse := passwordValidation(body)
		So(passwordResponse, ShouldBeFalse)
	})
}

func TestEmailConformsToExpectedFormat(t *testing.T) {

	Convey("The email conforms to the expected format and is validated", t, func() {

		body := AuthParams{
			Email: "email.email@ons.gov.uk",
		}

		emailResponse := emailValidation(body)
		So(emailResponse, ShouldBeTrue)
	})

	Convey("There isn't an email field in the body and the email isn't validated", t, func() {

		body := AuthParams{
			Password: "password",
		}

		emailResponse := emailValidation(body)
		So(emailResponse, ShouldBeFalse)
	})

	Convey("There isn't a email value in the body and the email isn't validated", t, func() {

		body := AuthParams{
			Email: "",
		}

		emailResponse := emailValidation(body)
		So(emailResponse, ShouldBeFalse)
	})

	Convey("The email doesn't conform to the expected format and it isn't validated", t, func() {

		body := AuthParams{
			Email: "email",
		}

		emailResponse := emailValidation(body)
		So(emailResponse, ShouldBeFalse)
	})
}

func TestWriteErrorResponse(t *testing.T) {
	Convey("A status code and an error body with two errors is written to a http response", t, func() {

		errorResponseBodyExample := `{"errors":[{"code":"Invalid email","description":"Unable to validate the email in the request"},{"code":"Invalid email","description":"Unable to validate the email in the request"}]}`

		var errorList []models.Error
		errorList = nil

		errInvalidEmail := errors.New("Invalid email")
		invalidErrorDescription := "Unable to validate the email in the request"
		invalidEmailErrorBody := apierrorsdeprecated.IndividualErrorBuilder(errInvalidEmail, invalidErrorDescription)
		errorList = append(errorList, invalidEmailErrorBody)
		errorList = append(errorList, invalidEmailErrorBody)

		ctx := context.Background()
		resp := httptest.NewRecorder()
		statusCode := 400
		errorResponseBody := apierrorsdeprecated.ErrorResponseBodyBuilder(errorList)

		apierrorsdeprecated.WriteErrorResponse(ctx, resp, statusCode, errorResponseBody)

		So(resp.Code, ShouldEqual, http.StatusBadRequest)
		So(resp.Body.String(), ShouldResemble, errorResponseBodyExample)
	})
}

func TestHandleUnexpectedError(t *testing.T) {
	Convey("An error and an error description is logged and written to a http response", t, func() {

		errorResponseBodyExample := `{"errors":[{"code":"unexpected error","description":"something unexpected has happened"}]}`

		ctx := context.Background()
		unexpectedError := errors.New("unexpected error")
		unexpectedErrorDescription := "something unexpected has happened"

		resp := httptest.NewRecorder()

		apierrorsdeprecated.HandleUnexpectedError(ctx, resp, unexpectedError, unexpectedErrorDescription)

		So(resp.Code, ShouldEqual, http.StatusInternalServerError)
		So(resp.Body.String(), ShouldResemble, errorResponseBodyExample)
	})
}
func TestCognitoRequestBuild(t *testing.T) {
	Convey("build Cognito Request, an authParams and Config is processed and Cognito Request is built", t, func() {

		authParams := AuthParams{
			Email:    "email.email@ons.gov.uk",
			Password: "password",
		}

		clientId := "awsclientid"
		clientSecret := "awsSectret"
		clientAuthFlow := "authflow"

		response := buildCognitoRequest(authParams, clientId, clientSecret, clientAuthFlow)

		So(*response.AuthParameters["USERNAME"], ShouldEqual, authParams.Email)
		So(*response.AuthParameters["PASSWORD"], ShouldEqual, authParams.Password)
		So(*response.AuthParameters["SECRET_HASH"], ShouldNotBeEmpty)
		So(*response.AuthFlow, ShouldResemble, "authflow")
		So(*response.ClientId, ShouldResemble, "awsclientid")
	})
}

func TestCognitoResponseHeaderBuild(t *testing.T) {
	Convey("build 201 response using an InitiateAuthOutput from Cognito", t, func() {
		w := httptest.NewRecorder()
		ctx := context.Background()
		accessToken := "accessToken"
		var expiration int64 = 123
		idToken := "idToken"
		Refresh := "refreshToken"

		initiateAuthOutput := &cognitoidentityprovider.InitiateAuthOutput{
			AuthenticationResult: &cognitoidentityprovider.AuthenticationResultType{
				AccessToken:  &accessToken,
				ExpiresIn:    &expiration,
				IdToken:      &idToken,
				RefreshToken: &Refresh,
			},
		}

		buildSucessfulResponse(initiateAuthOutput, w, ctx)

		So(w.Result().StatusCode, ShouldEqual, 201)
		So(w.Result().Header["Content-Type"], ShouldResemble, []string{"application/json"})
		So(w.Result().Header["Authorization"], ShouldResemble, []string{"Bearer " + accessToken})
		So(w.Result().Header["Id"], ShouldResemble, []string{idToken})
		So(w.Result().Header["Refresh"], ShouldResemble, []string{Refresh})

		var obj map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.String()), &obj)

		//there should be one entry in body
		So(len(obj), ShouldEqual, 1)

		type kv struct {
			Key   string
			Value interface{}
		}

		var ss []kv
		for k, v := range obj {
			ss = append(ss, kv{k, v})
		}
		str := fmt.Sprintf("%v", ss[0].Value)

		So(ss[0].Key, ShouldResemble, "expirationTime")
		So(str[:19], ShouldResemble, time.Now().UTC().Add(time.Second * 123).String()[:19])

	})

	Convey("build 500 response if the InitiateAuthOutput has an unexpected format", t, func() {
		w := httptest.NewRecorder()
		ctx := context.Background()

		initiateAuthOutput := &cognitoidentityprovider.InitiateAuthOutput{}
		buildSucessfulResponse(initiateAuthOutput, w, ctx)

		So(w.Result().StatusCode, ShouldEqual, 500)
	})
}

func TestBuildJson(t *testing.T) {
	w := httptest.NewRecorder()

	Convey("build json", t, func() {
		w := httptest.NewRecorder()
		ctx := context.Background()

		testBody := map[string]interface{}{"expirationTime": "123"}
		buildjson(testBody, w, ctx)
		So(w.Body.String(), ShouldResemble, "{\"expirationTime\":\"123\"}")

	})

	Convey("build json err", t, func() {

		ctx := context.Background()

		testBody := map[string]interface{}{
			"foo": make(chan int),
		}
		buildjson(testBody, w, ctx)
		So(w.Body.String(), ShouldResemble, "{\"errors\":[{\"code\":\"json: unsupported type: chan int\",\"description\":\"failed to marshal the error\"}]}")
		So(w.Result().StatusCode, ShouldEqual, 500)
		So(w.Result().Header["Content-Type"], ShouldResemble, []string{"application/json"})
	})
}

func TestBuildAdminSignoutRequest(t *testing.T) {
	Convey("build AdminSignout Request, an authParams and userPoolId is processed and AdminSignout Request is built", t, func() {
		authParams := AuthParams{
			Email:    "email.email@ons.gov.uk",
			Password: "password",
		}
		userPoolId := "userPoolId"

		response := buildAdminSignoutRequest(authParams, userPoolId)

		So(*response.Username, ShouldEqual, authParams.Email)
		So(*response.UserPoolId, ShouldEqual, userPoolId)
	})
}

func TestAdminUserGlobalSignOut(t *testing.T) {
	m := &mock.MockCognitoIdentityProviderClient{}
	authParams := AuthParams{
		Email:    "email.email@ons.gov.uk",
		Password: "password",
	}
	userPoolId := "userPoolId"

	Convey("Admin user global sign out returns an empty body signalling a 200 so no error is returned from the function", t, func() {

		// mock call to: AdminUserGlobalSignOut(adminUserGlobalSignOutInput *cognitoidentityprovider.AdminUserGlobalSignOutInput) (*cognitoidentityprovider.AdminUserGlobalSignOutOutput, error)
		m.AdminUserGlobalSignOutFunc = func(adminUserGlobalSignOutInput *cognitoidentityprovider.AdminUserGlobalSignOutInput) (*cognitoidentityprovider.AdminUserGlobalSignOutOutput, error) {
			return &cognitoidentityprovider.AdminUserGlobalSignOutOutput{}, nil
		}

		err := adminUserGlobalSignOut(authParams, userPoolId, m)

		So(err, ShouldBeNil)
	})

	Convey("Admin user global sign out returns an error so an error is returned from the function", t, func() {

		// mock call to: AdminUserGlobalSignOut(adminUserGlobalSignOutInput *cognitoidentityprovider.AdminUserGlobalSignOutInput) (*cognitoidentityprovider.AdminUserGlobalSignOutOutput, error)
		m.AdminUserGlobalSignOutFunc = func(adminUserGlobalSignOutInput *cognitoidentityprovider.AdminUserGlobalSignOutInput) (*cognitoidentityprovider.AdminUserGlobalSignOutOutput, error) {
			return nil, errors.New("InternalErrorException: Something went wrong")
		}

		err := adminUserGlobalSignOut(authParams, userPoolId, m)

		So(err, ShouldNotBeNil)
	})
}

func TestSignOutHandler(t *testing.T) {
	var (
		r                                                     = mux.NewRouter()
		ctx                                                   = context.Background()
		poolId, clientId, clientSecret, clientAuthFlow string = "us-west-11_bxushuds", "client-aaa-bbb", "secret-ccc-ddd", "authflow"
	)

	m := &mock.MockCognitoIdentityProviderClient{}

	// mock call to: GlobalSignOut(input *cognitoidentityprovider.GlobalSignOutInput) (*cognitoidentityprovider.GlobalSignOutOutput, error)
	m.GlobalSignOutFunc = func(signOutInput *cognitoidentityprovider.GlobalSignOutInput) (*cognitoidentityprovider.GlobalSignOutOutput, error) {
		return &cognitoidentityprovider.GlobalSignOutOutput{}, nil
	}

	api, _ := Setup(ctx, r, m, poolId, clientId, clientSecret, clientAuthFlow)

	Convey("Global Sign Out returns 204: successfully signed out user", t, func() {
		r := httptest.NewRequest(http.MethodDelete, signOutEndPoint, nil)
		r.Header.Set("Authorization", "Bearer zzzz-yyyy-xxxx")

		w := httptest.NewRecorder()

		api.Router.ServeHTTP(w, r)

		So(w.Code, ShouldEqual, http.StatusNoContent)
	})

	Convey("Global Sign Out returns 400: validate header structure", t, func() {
		headerValidationTests := []struct {
			authHeader string
		}{
			// missing Authorization header
			{
				"",
			},
			// malformed Authorization header
			{
				"Bearerzzzz-yyyy-xxxx",
			},
		}

		for _, tt := range headerValidationTests {
			r := httptest.NewRequest(http.MethodDelete, signOutEndPoint, nil)
			r.Header.Set("Authorization", tt.authHeader)

			w := httptest.NewRecorder()

			api.Router.ServeHTTP(w, r)

			So(w.Code, ShouldEqual, http.StatusBadRequest)
		}
	})

	Convey("Global Sign Out returns 500: Cognito internal error", t, func() {
		r := httptest.NewRequest(http.MethodDelete, signOutEndPoint, nil)
		r.Header.Set("Authorization", "Bearer zzzz-yyyy-xxxx")

		// mock failed call to: GlobalSignOut(input *cognitoidentityprovider.GlobalSignOutInput) (*cognitoidentityprovider.GlobalSignOutOutput, error)
		m.GlobalSignOutFunc = func(signOutInput *cognitoidentityprovider.GlobalSignOutInput) (*cognitoidentityprovider.GlobalSignOutOutput, error) {
			return nil, errors.New("InternalErrorException: Something went wrong")
		}

		w := httptest.NewRecorder()

		api.Router.ServeHTTP(w, r)

		So(w.Code, ShouldEqual, http.StatusInternalServerError)
	})

	Convey("Global Sign Out returns 400: request error", t, func() {
		r := httptest.NewRequest(http.MethodDelete, signOutEndPoint, nil)
		r.Header.Set("Authorization", "Bearer zzzz-yyyy-xxxx")

		// mock failed call to: GlobalSignOut(input *cognitoidentityprovider.GlobalSignOutInput) (*cognitoidentityprovider.GlobalSignOutOutput, error)
		m.GlobalSignOutFunc = func(signOutInput *cognitoidentityprovider.GlobalSignOutInput) (*cognitoidentityprovider.GlobalSignOutOutput, error) {
			return nil, errors.New("NotAuthorizedException: User is not authorized")
		}

		w := httptest.NewRecorder()

		api.Router.ServeHTTP(w, r)

		So(w.Code, ShouldEqual, http.StatusBadRequest)
	})
}
