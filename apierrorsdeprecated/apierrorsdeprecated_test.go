package apierrorsdeprecated

import (
	"errors"
	"testing"

	errModels "github.com/ONSdigital/dp-identity-api/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBuildingIndividualErrors(t *testing.T) {

	Convey("The individual error conforms to the expected structure", t, func() {

		err := errors.New("string, unchanging so devs can use this in code")
		message := "detailed explanation of error"
		sourceField := "reference to field like some.field or something"
		sourceParam := "query param causing issue"

		individualErrorExample := errModels.IndividualError{
			SpecificError: "string, unchanging so devs can use this in code",
			Message:       "detailed explanation of error",
			Source: errModels.Source{
				Field: "reference to field like some.field or something",
				Param: "query param causing issue"},
		}

		individualError := IndividualErrorBuilder(err, message, sourceField, sourceParam)

		So(individualError, ShouldResemble, individualErrorExample)

	})
}

func TestBuildingErrorStructure(t *testing.T) {
	Convey("An error structure is created from a list of errors", t, func() {

		listOfErrors := []errModels.IndividualError{
			{
				SpecificError: "string, unchanging so devs can use this in code",
				Message:       "detailed explanation of error",
				Source: errModels.Source{
					Field: "reference to field like some.field or something",
					Param: "query param causing issue"},
			},
		}

		errorResponseBodyExample := errModels.ErrorStructure{
			Errors: []errModels.IndividualError{
				{
					SpecificError: "string, unchanging so devs can use this in code",
					Message:       "detailed explanation of error",
					Source: errModels.Source{
						Field: "reference to field like some.field or something",
						Param: "query param causing issue"}},
			},
		}

		errorResponseBody := ErrorResponseBodyBuilder(listOfErrors)

		So(errorResponseBody, ShouldResemble, errorResponseBodyExample)
	})
}

func TestIdentifyInternalErrors(t *testing.T) {
	Convey("True is returned if an internal error is identified", t, func() {
		authError := errors.New("RequestError: send request failed")
		isInternalError := IdentifyInternalError(authError)

		So(isInternalError, ShouldBeTrue)
	})

	Convey("False is returned if an internal error is not identified", t, func() {

		authError := errors.New("NotAuthorizedException: Incorrect username or password.")
		isInternalError := IdentifyInternalError(authError)

		So(isInternalError, ShouldBeFalse)
	})
}
