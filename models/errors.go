package models

import (
	"context"
	"errors"
	"github.com/ONSdigital/log.go/log"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

type ErrorResponse struct {
	Errors []error `json:"errors"`
	Status int     `json:"-"`
}

type Error struct {
	Cause       error  `json:"-"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return e.Code + ": " + e.Description
}

func NewError(ctx context.Context, cause error, code string, description string) *Error {
	log.Event(ctx, description, log.Error(cause), log.ERROR)
	return &Error{
		Cause:       cause,
		Code:        code,
		Description: description,
	}
}

func NewValidationError(ctx context.Context, code string, description string) *Error {
	log.Event(ctx, description, log.ERROR, log.Data{"code": code})
	return &Error{
		Cause:       errors.New(code),
		Code:        code,
		Description: description,
	}
}

func NewCognitoError(ctx context.Context, err error, errContext string) *Error {
	log.Event(ctx, errContext, log.Error(err), log.ERROR)
	var cognitoErr awserr.Error
	if !errors.As(err, &cognitoErr) {
		log.Event(ctx, CastingAWSErrorFailedDescription, log.Error(err), log.ERROR)
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
		Description: cognitoErr.Message(),
	}
}

func MapCognitoErrorToLocalError(ctx context.Context, cognitoErr awserr.Error) string {
	if val, ok := CognitoErrorMapping[cognitoErr.Code()]; ok {
		return val
	}
	log.Event(ctx, "unknown Cognito error code received", log.ERROR)
	return InternalError
}
