package apierrorsdeprecated

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/ONSdigital/log.go/log"
)

const (
	MissingTokenDescription              = "no Authorization token was provided"
	MalformedTokenDescription            = "the provided token does not meet the required format"
	InvalidUserNameDescription           = "Unable to validate the username in the request"
	InvalidPasswordDescription           = "Unable to validate the password in the request"
	InvalidForenameErrorDescription      = "Unable to validate the user's forename in the request"
	InvalidSurnameErrorDescription       = "Unable to validate the user's surname in the request"
	InvalidErrorDescription              = "Unable to validate the email in the request"
	PasswordErrorDescription             = "failed to generate password"
	RequestErrorDescription              = "api endpoint POST user returned an error reading request body"
	UnmarshallingErrorDescription        = "api endpoint POST user returned an error unmarshalling request body"
	DuplicateEmailFound                  = "duplicate email address found"
	NewUserModelErrorDescription         = "Failed to create new user model"
	ListUsersErrorDescription            = "Error in checking duplicate email address"
	AdminCreateUserErrorDescription      = "Failed to create new user in user pool"
	MarshallingNewUserErrorDescription   = "Failed to marshall json response"
	HttpResponseErrorDescription         = "Failed to write http response"
	RequiredParameterNotFoundDescription = "error in parsing api setup arguments - missing parameter"
	InternalErrorException               = "InternalErrorException"
	MissingRefreshTokenMessage           = "no refresh token was provided"
	TokenExpiredMessage                  = "refresh token has expired"
	MissingIDTokenMessage                = "no ID token was provided"
	MalformedIDTokenMessage              = "the ID token could not be parsed"
	InternalErrorMessage                 = "an internal error has occurred"
)

var InvalidTokenError = errors.New("invalid token")

var InvalidRefreshTokenError = errors.New("invalid refresh token")

var InvalidIDTokenError = errors.New("invalid ID token")

var ErrInvalidUserName = errors.New("invalid username")

var ErrInvalidPassword = errors.New("invalid password")

var ErrInvalidForename = errors.New("invalid forename")

var ErrInvalidSurname = errors.New("invalid surname")

var ErrInvalidEmail = errors.New("invalid email")

var ErrDuplicateEmail = errors.New("duplicate email")

func IndividualErrorBuilder(err error, description string) (apiErr Error) {

	apiErr = Error{
		Code:        error.Error(err),
		Description: description,
	}

	return apiErr
}

func ErrorResponseBodyBuilder(listOfErrors []Error) (errorResponseBody ErrorList) {

	errorResponseBody = ErrorList{
		Errors: listOfErrors,
	}

	return errorResponseBody
}

func WriteErrorResponse(ctx context.Context, w http.ResponseWriter, status int, errorResponseBody interface{}) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	jsonResponse, err := json.Marshal(errorResponseBody)
	if err != nil {
		log.Event(ctx, "failed to marshal the error", log.Error(err), log.ERROR)
		http.Error(w, "failed to marshal the error", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Event(ctx, "writing response failed", log.Error(err), log.ERROR)
		http.Error(w, "failed to write http response", http.StatusInternalServerError)
		return
	}
}

func HandleUnexpectedError(ctx context.Context, w http.ResponseWriter, err error, description string) {

	var errorList []Error

	internalServerErrorBody := IndividualErrorBuilder(err, description)
	errorList = append(errorList, internalServerErrorBody)
	errorResponseBody := ErrorResponseBodyBuilder(errorList)

	log.Event(ctx, description, log.ERROR, log.Error(err))
	WriteErrorResponse(ctx, w, http.StatusInternalServerError, errorResponseBody)
}

func IdentifyInternalError(authErr error) (isInternalError bool) {
	possibleInternalErrors := []string{
		"InternalErrorException",
		"SerializationError",
		"ReadError",
		"ResponseTimeout",
		"InvalidPresignExpireError",
		"RequestCanceled",
		"RequestError",
	}

	for _, internalError := range possibleInternalErrors {
		if strings.Contains(authErr.Error(), internalError) {
			return true
		}
	}
	return false
}
