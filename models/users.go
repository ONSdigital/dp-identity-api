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

type UserParams struct {
	Forename string `json:"forename"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func (p UserParams) GeneratePassword(ctx context.Context) (*string, error) {
	tempPassword, err := password.Generate(14, 1, 1, false, false)
	if err != nil {
		return nil, NewError(ctx, err, InternalError, PasswordGenerationErrorDescription)
	}
	return &tempPassword, nil
}

func (p UserParams) ValidateRegistration(ctx context.Context) []error {
	var validationErrs []error
	if p.Forename == "" {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidForenameError, InvalidForenameErrorDescription))
	}

	if p.Surname == "" {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidSurnameError, InvalidSurnameErrorDescription))
	}

	if !validation.ValidateONSEmail(p.Email) {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidEmailError, InvalidEmailDescription))
	}
	return validationErrs
}

func (p UserParams) CheckForDuplicateEmail(ctx context.Context, listUserResp *cognitoidentityprovider.ListUsersOutput) error {
	if len(listUserResp.Users) == 0 {
		return nil
	}
	return NewValidationError(ctx, InvalidEmailError, DuplicateEmailDescription)
}

func (p UserParams) BuildListUserRequest(filterString string, requiredAttribute string, limit int64, userPoolId *string) *cognitoidentityprovider.ListUsersInput {
	return &cognitoidentityprovider.ListUsersInput{
		AttributesToGet: []*string{
			&requiredAttribute,
		},
		Filter:     &filterString,
		Limit:      &limit,
		UserPoolId: userPoolId,
	}
}

func (p UserParams) BuildCreateUserRequest(userId string, userPoolId string) *cognitoidentityprovider.AdminCreateUserInput {
	var (
		deliveryMethod, forenameAttrName, surnameAttrName, emailAttrName, emailVerifiedAttrName, emailVerifiedValue string = "EMAIL", "name", "family_name", "email", "email_verified", "true"
	)

	return &cognitoidentityprovider.AdminCreateUserInput{
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  &forenameAttrName,
				Value: &p.Forename,
			},
			{
				Name:  &surnameAttrName,
				Value: &p.Surname,
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

func (p UserParams) BuildSuccessfulJsonResponse(ctx context.Context, createdUser *cognitoidentityprovider.AdminCreateUserOutput) ([]byte, error) {
	jsonResponse, err := json.Marshal(createdUser)
	if err != nil {
		return nil, NewError(ctx, err, JSONMarshalError, ErrorMarshalFailedDescription)
	}
	return jsonResponse, nil
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

func (p *PasswordReset) BuildSuccessfulJsonResponse(ctx context.Context, result *cognitoidentityprovider.ForgotPasswordOutput) error {
	if result.CodeDeliveryDetails == nil {
		return NewValidationError(ctx, InternalError, UnrecognisedCognitoResponseDescription)
	}
	return nil
}
