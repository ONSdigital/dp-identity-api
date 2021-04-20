package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPasswordHasBeenProvided(t *testing.T) {

	Convey("A password has been provided", t, func() {
		password := "password"

		body := make(map[string]string)
		body["password"] = password

		passwordResponse := passwordValidation(body)

		So(passwordResponse, ShouldBeTrue)
	})

	Convey("There isn't a password field in the body and the password isn't validated", t, func() {

		body := make(map[string]string)

		passwordResponse := passwordValidation(body)
		So(passwordResponse, ShouldBeFalse)
	})

	Convey("There isn't a password value in the body and the password isn't validated", t, func() {

		password := ""
		body := make(map[string]string)
		body["password"] = password

		passwordResponse := passwordValidation(body)
		So(passwordResponse, ShouldBeFalse)
	})
}

func TestEmailConformsToExpectedFormat(t *testing.T) {

	Convey("The email conforms to the expected format and is validated", t, func() {

		email := "email.email@ons.gov.uk"

		body := make(map[string]string)
		body["email"] = email

		emailResponse := emailValidation(body)
		So(emailResponse, ShouldBeTrue)
	})

	Convey("There isn't an email field in the body and the email isn't validated", t, func() {

		password := "password"
		body := make(map[string]string)
		body["password"] = password

		emailResponse := emailValidation(body)
		So(emailResponse, ShouldBeFalse)
	})

	Convey("There isn't a email value in the body and the email isn't validated", t, func() {

		email := ""
		body := make(map[string]string)
		body["email"] = email

		passwordResponse := emailValidation(body)
		So(passwordResponse, ShouldBeFalse)
	})

	Convey("The email doesn't conform to the expected format and it isn't validated", t, func() {

		email := "email"

		body := make(map[string]string)
		body["email"] = email

		emailResponse := emailValidation(body)
		So(emailResponse, ShouldBeFalse)
	})
}

func TestBuildingIndividualErrors(t *testing.T) {

	Convey("The individual error conforms to the expected structure", t, func() {

		err := errors.New("string, unchanging so devs can use this in code")
		message := "detailed explanation of error"
		sourceField := "reference to field like some.field or something"
		sourceParam := "query param causing issue"

		individualErrorExample := IndividualError{
			SpecificError: "string, unchanging so devs can use this in code",
			Message:       "detailed explanation of error",
			Source: Source{
				Field: "reference to field like some.field or something",
				Param: "query param causing issue"},
		}

		individualError := individualErrorBuilder(err, message, sourceField, sourceParam)

		So(individualError, ShouldResemble, individualErrorExample)

	})
}

func TestBuildingErrorStructure(t *testing.T) {
	Convey("An error structure is created from a list of errors", t, func() {

		listOfErrors := []IndividualError{
			{
				SpecificError: "string, unchanging so devs can use this in code",
				Message:       "detailed explanation of error",
				Source: Source{
					Field: "reference to field like some.field or something",
					Param: "query param causing issue"},
			},
		}

		errorResponseBodyExample := ErrorStructure{
			Errors: []IndividualError{
				{
					SpecificError: "string, unchanging so devs can use this in code",
					Message:       "detailed explanation of error",
					Source: Source{
						Field: "reference to field like some.field or something",
						Param: "query param causing issue"}},
			},
		}

		errorResponseBody := errorResponseBodyBuilder(listOfErrors)

		So(errorResponseBody, ShouldResemble, errorResponseBodyExample)
	})
}

func TestWriteErrorResponse(t *testing.T) {
	Convey("A status code and an error body with two errors is written to a http response", t, func() {

		errorResponseBodyExample := `{"errors":[{"error":"Invalid email","message":"Unable to validate the email in the request","source":{"field":"","param":""}},{"error":"Invalid email","message":"Unable to validate the email in the request","source":{"field":"","param":""}}]}`

		invalidEmailError := errors.New("Invalid email")
		invalidErrorMessage := "Unable to validate the email in the request"
		field := ""
		param := ""
		invalidEmailErrorBody := individualErrorBuilder(invalidEmailError, invalidErrorMessage, field, param)
		errorList = append(errorList, invalidEmailErrorBody)
		errorList = append(errorList, invalidEmailErrorBody)

		resp := httptest.NewRecorder()
		statusCode := 400
		errorResponseBody := errorResponseBodyBuilder(errorList)

		writeErrorResponse(resp, statusCode, errorResponseBody)

		So(resp.Code, ShouldEqual, http.StatusBadRequest)
		So(resp.Body.String(), ShouldResemble, errorResponseBodyExample)
	})
}
