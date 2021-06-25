package models

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ONSdigital/dp-identity-api/utilities"
	"github.com/ONSdigital/dp-identity-api/validation"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/sethvargo/go-password/password"
)

const (
	NewPasswordRequiredType = "NewPasswordRequired"
	ForgottenPasswordType   = "ForgottenPassword"
)

type UsersList struct {
	Users []UserParams `json:"users"`
	Count int          `json:"count"`
}

//BuildListUserRequest generates a ListUsersInput object for Cognito
func (p UsersList) BuildListUserRequest(filterString string, requiredAttribute string, limit int64, userPoolId *string) *cognitoidentityprovider.ListUsersInput {
	requestInput := &cognitoidentityprovider.ListUsersInput{
		UserPoolId: userPoolId,
	}
	if requiredAttribute != "" {
		requestInput.AttributesToGet = []*string{
			&requiredAttribute,
		}
	}
	if filterString != "" {
		requestInput.Filter = &filterString
	}
	if limit != 0 {
		requestInput.Limit = &limit
	}

	return requestInput
}

//MapCognitoUsers maps the users from the cognito response into the UsersList Users attribute and sets the Count attribute
func (p *UsersList) MapCognitoUsers(cognitoResults *cognitoidentityprovider.ListUsersOutput) {
	var usersList []UserParams
	for _, user := range cognitoResults.Users {
		usersList = append(usersList, UserParams{}.MapCognitoDetails(user))
	}
	p.Users = usersList
	p.Count = len(p.Users)
}

//BuildSuccessfulJsonResponse builds the UsersList response json for client responses
func (p *UsersList) BuildSuccessfulJsonResponse(ctx context.Context) ([]byte, error) {
	jsonResponse, err := json.Marshal(p)
	if err != nil {
		return nil, NewError(ctx, err, JSONMarshalError, ErrorMarshalFailedDescription)
	}
	return jsonResponse, nil
}

type UserParams struct {
	Forename string   `json:"forename"`
	Lastname string   `json:"lastname"`
	Email    string   `json:"email"`
	Password string   `json:"-"`
	Groups   []string `json:"groups"`
	Status   string   `json:"status"`
	ID       string   `json:"id"`
}

//GeneratePassword creates a password for the user and assigns it to the struct
func (p *UserParams) GeneratePassword(ctx context.Context) error {
	tempPassword, err := password.Generate(14, 1, 1, false, false)
	if err != nil {
		return NewError(ctx, err, InternalError, PasswordGenerationErrorDescription)
	}
	p.Password = tempPassword
	return nil
}

//ValidateRegistration validates the required fields for user creation, returning validation errors for any failures
func (p UserParams) ValidateRegistration(ctx context.Context) []error {
	var validationErrs []error
	if p.Forename == "" {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidForenameError, InvalidForenameErrorDescription))
	}

	if p.Lastname == "" {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidSurnameError, InvalidSurnameErrorDescription))
	}

	if !validation.ValidateONSEmail(p.Email) {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidEmailError, InvalidEmailDescription))
	}
	return validationErrs
}

//CheckForDuplicateEmail checks the Cognito response for users already using the email address, returning a validation error if found
func (p UserParams) CheckForDuplicateEmail(ctx context.Context, listUserResp *cognitoidentityprovider.ListUsersOutput) error {
	if len(listUserResp.Users) == 0 {
		return nil
	}
	return NewValidationError(ctx, InvalidEmailError, DuplicateEmailDescription)
}

//BuildCreateUserRequest generates a AdminCreateUserInput for Cognito
func (p UserParams) BuildCreateUserRequest(userId string, userPoolId string) *cognitoidentityprovider.AdminCreateUserInput {
	var (
		deliveryMethod, forenameAttrName, surnameAttrName, emailAttrName, emailVerifiedAttrName, emailVerifiedValue string = "EMAIL", "given_name", "family_name", "email", "email_verified", "true"
	)

	return &cognitoidentityprovider.AdminCreateUserInput{
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  &forenameAttrName,
				Value: &p.Forename,
			},
			{
				Name:  &surnameAttrName,
				Value: &p.Lastname,
			},
			{
				Name:  &emailAttrName,
				Value: &p.Email,
			},
			{
				Name:  &emailVerifiedAttrName,
				Value: &emailVerifiedValue,
			},
		},
		DesiredDeliveryMediums: []*string{
			&deliveryMethod,
		},
		TemporaryPassword: &p.Password,
		UserPoolId:        &userPoolId,
		Username:          &userId,
	}
}

