package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-identity-api/apierrors"

	. "github.com/smartystreets/goconvey/convey"
)

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
