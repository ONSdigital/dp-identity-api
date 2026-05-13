package models_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/smithy-go"

	"github.com/ONSdigital/dp-identity-api/v2/models"
	. "github.com/smartystreets/goconvey/convey"
)

const serverError = smithy.ErrorFault(1)

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

	Convey("logs additional data when a single logData is provided", t, func() {
		cause := errors.New("TestError")
		errorCode := "TestErrorCode"
		errorDescription := "description of the error"
		var destination bytes.Buffer

		log.SetDestination(&destination, &destination)
		defer log.SetDestination(os.Stdout, os.Stderr)

		err := models.NewError(ctx, cause, errorCode, errorDescription, log.Data{
			"request_id": "abc-123",
			"endpoint":   "/tokens",
		})

		So(err, ShouldNotBeNil)

		var payload map[string]interface{}
		unmarshalErr := json.Unmarshal([]byte(strings.TrimSpace(destination.String())), &payload)
		So(unmarshalErr, ShouldBeNil)

		data, ok := payload["data"].(map[string]interface{})
		So(ok, ShouldBeTrue)
		So(data["request_id"], ShouldEqual, "abc-123")
		So(data["endpoint"], ShouldEqual, "/tokens")
	})

	Convey("merges multiple logData maps and prioritises later values for duplicate keys", t, func() {
		cause := errors.New("TestError")
		errorCode := "TestErrorCode"
		errorDescription := "description of the error"
		var destination bytes.Buffer

		log.SetDestination(&destination, &destination)
		defer log.SetDestination(os.Stdout, os.Stderr)

		err := models.NewError(
			ctx,
			cause,
			errorCode,
			errorDescription,
			log.Data{"service": "identity-api", "request_id": "first"},
			log.Data{"request_id": "second", "attempt": 2.0},
		)

		So(err, ShouldNotBeNil)

		var payload map[string]interface{}
		unmarshalErr := json.Unmarshal([]byte(strings.TrimSpace(destination.String())), &payload)
		So(unmarshalErr, ShouldBeNil)

		data, ok := payload["data"].(map[string]interface{})
		So(ok, ShouldBeTrue)
		So(data["service"], ShouldEqual, "identity-api")
		So(data["request_id"], ShouldEqual, "second")
		So(data["attempt"], ShouldEqual, 2.0)
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
		awsErr := &smithy.GenericAPIError{
			Code:    awsErrCode,
			Message: awsErrMessage,
			Fault:   serverError,
		}
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
			awsErr := &smithy.GenericAPIError{
				Code:    cognitoCode,
				Message: awsErrMessage,
			}
			actualONSCode := models.MapCognitoErrorToLocalError(ctx, awsErr)

			So(actualONSCode, ShouldEqual, expectedONSCode)
		}
	})

	Convey("returns internal error code when unmapped Cognito code provided", t, func() {
		cognitoCode := "EverythingIsOK"
		awsErrMessage := "an error occurred"
		awsErr := &smithy.GenericAPIError{
			Code:    cognitoCode,
			Message: awsErrMessage,
		}
		actualONSCode := models.MapCognitoErrorToLocalError(ctx, awsErr)

		So(actualONSCode, ShouldEqual, models.InternalError)
	})
}
