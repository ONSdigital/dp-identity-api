package models

import (
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

// API error codes
const (
	BodyReadError              = "RequestBodyReadError"
	JSONMarshalError           = "JSONMarshalError"
	JSONUnmarshalError         = "JSONUnmarshalError"
	WriteResponseError         = "WriteResponseError"
	InvalidEmailError          = "InvalidEmail"
	InvalidTokenError          = "InvalidToken"
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
	ErrorUnmarshalFailedDescription        = "failed to unmarshal the request body"
	WriteResponseFailedDescription         = "failed to write http response"
	CastingAWSErrorFailedDescription       = "failed to cast error to AWS error"
	UnrecognisedCognitoResponseDescription = "unexpected response from cognito"
	BodyReadFailedDescription              = "endpoint returned an error reading the request body"
	InvalidPasswordDescription             = "the submitted password could not be validated"
	InvalidEmailDescription                = "the submitted email could not be validated"
	SignInFailedDescription                = "Incorrect username or password"
	SignInAttemptsExceededDescription      = "Password attempts exceeded"
)

// Mapping Cognito error codes to API error codes
var CognitoErrorMapping = map[string]string{
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
	request.ErrCodeSerialization:                                   InternalError,
	request.ErrCodeRead:                                            InternalError,
	request.ErrCodeResponseTimeout:                                 InternalError,
	request.ErrCodeInvalidPresignExpire:                            InternalError,
	request.CanceledErrorCode:                                      InternalError,
	request.ErrCodeRequestError:                                    InternalError,
}
