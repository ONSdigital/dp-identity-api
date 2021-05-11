package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-identity-api/apierrors"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPasswordHasBeenProvided(t *testing.T) {

	Convey("A password has been provided", t, func() {

		body := AuthParams{
			Password: "password",
		}

		passwordResponse := passwordValidation(body)
		// asstert 1
		So(passwordResponse, ShouldBeTrue)
	})

	Convey("There isn't a password field in the body and the password isn't validated", t, func() {

		body := AuthParams{}

		passwordResponse := passwordValidation(body)
		// asstert 2
		So(passwordResponse, ShouldBeFalse)
	})

	Convey("There isn't a password value in the body and the password isn't validated", t, func() {

		body := AuthParams{
			Password: "",
		}

		passwordResponse := passwordValidation(body)
		//asssert 3
		So(passwordResponse, ShouldBeFalse)
	})
}

func TestEmailConformsToExpectedFormat(t *testing.T) {

	Convey("The email conforms to the expected format and is validated", t, func() {

		body := AuthParams{
			Email: "email.email@ons.gov.uk",
		}

		emailResponse := emailValidation(body)
		// assert 4
		So(emailResponse, ShouldBeTrue)
	})

	Convey("There isn't an email field in the body and the email isn't validated", t, func() {

		body := AuthParams{
			Password: "password",
		}

		emailResponse := emailValidation(body)
		// assert 5
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

		var errorList []apierrors.IndividualError
		errorList = nil

		invalidEmailError := errors.New("Invalid email")
		invalidErrorMessage := "Unable to validate the email in the request"
		field := ""
		param := ""
		invalidEmailErrorBody := apierrors.IndividualErrorBuilder(invalidEmailError, invalidErrorMessage, field, param)
		errorList = append(errorList, invalidEmailErrorBody)
		errorList = append(errorList, invalidEmailErrorBody)

		ctx := context.Background()
		resp := httptest.NewRecorder()
		statusCode := 400
		errorResponseBody := apierrors.ErrorResponseBodyBuilder(errorList)

		writeErrorResponse(ctx, resp, statusCode, errorResponseBody)

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

		handleUnexpectedError(ctx, resp, unexpectedError, unexpectedErrorMessage, field, param)

		So(resp.Code, ShouldEqual, http.StatusInternalServerError)
		So(resp.Body.String(), ShouldResemble, errorResponseBodyExample)
	})
}
func TestCognitoResponseBuild(t *testing.T) {
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

func TestCognitoRequestHeaderBuild(t *testing.T) {
	Convey("build Cognito Request, an authParams and Config is processed and Cognito Request is built", t, func() {
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

		//So(reflect.TypeOf(w.Body.String()), ShouldResemble, true)
	})
}

func TestBuildJson(t *testing.T) {
	Convey("build json", t, func() {
		w := httptest.NewRecorder()
		ctx := context.Background()

		testBody := map[string]interface{}{"expirationTime": "123"}
		buildjson(testBody, w, ctx)
		So(w.Body.String(), ShouldResemble, "{\"expirationTime\":\"123\"}")

	})

	Convey("build json err", t, func() {
		w := httptest.NewRecorder()
		ctx := context.Background()

		testBody := map[string]interface{}{
			"foo": make(chan int),
		}
		buildjson(testBody, w, ctx)
		So(w.Body.String(), ShouldResemble, "failed to marshal the error\n")
	})

	// Convey("write HttpResponse err", t, func() {
	// 	w := httptest.NewRecorder()
	// 	ctx := context.Background()
	// 	testBody := map[string]interface{}{
	// 		"foo": "TestString",
	// 	}

	// 	buildjson(testBody, w, ctx)
	// 	So(w.Body.String(), ShouldResemble, "failed to write http response\n")
	// })
}
