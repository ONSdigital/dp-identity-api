package api

import (
	"errors"
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

		email := "email@ons2.gov.uk"

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

func TestBuildingListOfErrors(t *testing.T) {

	Convey("An eror can be added to a list of errors", t, func() {

		individualErrorOne := IndividualError{
			SpecificError: "string, unchanging so devs can use this in code",
			Message:       "detailed explanation of error",
			Source: Source{
				Field: "reference to field like some.field or something",
				Param: "query param causing issue"},
		}

		individualErrorTwo := IndividualError{
			SpecificError: "string, unchanging so devs can use this in code",
			Message:       "detailed explanation of error two",
			Source: Source{
				Field: "reference to field like some.field or something",
				Param: "query param causing issue"},
		}

		listOfErrorsExample := []IndividualError{
			{
				SpecificError: "string, unchanging so devs can use this in code",
				Message:       "detailed explanation of error",
				Source: Source{
					Field: "reference to field like some.field or something",
					Param: "query param causing issue"},
			},
			{
				SpecificError: "string, unchanging so devs can use this in code",
				Message:       "detailed explanation of error two",
				Source: Source{
					Field: "reference to field like some.field or something",
					Param: "query param causing issue"},
			},
		}

		listOfErrors := errorListBuilder(nil, individualErrorOne)
		listOfErrors = errorListBuilder(listOfErrors, individualErrorTwo)

		So(listOfErrors, ShouldResemble, listOfErrorsExample)

	})
}

func TestBuildingErrorStructure(t *testing.T) {
	Convey("An error structure is created from a list of errors ", t, func() {

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
