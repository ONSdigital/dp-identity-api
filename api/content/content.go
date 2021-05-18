package content

// message, field and param error constants
const PasswordErrorMessage = "failed to generate password"
const PasswordErrorField = "temp password"
const PasswordErrorParam = "error creating temp password"

const RequestErrorMessage = "api endpoint POST user returned an error reading request body"
const RequestErrorField = "request body"
const RequestErrorParam = "error in reading request body"

const UnmarshallingErrorMessage = "api endpoint POST user returned an error unmarshalling request body"
const UnmarshallingErrorField = "unmarshalling"
const UnmarshallingErrorParam = "error unmarshalling request body"

const ValidForenameErrorField = "validating forename"
const ValidForenameErrorParam = "error validating username"

const ValidSurnameErrorField = "validating surname"
const ValidSurnameErrorParam = "error validating surname"

const ValidEmailErrorField = "validating email"
const ValidEmailErrorParam = "error validating email"

const DuplicateEmailFound = "duplicate email address found"

const NewUserModelErrorMessage = "Failed to create new user model"
const NewUserModelErrorField = "create new user model"
const NewUserModelErrorParam = "error creating new user model"

const ListUsersErrorMessage = "Error in checking duplicate email address"
const ListUsersErrorField = "duplicate email address check"
const ListUsersErrorParam = "error checking duplicate email address"

const AdminCreateUserErrorMessage = "Failed to create new user in user pool"
const AdminCreateUserErrorField = "create new user pool user"
const AdminCreateUserErrorParam = "error creating new user pool user"

const MarshallingNewUserErrorMessage = "Failed to marshall json response"
const MarshallingNewUserErrorField = "marshalling"
const MarshallingNewUserErrorParam = "error marshalling new user response"

const HttpResponseErrorMessage = "Failed to write http response"
const HttpResponseErrorField = "response"
const HttpResponseErrorParam = "error writing response"

const RequiredParameterNotFoundMessage = "error in parsing api setup arguments - missing parameter"
const InternalErrorException = "InternalErrorException"
