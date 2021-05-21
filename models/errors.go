package models

import (
	"context"
	"errors"
	"github.com/ONSdigital/log.go/log"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

type ErrorList struct {
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

func NewError(cause error, code string, description string) BaseError {
	return &Error{
		Cause:       cause,
		Code:        code,
		Description: description,
	}
}

func NewValidationError(ctx context.Context, code string, description string) BaseError {
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
	code := mapCognitoErrorToLocalError(cognitoErr)
	return CognitoError{
		Cause:       err,
		Code:        code,
		Description: cognitoErr.Message(),
	}
}

func mapCognitoErrorToLocalError(cognitoErr awserr.Error) string {
	if val, ok := cognitoErrorMapping[cognitoErr.Code()]; ok {
		return val
	}
	return InternalError
}