//BuildSuccessfulJsonResponse builds the UserParams response json for client responses
func (p UserParams) BuildSuccessfulJsonResponse(ctx context.Context) ([]byte, error) {
	jsonResponse, err := json.Marshal(p)
	if err != nil {
		return nil, NewError(ctx, err, JSONMarshalError, ErrorMarshalFailedDescription)
	}
	return jsonResponse, nil
}

//BuildCreateUserRequest generates a AdminCreateUserInput for Cognito
func (p UserParams) BuildAdminGetUserRequest(userPoolId string) *cognitoidentityprovider.AdminGetUserInput {
	return &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: &userPoolId,
		Username:   &p.ID,
	}
}

//MapCognitoDetails maps the details from the Cognito User model to the UserParams model
func (p UserParams) MapCognitoDetails(userDetails *cognitoidentityprovider.UserType) UserParams {
	var forename, surname, email string
	for _, attr := range userDetails.Attributes {
		if *attr.Name == "given_name" {
			forename = *attr.Value
		} else if *attr.Name == "family_name" {
			surname = *attr.Value
		} else if *attr.Name == "email" {
			email = *attr.Value
		}
	}
	return UserParams{
		Forename: forename,
		Lastname: surname,
		Email:    email,
		Groups:   []string{},
		Status:   *userDetails.UserStatus,
		ID:       *userDetails.Username,
	}
}

type CreateUserInput struct {
	UserInput *cognitoidentityprovider.AdminCreateUserInput
}
type CreateUserOutput struct {
	UserOutput *cognitoidentityprovider.AdminCreateUserOutput
}
type ListUsersInput struct {
	ListUsersInput *cognitoidentityprovider.ListUsersInput
}
type ListUsersOutput struct {
	ListUsersOutput *cognitoidentityprovider.ListUsersOutput
}

type UserSignIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//ValidateCredentials validates the required fields have been submitted and meet the basic structure requirements
func (p *UserSignIn) ValidateCredentials(ctx context.Context) *[]error {
	var validationErrors []error
	if !validation.IsPasswordValid(p.Password) {
		validationErrors = append(validationErrors, NewValidationError(ctx, InvalidPasswordError, InvalidPasswordDescription))
	}
	if !validation.IsEmailValid(p.Email) {
		validationErrors = append(validationErrors, NewValidationError(ctx, InvalidEmailError, InvalidEmailDescription))
	}
	if len(validationErrors) == 0 {
		return nil
	}
	return &validationErrors
}

//BuildCognitoRequest generates a InitiateAuthInput for Cognito
func (p *UserSignIn) BuildCognitoRequest(clientId string, clientSecret string, clientAuthFlow string) *cognitoidentityprovider.InitiateAuthInput {
	secretHash := utilities.ComputeSecretHash(clientSecret, p.Email, clientId)

	authParameters := map[string]*string{
		"USERNAME":    &p.Email,
		"PASSWORD":    &p.Password,
		"SECRET_HASH": &secretHash,
	}

	return &cognitoidentityprovider.InitiateAuthInput{
		AnalyticsMetadata: &cognitoidentityprovider.AnalyticsMetadataType{},
		AuthFlow:          &clientAuthFlow,
		AuthParameters:    authParameters,
		ClientId:          &clientId,
		ClientMetadata:    map[string]*string{},
		UserContextData:   &cognitoidentityprovider.UserContextDataType{},
	}
}

