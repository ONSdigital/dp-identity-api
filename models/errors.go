package models

import (
	"context"
	"errors"
	"github.com/ONSdigital/log.go/log"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

type ErrorResponse struct {
	Errors []BaseError `json:"errors"`
	Status int         `json:"-"`
}

type BaseError interface {
	Error() string
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
	log.Event(ctx, code+": "+description, log.ERROR)
	return &Error{
		Cause:       errors.New(code),
		Code:        code,
		Description: description,
	}
}

type CognitoError struct {
	Cause       error  `json:"-"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

func (e *CognitoError) Error() string {
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return e.Code + ": " + e.Description
}

func NewCognitoError(ctx context.Context, err error, errContext string) CognitoError {
	log.Event(ctx, errContext+" - "+err.Error(), log.ERROR)
	cognitoErr := err.(awserr.Error)
	code := MapCognitoErrorToLocalError(cognitoErr)
	return CognitoError{
		Cause:       cognitoErr,
		Code:        code,
		Description: cognitoErr.Message(),
	}
}

func MapCognitoErrorToLocalError(cognitoErr awserr.Error) string {
	if val, ok := CognitoErrorMapping[cognitoErr.Code()]; ok {
		return val
	}
	return InternalError
}
