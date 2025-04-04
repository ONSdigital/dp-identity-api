package models

import (
	"context"
	"errors"
	"github.com/aws/smithy-go"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

// Error represents a custom error type with additional context and description.
type Error struct {
	Cause       error  `json:"-"`           // The underlying error, if available.
	Code        string `json:"code"`        // Error code representing the type of error.
	Description string `json:"description"` // Detailed description of the error.
}

// Error returns the error message string for the custom Error type.
func (e *Error) Error() string {
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return e.Code + ": " + e.Description
}

// NewError creates and logs a new Error with the provided context, cause, code, and description.
func NewError(ctx context.Context, cause error, code, description string) *Error {
	err := &Error{
		Cause:       cause,
		Code:        code,
		Description: description,
	}
	log.Error(ctx, description, err)
	return err
}

// NewValidationError creates a new Error specifically for validation errors with a code and description.
func NewValidationError(ctx context.Context, code, description string) *Error {
	err := &Error{
		Cause:       errors.New(code),
		Code:        code,
		Description: description,
	}
	log.Error(ctx, description, err, log.Data{"code": code})
	return err
}

// IsGroupExistsError checks if the given error is a Cognito GroupExistsException error.
func IsGroupExistsError(err error) bool {
	//var cognitoErr smithy.APIError
	//if errors.As(err, &cognitoErr) && cognitoErr.ErrorCode() == cognitoidentityprovider.ErrCodeGroupExistsException {
	//	return true
	//}
	//return false
	var groupExistsErr *types.GroupExistsException
	return errors.As(err, &groupExistsErr)
}

// NewCognitoError creates a new Error for errors returned from AWS Cognito, mapping it to a local error code.
func NewCognitoError(ctx context.Context, err error, errContext string) *Error {
	log.Error(ctx, errContext, err)
	var cognitoErr smithy.APIError
	if !errors.As(err, &cognitoErr) {
		log.Error(ctx, CastingAWSErrorFailedDescription, err)
		return &Error{
			Cause:       err,
			Code:        InternalError,
			Description: CastingAWSErrorFailedDescription,
		}
	}
	code := MapCognitoErrorToLocalError(ctx, cognitoErr)
	return &Error{
		Cause:       cognitoErr,
		Code:        code,
		Description: cognitoErr.ErrorMessage(),
	}
}

// MapCognitoErrorToLocalError maps an AWS Cognito error to a local error code.
func MapCognitoErrorToLocalError(ctx context.Context, cognitoErr smithy.APIError) string {
	if val, ok := CognitoErrorMapping[cognitoErr.ErrorCode()]; ok {
		return val
	}
	log.Error(ctx, "unknown Cognito error code received", errors.New("unknown Cognito error code received"))
	return InternalError
}
