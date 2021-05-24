package models

import "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

// API error codes
const (
	InvalidTokenError          = "InvalidToken"
	JSONMarshalError           = "JSONMarshalError"
	WriteResponseError         = "WriteResponseError"
	InternalError              = "InternalServerError"
	NotFoundError              = "NotFound"
	AlreadyExistsError         = "AlreadyExists"
	DeliveryFailureError       = "DeliveryFailure"
	InvalidCodeError           = "InvalidCode"
	ExpiredCodeError           = "ExpiredCode"
	InvalidFieldError          = "InvalidField"
	InvalidPasswordError       = "InvalidPassword"
	LimitExceededError         = "LimitExceeded"
	NotAuthorisedError         = "NotAuthorised"
	PasswordResetRequiredError = "PasswordResetRequired"
	TooManyFailedAttemptsError = "TooManyFailedAttempts"
	TooManyRequestsError       = "TooManyRequests"
	UserNotConfirmedError      = "UserNotConfirmed"
	UsernameExistsError        = "UsernameExists"
)

// API error descriptions
const (
	MissingAuthorizationTokenDescription   = "no Authorization token was provided"
	MalformedAuthorizationTokenDescription = "the authorization token does not meet the required format"
	ErrorMarshalFailedDescription          = "failed to marshal the error"
	WriteResponseFailedDescription         = "failed to write http response"
)

// Mapping Cognito error codes to API error codes
var cognitoErrorMapping = map[string]string{
	cognitoidentityprovider.ErrCodeInternalErrorException:          InternalError,
	cognitoidentityprovider.ErrCodeCodeDeliveryFailureException:    DeliveryFailureError,
	cognitoidentityprovider.ErrCodeCodeMismatchException:           InvalidCodeError,
	cognitoidentityprovider.ErrCodeConcurrentModificationException: InternalError,
	cognitoidentityprovider.ErrCodeExpiredCodeException:            ExpiredCodeError,
	cognitoidentityprovider.ErrCodeGroupExistsException:            AlreadyExistsError,
	cognitoidentityprovider.ErrCodeInvalidOAuthFlowException:       InternalError,
	cognitoidentityprovider.ErrCodeInvalidParameterException:       InvalidFieldError,
	cognitoidentityprovider.ErrCodeInvalidPasswordException:        InvalidPasswordError,
	cognitoidentityprovider.ErrCodeLimitExceededException:          LimitExceededError,
	cognitoidentityprovider.ErrCodeNotAuthorizedException:          NotAuthorisedError,
	cognitoidentityprovider.ErrCodePasswordResetRequiredException:  PasswordResetRequiredError,
	cognitoidentityprovider.ErrCodeResourceNotFoundException:       NotFoundError,
	cognitoidentityprovider.ErrCodeTooManyFailedAttemptsException:  TooManyFailedAttemptsError,
	cognitoidentityprovider.ErrCodeTooManyRequestsException:        TooManyRequestsError,
	cognitoidentityprovider.ErrCodeUserNotConfirmedException:       UserNotConfirmedError,
	cognitoidentityprovider.ErrCodeUserNotFoundException:           NotFoundError,
	cognitoidentityprovider.ErrCodeUsernameExistsException:         UsernameExistsError,
}
