package models

import "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

const (
	InvalidTokenError                      = "InvalidToken"
	MissingAuthorizationTokenDescription   = "no Authorization token was provided"
	MalformedAuthorizationTokenDescription = "the authorization token does not meet the required format"
	InternalError                          = "InternalServerError"
	NotFoundError                          = "NotFound"
	AlreadyExistsError                     = "AlreadyExists"
	DeliveryFailureError                   = "DeliveryFailure"
	InvalidCodeError                       = "InvalidCode"
	ExpiredCodeError                       = "ExpiredCode"
	InvalidFieldError                      = "InvalidField"
	InvalidPasswordError                   = "InvalidPassword"
	LimitExceededError                     = "LimitExceeded"
	NotAuthorisedError                     = "NotAuthorised"
	PasswordResetRequiredError             = "PasswordResetRequired"
	TooManyFailedAttemptsError             = "TooManyFailedAttempts"
	TooManyRequestsError                   = "TooManyRequests"
	UserNotConfirmedError                  = "UserNotConfirmed"
	UsernameExistsError                    = "UsernameExists"
)

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
