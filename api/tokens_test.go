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

	"github.com/ONSdigital/dp-identity-api/apierrors"
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

		errorResponseBodyExample := `{"errors":[{"error":"Invalid email","message":"Unable to validate the email in the request","source":{"field":"","param":""}},{"error":"Invalid email","message":"Unable to validate the email in the request","source":{"field":"","param":""}}]}`

		var errorList []models.IndividualError
		errorList = nil
		
		errInvalidEmail := errors.New("Invalid email")
		invalidErrorMessage := "Unable to validate the email in the request"
		field := ""
		param := ""
		invalidEmailErrorBody := apierrors.IndividualErrorBuilder(errInvalidEmail, invalidErrorMessage, field, param)
		errorList = append(errorList, invalidEmailErrorBody)
		errorList = append(errorList, invalidEmailErrorBody)

		ctx := context.Background()
		resp := httptest.NewRecorder()
		statusCode := 400
		errorResponseBody := apierrors.ErrorResponseBodyBuilder(errorList)

		apierrors.WriteErrorResponse(ctx, resp, statusCode, errorResponseBody)

		So(resp.Code, ShouldEqual, http.StatusBadRequest)
		So(resp.Body.String(), ShouldResemble, errorResponseBodyExample)
	})
}

func TestHandleUnexpectedError(t *testing.T) {
	Convey("An error, an error message, a field and a param is logged and written to a http response", t, func() {

		errorResponseBodyExample := `{"errors":[{"error":"unexpected error","message":"something unexpected has happened","source":{"field":"","param":""}}]}`

		ctx := context.Background()
		unexpectedError := errors.New("unexpected error")
		unexpectedErrorMessage := "something unexpected has happened"
		field := ""
		param := ""

		resp := httptest.NewRecorder()

		apierrors.HandleUnexpectedError(ctx, resp, unexpectedError, unexpectedErrorMessage, field, param)

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
		So(w.Body.String(), ShouldResemble, "{\"errors\":[{\"error\":\"json: unsupported type: chan int\",\"message\":\"failed to marshal the error\",\"source\":{\"field\":\"\",\"param\":\"\"}}]}")
		So(w.Result().StatusCode, ShouldEqual, 500)
		So(w.Result().Header["Content-Type"], ShouldResemble, []string{"application/json"})
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