//BuildSuccessfulJsonResponse builds the UserSignIn response json for client responses
func (p *UserSignIn) BuildSuccessfulJsonResponse(ctx context.Context, result *cognitoidentityprovider.InitiateAuthOutput) ([]byte, error) {
	if result.AuthenticationResult != nil {
		tokenDuration := time.Duration(*result.AuthenticationResult.ExpiresIn)
		expirationTime := time.Now().UTC().Add(time.Second * tokenDuration).String()

		postBody := map[string]interface{}{"expirationTime": expirationTime}

		jsonResponse, err := json.Marshal(postBody)
		if err != nil {
			return nil, NewError(ctx, err, JSONMarshalError, ErrorMarshalFailedDescription)
		}
		return jsonResponse, nil
	} else if result.ChallengeName != nil && *result.ChallengeName == "NEW_PASSWORD_REQUIRED" {
		postBody := map[string]interface{}{
			"new_password_required": "true",
			"session":               *result.Session,
		}

		jsonResponse, err := json.Marshal(postBody)
		if err != nil {
			return nil, NewError(ctx, err, JSONMarshalError, ErrorMarshalFailedDescription)
		}
		return jsonResponse, nil
	} else {
		return nil, NewValidationError(ctx, InternalError, UnrecognisedCognitoResponseDescription)
	}
}

type ChangePassword struct {
	ChangeType  string `json:"type"`
	Session     string `json:"session"`
	Email       string `json:"email"`
	NewPassword string `json:"password"`
}

//ValidateNewPasswordRequiredRequest validates the required fields have been submitted and meet the basic structure requirements
func (p ChangePassword) ValidateNewPasswordRequiredRequest(ctx context.Context) []error {
	var validationErrs []error
	if !validation.IsPasswordValid(p.NewPassword) {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidPasswordError, InvalidPasswordDescription))
	}
	if !validation.IsEmailValid(p.Email) {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidEmailError, InvalidEmailDescription))
	}
	if p.Session == "" {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidChallengeSessionError, InvalidChallengeSessionDescription))
	}
	return validationErrs
}

//BuildAuthChallengeResponseRequest generates a RespondToAuthChallengeInput for Cognito
func (p ChangePassword) BuildAuthChallengeResponseRequest(clientSecret string, clientId string, challengeName string) *cognitoidentityprovider.RespondToAuthChallengeInput {
	secretHash := utilities.ComputeSecretHash(clientSecret, p.Email, clientId)

	challengeResponses := map[string]*string{
		"USERNAME":     &p.Email,
		"NEW_PASSWORD": &p.NewPassword,
		"SECRET_HASH":  &secretHash,
	}

	return &cognitoidentityprovider.RespondToAuthChallengeInput{
		ClientId:           &clientId,
		ChallengeName:      &challengeName,
		Session:            &p.Session,
		ChallengeResponses: challengeResponses,
	}
}

//BuildAuthChallengeSuccessfulJsonResponse builds the ChangePassword response json for client responses to NewPasswordRequired changes
func (p ChangePassword) BuildAuthChallengeSuccessfulJsonResponse(ctx context.Context, result *cognitoidentityprovider.RespondToAuthChallengeOutput) ([]byte, error) {
	if result.AuthenticationResult != nil {
		tokenDuration := time.Duration(*result.AuthenticationResult.ExpiresIn)
		expirationTime := time.Now().UTC().Add(time.Second * tokenDuration).String()

		postBody := map[string]interface{}{"expirationTime": expirationTime}

		jsonResponse, err := json.Marshal(postBody)
		if err != nil {
			return nil, NewError(ctx, err, JSONMarshalError, ErrorMarshalFailedDescription)
		}
		return jsonResponse, nil
	} else {
		return nil, NewValidationError(ctx, InternalError, UnrecognisedCognitoResponseDescription)
	}
}

type PasswordReset struct {
	Email string `json:"email"`
}

func (p *PasswordReset) Validate(ctx context.Context) error {
	if !validation.IsEmailValid(p.Email) {
		return NewValidationError(ctx, InvalidEmailError, InvalidEmailDescription)
	}
	return nil
}

func (p PasswordReset) BuildCognitoRequest(clientSecret string, clientId string) *cognitoidentityprovider.ForgotPasswordInput {
	secretHash := utilities.ComputeSecretHash(clientSecret, p.Email, clientId)
	return &cognitoidentityprovider.ForgotPasswordInput{
		ClientId:   &clientId,
		SecretHash: &secretHash,
		Username:   &p.Email,
	}
}
