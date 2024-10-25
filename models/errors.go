package models

import (
	"context"
	"errors"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

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

func NewError(ctx context.Context, cause error, code, description string) *Error {
	err := &Error{
		Cause:       cause,
		Code:        code,
		Description: description,
	}
	log.Error(ctx, description, err)
	return err
}

func NewValidationError(ctx context.Context, code, description string) *Error {
	err := &Error{
		Cause:       errors.New(code),
		Code:        code,
		Description: description,
	}

	log.Error(ctx, description, err, log.Data{"code": code})
	return err
}

// IsGroupExistsError validates if the err occurred because group already exists
func IsGroupExistsError(err error) bool {
	var cognitoErr awserr.Error
	if errors.As(err, &cognitoErr) && cognitoErr.Code() == cognitoidentityprovider.ErrCodeGroupExistsException {
		return true
	}
	return false
}

func NewCognitoError(ctx context.Context, err error, errContext string) *Error {
	log.Error(ctx, errContext, err)
	var cognitoErr awserr.Error
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
		Description: cognitoErr.Message(),
	}
}

func MapCognitoErrorToLocalError(ctx context.Context, cognitoErr awserr.Error) string {
	if val, ok := CognitoErrorMapping[cognitoErr.Code()]; ok {
		return val
	}
	log.Error(ctx, "unknown Cognito error code received", errors.New("unknown Cognito error code received"))
	return InternalError
}
