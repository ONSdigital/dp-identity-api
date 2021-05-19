package apierrors

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	errModels "github.com/ONSdigital/dp-identity-api/models"
	"github.com/ONSdigital/log.go/log"
)

var InvalidAuthorizationTokenError = errors.New("invalid token")
var MissingAuthorizationTokenMessage = "no Authorization token was provided"
var MalformedTokenMessage = "the provided token does not meet the required format"

var InvalidRefreshTokenError = errors.New("invalid refresh token")
var MissingRefreshTokenMessage = "no refresh token was provided"

var InvalidIDTokenError = errors.New("invalid ID token")
var MissingIDTokenMessage = "no ID token was provided"

var ErrInvalidUserName = errors.New("invalid username")
var InvalidUserNameMessage = "Unable to validate the username in the request"

var ErrInvalidPassword = errors.New("invalid password")
var InvalidPasswordMessage = "Unable to validate the password in the request"

var ErrInvalidForename = errors.New("invalid forename")
var InvalidForenameErrorMessage = "Unable to validate the user's forename in the request"

var ErrInvalidSurname = errors.New("invalid surname")
var InvalidSurnameErrorMessage = "Unable to validate the user's surname in the request"

var ErrInvalidEmail = errors.New("invalid email")
var InvalidErrorMessage = "Unable to validate the email in the request"

var ErrDuplicateEmail = errors.New("duplicate email")

func IndividualErrorBuilder(err error, message, sourceField, sourceParam string) (individualError errModels.IndividualError) {

	individualError = errModels.IndividualError{
		SpecificError: error.Error(err),
		Message:       message,
		Source: errModels.Source{
			Field: sourceField,
			Param: sourceParam},
	}

	return individualError
}

func ErrorResponseBodyBuilder(listOfErrors []errModels.IndividualError) (errorResponseBody errModels.ErrorStructure) {

	errorResponseBody = errModels.ErrorStructure{
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

func HandleUnexpectedError(ctx context.Context, w http.ResponseWriter, err error, message, sourceField, sourceParam string) {

	var errorList []errModels.IndividualError

	internalServerErrorBody := IndividualErrorBuilder(err, message, sourceField, sourceParam)
	errorList = append(errorList, internalServerErrorBody)
	errorResponseBody := ErrorResponseBodyBuilder(errorList)

	log.Event(ctx, message, log.ERROR, log.Error(err))
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
