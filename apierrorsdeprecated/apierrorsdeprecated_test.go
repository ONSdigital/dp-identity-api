package apierrorsdeprecated

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBuildingIndividualErrors(t *testing.T) {

	Convey("The individual error conforms to the expected structure", t, func() {

		err := errors.New("SomeError")
		description := "detailed explanation of error"

		individualErrorExample := Error{
			Code:        "SomeError",
			Description: "detailed explanation of error",
		}

		Error := IndividualErrorBuilder(err, description)

		So(Error, ShouldResemble, individualErrorExample)

	})
}

func TestBuildingErrorStructure(t *testing.T) {
	Convey("An error structure is created from a list of errors", t, func() {

		listOfErrors := []Error{
			{
				Code:        "SomeError",
				Description: "detailed explanation of error",
			},
		}

		errorResponseBodyExample := ErrorList{
			Errors: []Error{
				{
					Code:        "SomeError",
					Description: "detailed explanation of error",
				},
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
