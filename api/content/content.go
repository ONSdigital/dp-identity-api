package content

// message, field and param error constants
const PasswordErrorMessage           = "failed to generate password"
const PasswordErrorField             = "temp password"
const PasswordErrorParam             = "error creating temp password"

const RequestErrorMessage            = "api endpoint POST user returned an error reading request body"
const RequestErrorField              = "request body"
const RequestErrorParam              = "error in reading request body"

const UnmarshallingErrorMessage      = "api endpoint POST user returned an error unmarshalling request body"
const UnmarshallingErrorField        = "unmarshalling"
const UnmarshallingErrorParam        = "error unmarshalling request body"

const ValidUserNameErrorField        = "validating username"
const ValidUserNameErrorParam        = "error validating username"

const ValidEmailErrorField           = "validating email"
const ValidEmailErrorParam           = "error validating email"

const NewUserModelErrorMessage       = "Failed to create new user model"
const NewUserModelErrorField         = "create new user model"
const NewUserModelErrorParam         = "error creating new user model"

const AdminCreateUserErrorMessage    = "Failed to create new user in user pool"
const AdminCreateUserErrorField      = "create new user pool user"
const AdminCreateUserErrorParam      = "error creating new user pool user"

const MarshallingNewUserErrorMessage = "Failed to marshall json response"
const MarshallingNewUserErrorField   = "marshalling"
const MarshallingNewUserErrorParam   = "error marshalling new user response"

const HttpResponseErrorMessage       = "Failed to write http response"
const HttpResponseErrorField         = "response"
const HttpResponseErrorParam         = "error writing response"

const UserPoolIdNotFoundMessage      = "userPoolId must not be an empty string"
const InternalErrorException         = "InternalErrorException"