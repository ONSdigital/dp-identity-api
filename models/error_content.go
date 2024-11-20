package models

import (
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

// API error codes
const (
	BodyReadError                = "RequestBodyReadError"
	JSONMarshalError             = "JSONMarshalError"
	JSONUnmarshalError           = "JSONUnmarshalError"
	WriteResponseError           = "WriteResponseError"
	InvalidUserIDError           = "InvalidUserId"
	InvalidGroupIDError          = "InvalidGroupID"
	InvalidForenameError         = "InvalidForename"
	InvalidSurnameError          = "InvalidSurname"
	InvalidStatusNotesError      = "InvalidStatusNotes"
	InvalidEmailError            = "InvalidEmail"
	InvalidTokenError            = "InvalidToken"
	InternalError                = "InternalServerError"
	NotFoundError                = "NotFound"
	UserNotFoundError            = "UserNotFound"
	GroupExistsError             = "GroupExists"
	GroupNotFoundError           = "GroupNotFound"
	DeliveryFailureError         = "DeliveryFailure"
	InvalidCodeError             = "InvalidCode"
	ExpiredCodeError             = "ExpiredCode"
	InvalidFieldError            = "InvalidField"
	InvalidPasswordError         = "InvalidPassword"
	LimitExceededError           = "LimitExceeded"
	NotAuthorisedError           = "NotAuthorised"
	PasswordResetRequiredError   = "PasswordResetRequired"
	TooManyFailedAttemptsError   = "TooManyFailedAttempts"
	TooManyRequestsError         = "TooManyRequests"
	UserNotConfirmedError        = "UserNotConfirmed"
	UsernameExistsError          = "UsernameExists"
	MissingConfigError           = "MissingConfig"
	UnknownRequestTypeError      = "UnknownRequestType"
	NotImplementedError          = "NotImplemented"
	InvalidChallengeSessionError = "InvalidChallengeSession"
	InvalidUserPoolError         = "InvalidUserPool"
	BodyCloseError               = "BodyCloseError"
	InvalidGroupName             = "InvalidGroupName"
	InvalidGroupPrecedence       = "InvalidGroupPrecedence"
	InvalidFilterQuery           = "InvalidFilterQuery"
	JWKSParseError               = "JWKSParseError"
)

// API error descriptions
const (
	MissingAuthorizationTokenDescription   = "no Authorization token was provided"
	MissingRefreshTokenDescription         = "no Refresh token was provided"
	MissingIDTokenDescription              = "no ID token was provided"         //nolint:gosec // not a hardcoded secret
	MalformedIDTokenDescription            = "the ID token could not be parsed" //nolint:gosec // not a hardcoded secret
	MalformedAuthorizationTokenDescription = "the authorization token does not meet the required format"
	ErrorMarshalFailedDescription          = "failed to marshal the error"
	ErrorUnmarshalFailedDescription        = "failed to unmarshal the request body"
	WriteResponseFailedDescription         = "failed to write http response"
	CastingAWSErrorFailedDescription       = "failed to cast error to AWS error"
	UnrecognisedCognitoResponseDescription = "unexpected response from cognito"
	BodyReadFailedDescription              = "endpoint returned an error reading the request body"
	InvalidPasswordDescription             = "the submitted password could not be validated"
	PasswordGenerationErrorDescription     = "failed to generate a valid password"
	MissingGroupIDErrorDescription         = "the group ID was missing"
	MissingUserIDErrorDescription          = "the user id was missing"
	InvalidForenameErrorDescription        = "the submitted user's forename could not be validated"
	InvalidSurnameErrorDescription         = "the submitted user's lastname could not be validated"
	InvalidEmailDescription                = "the submitted email could not be validated"
	DuplicateEmailDescription              = "account using email address found"
	SignInFailedDescription                = "Incorrect username or password."
	SignInAttemptsExceededDescription      = "Password attempts exceeded"
	MissingConfigDescription               = "required configuration setting is missing"
	UnknownPasswordChangeTypeDescription   = "unknown password change type received"
	NotImplementedDescription              = "this feature has not been implemented yet"
	InvalidChallengeSessionDescription     = "no valid auth challenge session was provided"
	InvalidTokenDescription                = "the submitted token could not be validated"
	TooLongStatusNotesDescription          = "the status notes are too long"
	InvalidUserPoolDescription             = "dummy data load being run against non local userpool"
	BodyClosedFailedDescription            = "the request body failed to close"
	MissingGroupName                       = "the group name was not found"
	MissingGroupPrecedence                 = "the group precedence was not found"
	GroupPrecedenceIncorrect               = "the group precedence needs to be a minumum of 10 and maximum of 100"
	IncorrectPatternInGroupName            = "a group name cannot start with 'role-' or 'ROLE-'"
	GroupAlreadyExistsDescription          = "a group with the name already exists"
	InvalidFilterQueryDescription          = "the submitted query could not be validated"
	InternalErrorDescription               = "Internal Server Error"
	JWKSParseErrorDescription              = "error encountered when parsing the json web key set (jwks)"
	JWKSUnsupportedKeyTypeDescription      = "unsupported key type. Must be rsa key"
	JWKSErrorDecodingDescription           = "error decoding json web key"
	JWKSExponentErrorDescription           = "unexpected exponent: unable to decode JWK"
	JWKSEmptyWebKeySetDescription          = "empty json web key set"
)

// CognitoErrorMapping mapping Cognito error codes to API error codes
var CognitoErrorMapping = map[string]string{
	cognitoidentityprovider.ErrCodeInternalErrorException:          InternalError,
	cognitoidentityprovider.ErrCodeCodeDeliveryFailureException:    DeliveryFailureError,
	cognitoidentityprovider.ErrCodeCodeMismatchException:           InvalidCodeError,
	cognitoidentityprovider.ErrCodeConcurrentModificationException: InternalError,
	cognitoidentityprovider.ErrCodeExpiredCodeException:            ExpiredCodeError,
	cognitoidentityprovider.ErrCodeGroupExistsException:            GroupExistsError,
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
	cognitoidentityprovider.ErrCodeUserNotFoundException:           UserNotFoundError,
	cognitoidentityprovider.ErrCodeUsernameExistsException:         UsernameExistsError,
	request.ErrCodeSerialization:                                   InternalError,
	request.ErrCodeRead:                                            InternalError,
	request.ErrCodeResponseTimeout:                                 InternalError,
	request.ErrCodeInvalidPresignExpire:                            InternalError,
	request.CanceledErrorCode:                                      InternalError,
	request.ErrCodeRequestError:                                    InternalError,
}
