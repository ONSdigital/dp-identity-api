package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-identity-api/cognito/mock"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gorilla/mux"

	"github.com/ONSdigital/dp-identity-api/apierrors"
	errModels "github.com/ONSdigital/dp-identity-api/models"

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

		var errorList []errModels.IndividualError
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

func TestSignOutHandler(t *testing.T) {
	var (
		r                                     = mux.NewRouter()
		ctx                                   = context.Background()
		poolId, clientId, clientSecret string = "us-west-11_bxushuds", "client-aaa-bbb", "secret-ccc-ddd"
	)

	m := &mock.MockCognitoIdentityProviderClient{}

	// mock call to: GlobalSignOut(input *cognitoidentityprovider.GlobalSignOutInput) (*cognitoidentityprovider.GlobalSignOutOutput, error)
	m.GlobalSignOutFunc = func(signOutInput *cognitoidentityprovider.GlobalSignOutInput) (*cognitoidentityprovider.GlobalSignOutOutput, error) {
		return &cognitoidentityprovider.GlobalSignOutOutput{}, nil
	}

	api, _ := Setup(ctx, r, m, poolId, clientId, clientSecret)

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
