package models_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/aws/aws-sdk-go/aws/awserr"
	. "github.com/smartystreets/goconvey/convey"
)

func TestError_Error(t *testing.T) {
	Convey("returns the cause Error value when a cause is set", t, func() {
		originalErr := errors.New("OriginalErrorCause")
		errorCode := "TestErrorCode"
		errorDescription := "description of the error"

		err := models.Error{
			Cause:       originalErr,
			Code:        errorCode,
			Description: errorDescription,
		}

		So(err.Error(), ShouldEqual, originalErr.Error())
	})

	Convey("returns the Code and Description when a cause is not set", t, func() {
		errorCode := "TestErrorCode"
		errorDescription := "description of the error"

		err := models.Error{
			Code:        errorCode,
			Description: errorDescription,
		}

		So(err.Error(), ShouldEqual, errorCode+": "+errorDescription)
	})
}

func TestNewError(t *testing.T) {
	var ctx = context.Background()

	Convey("successfully constructs an Error object", t, func() {
		cause := errors.New("TestError")
		errorCode := "TestErrorCode"
		errorDescription := "description of the error"

		err := models.NewError(ctx, cause, errorCode, errorDescription)

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, cause.Error())
		So(err.Code, ShouldEqual, errorCode)
		So(err.Description, ShouldEqual, errorDescription)
	})
}

func TestNewValidationError(t *testing.T) {
	var ctx = context.Background()

	Convey("successfully constructs an Error object", t, func() {
		errorCode := "TestErrorCode"
		errorDescription := "description of the error"

		err := models.NewValidationError(ctx, errorCode, errorDescription)

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, errorCode)
		So(err.Code, ShouldEqual, errorCode)
		So(err.Description, ShouldEqual, errorDescription)
	})
}

func TestNewCognitoError(t *testing.T) {
	var ctx = context.Background()

	Convey("successfully constructs a CognitoError object", t, func() {
		awsErrCode := "InternalErrorException"
		awsErrMessage := "Something strange happened"
		awsOrigErr := errors.New(awsErrCode)
		awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
		errorContext := "dp-identity-api calling AWS Cognito"

		err := models.NewCognitoError(ctx, awsErr, errorContext)

		expectedErrorCode := models.CognitoErrorMapping[awsErrCode]

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, awsErr.Error())
		So(err.Code, ShouldEqual, expectedErrorCode)
		So(err.Description, ShouldEqual, awsErrMessage)
	})

	Convey("constructs a internal error CognitoError object if cast to AWS error fails", t, func() {
		originalErr := errors.New("TestError")
		errorContext := "dp-identity-api calling AWS Cognito"

		err := models.NewCognitoError(ctx, originalErr, errorContext)

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, originalErr.Error())
		So(err.Code, ShouldEqual, models.InternalError)
		So(err.Description, ShouldEqual, models.CastingAWSErrorFailedDescription)
	})
}

func TestMapCognitoErrorToLocalError(t *testing.T) {
	var ctx = context.Background()

	Convey("correctly maps the Cognito error code to the ONS error code", t, func() {
		for cognitoCode, expectedONSCode := range models.CognitoErrorMapping {
			awsErrMessage := "an error occurred"
			awsOrigErr := errors.New(cognitoCode)
			awsErr := awserr.New(cognitoCode, awsErrMessage, awsOrigErr)
			actualONSCode := models.MapCognitoErrorToLocalError(ctx, awsErr)

			So(actualONSCode, ShouldEqual, expectedONSCode)
		}
	})

	Convey("returns internal error code when unmapped Cognito code provided", t, func() {
		cognitoCode := "EverythingIsOK"
		awsErrMessage := "an error occurred"
		awsOrigErr := errors.New(cognitoCode)
		awsErr := awserr.New(cognitoCode, awsErrMessage, awsOrigErr)
		actualONSCode := models.MapCognitoErrorToLocalError(ctx, awsErr)

		So(actualONSCode, ShouldEqual, models.InternalError)
	})
}
