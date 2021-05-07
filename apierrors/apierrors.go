package apierrors

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/log.go/log"
)

type ErrorStructure struct {
	Errors []IndividualError `json:"errors"`
}

type IndividualError struct {
	SpecificError string `json:"error"`
	Message       string `json:"message"`
	Source        Source `json:"source"`
}

type Source struct {
	Field string `json:"field"`
	Param string `json:"param"`
}

func IndividualErrorBuilder(err error, message, sourceField, sourceParam string) (individualError IndividualError) {

	individualError = IndividualError{
		SpecificError: error.Error(err),
		Message:       message,
		Source: Source{
			Field: sourceField,
			Param: sourceParam},
	}

	return individualError
}

func ErrorResponseBodyBuilder(listOfErrors []IndividualError) (errorResponseBody ErrorStructure) {

	errorResponseBody = ErrorStructure{
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

	var errorList []IndividualError

	internalServerErrorBody := IndividualErrorBuilder(err, message, sourceField, sourceParam)
	errorList = append(errorList, internalServerErrorBody)
	errorResponseBody := ErrorResponseBodyBuilder(errorList)

	log.Event(ctx, message, log.ERROR, log.Error(err))
	WriteErrorResponse(ctx, w, http.StatusInternalServerError, errorResponseBody)
	return
}
