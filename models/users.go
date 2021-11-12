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
	MaxStatusNotesLength    = 512
	SecondsInDay            = 86400
)

type UsersList struct {
	Count           int          `json:"count"`
	Users           []UserParams `json:"users"`
	PaginationToken string
}

// ListGroupsForUser output structure from cognitoidentityprovider.AdminListGroupsForUserOutput but changing the
// json output
type ListUserGroupType struct {
	CreationDate     *time.Time `type:"timestamp" json:"creation_date"`
	Description      *string    `type:"string" json:"description"`
	GroupName        *string    `min:"1" type:"string" json:"group_name"`
	LastModifiedDate *time.Time `type:"timestamp" json:"last_modified_date"`
	Precedence       *int64     `type:"integer" json:"precedence"`
	RoleArn          *string    `min:"20" type:"string" json:"role_arn"`
	UserPoolId       *string    `min:"1" type:"string" json:"user_pool_id"`
}

// List groups for user output structure from cognitoidentityprovider.AdminListGroupsForUserOutput
// with count of total groups returned
type ListUserGroups struct {
	Groups    []*ListUserGroupType `json:"groups"`
	NextToken *string              `json:"next_token"`
	Count     int                  `json:"count"`
}

//BuildListUserRequest generates a ListUsersInput object for Cognito
func (p UsersList) BuildListUserRequest(filterString string, requiredAttribute string, limit int64, paginationToken *string, userPoolId *string) *cognitoidentityprovider.ListUsersInput {
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
	if paginationToken != nil {
		requestInput.PaginationToken = paginationToken
	}

	return requestInput
}

//MapCognitoUsers maps the users from the cognito response into the UsersList Users attribute and sets the Count attribute
func (p *UsersList) MapCognitoUsers(cognitoResults *[]*cognitoidentityprovider.UserType) {
	p.Users = []UserParams{}
	for _, user := range *cognitoResults {
		p.Users = append(p.Users, UserParams{}.MapCognitoDetails(user))
	}
	p.Count = len(p.Users)
}

//SetUsers sets the UsersList Users attribute and sets the Count attribute
func (p *UsersList) SetUsers(usersList *[]UserParams) {
	p.Users = *usersList
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

//Model for the User
type UserParams struct {
	Forename    string   `json:"forename"`
	Lastname    string   `json:"lastname"`
	Email       string   `json:"email"`
	Password    string   `json:"-"`
	Groups      []string `json:"groups"`
	Status      string   `json:"status"`
	Active      bool     `json:"active"`
	ID          string   `json:"id"`
	StatusNotes string   `json:"status_notes"`
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
func (p UserParams) ValidateRegistration(ctx context.Context, allowedDomains []string) []error {
	var validationErrs []error
	if p.Forename == "" {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidForenameError, InvalidForenameErrorDescription))
	}

	if p.Lastname == "" {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidSurnameError, InvalidSurnameErrorDescription))
	}

	if !validation.IsAllowedEmailDomain(p.Email, allowedDomains) {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidEmailError, InvalidEmailDescription))
	}
	return validationErrs
}

//ValidateUpdate validates the required fields for user update, returning validation errors for any failures
func (p UserParams) ValidateUpdate(ctx context.Context) []error {
	var validationErrs []error
	if p.Forename == "" {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidForenameError, InvalidForenameErrorDescription))
	}

	if p.Lastname == "" {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidSurnameError, InvalidSurnameErrorDescription))
	}

	if len(p.StatusNotes) > MaxStatusNotesLength {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidStatusNotesError, TooLongStatusNotesDescription))
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

//BuildUpdateUserRequest generates a AdminUpdateUserAttributesInput for Cognito
func (p UserParams) BuildUpdateUserRequest(userPoolId string) *cognitoidentityprovider.AdminUpdateUserAttributesInput {
	var (
		forenameAttrName, surnameAttrName, statusNotesAttrName string = "given_name", "family_name", "custom:status_notes"
	)

	return &cognitoidentityprovider.AdminUpdateUserAttributesInput{
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
				Name:  &statusNotesAttrName,
				Value: &p.StatusNotes,
			},
		},
		UserPoolId: &userPoolId,
		Username:   &p.ID,
	}
}

//BuildEnableUserRequest generates a AdminEnableUserInput for Cognito
func (p UserParams) BuildEnableUserRequest(userPoolId string) *cognitoidentityprovider.AdminEnableUserInput {
	return &cognitoidentityprovider.AdminEnableUserInput{
		UserPoolId: &userPoolId,
		Username:   &p.ID,
	}
}

//BuildDisableUserRequest generates a AdminDisableUserInput for Cognito
func (p UserParams) BuildDisableUserRequest(userPoolId string) *cognitoidentityprovider.AdminDisableUserInput {
	return &cognitoidentityprovider.AdminDisableUserInput{
		UserPoolId: &userPoolId,
		Username:   &p.ID,
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

//BuildAdminGetUserRequest generates a AdminGetUserInput for Cognito
func (p UserParams) BuildAdminGetUserRequest(userPoolId string) *cognitoidentityprovider.AdminGetUserInput {
	return &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: &userPoolId,
		Username:   &p.ID,
	}
}

//MapCognitoDetails maps the details from the Cognito ListUser User model to the UserParams model
func (p UserParams) MapCognitoDetails(userDetails *cognitoidentityprovider.UserType) UserParams {
	var forename, surname, email, statusNotes string
	for _, attr := range userDetails.Attributes {
		if *attr.Name == "given_name" {
			forename = *attr.Value
		} else if *attr.Name == "family_name" {
			surname = *attr.Value
		} else if *attr.Name == "email" {
			email = *attr.Value
		} else if *attr.Name == "custom:status_notes" {
			statusNotes = *attr.Value
		}
	}
	return UserParams{
		Forename:    forename,
		Lastname:    surname,
		Email:       email,
		Groups:      []string{},
		Status:      *userDetails.UserStatus,
		ID:          *userDetails.Username,
		StatusNotes: statusNotes,
		Active:      *userDetails.Enabled,
	}
}

//MapCognitoGetResponse maps the details from the Cognito GetUser User model to the UserParams model
func (p *UserParams) MapCognitoGetResponse(userDetails *cognitoidentityprovider.AdminGetUserOutput) {
	for _, attr := range userDetails.UserAttributes {
		if *attr.Name == "given_name" {
			p.Forename = *attr.Value
		} else if *attr.Name == "family_name" {
			p.Lastname = *attr.Value
		} else if *attr.Name == "email" {
			p.Email = *attr.Value
		} else if *attr.Name == "custom:status_notes" {
			p.StatusNotes = *attr.Value
		}
	}
	p.Status = *userDetails.UserStatus
	p.Groups = []string{}
	p.Active = *userDetails.Enabled
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
func (p *UserSignIn) BuildSuccessfulJsonResponse(ctx context.Context, result *cognitoidentityprovider.InitiateAuthOutput, refreshTokenTTL int) ([]byte, error) {
	if result.AuthenticationResult != nil {
		tokenDuration := time.Duration(*result.AuthenticationResult.ExpiresIn)
		expirationTime := time.Now().UTC().Add(time.Second * tokenDuration).String()
		refreshTokenDuration := time.Duration(SecondsInDay * refreshTokenTTL)
		refreshTokenExpirationTime := time.Now().UTC().Add(time.Second * refreshTokenDuration).String()

		postBody := map[string]interface{}{"expirationTime": expirationTime, "refreshTokenExpirationTime": refreshTokenExpirationTime}

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
	ChangeType        string `json:"type"`
	Session           string `json:"session"`
	Email             string `json:"email"`
	NewPassword       string `json:"password"`
	VerificationToken string `json:"verification_token"`
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
func (p ChangePassword) BuildAuthChallengeSuccessfulJsonResponse(ctx context.Context, result *cognitoidentityprovider.RespondToAuthChallengeOutput, refreshTokenTTL int) ([]byte, error) {
	if result.AuthenticationResult != nil {
		tokenDuration := time.Duration(*result.AuthenticationResult.ExpiresIn)
		expirationTime := time.Now().UTC().Add(time.Second * tokenDuration).String()
		refreshTokenDuration := time.Duration(SecondsInDay * refreshTokenTTL)
		refreshTokenExpirationTime := time.Now().UTC().Add(time.Second * refreshTokenDuration).String()

		postBody := map[string]interface{}{"expirationTime": expirationTime, "refreshTokenExpirationTime": refreshTokenExpirationTime}

		jsonResponse, err := json.Marshal(postBody)
		if err != nil {
			return nil, NewError(ctx, err, JSONMarshalError, ErrorMarshalFailedDescription)
		}
		return jsonResponse, nil
	} else {
		return nil, NewValidationError(ctx, InternalError, UnrecognisedCognitoResponseDescription)
	}
}

func (p ChangePassword) ValidateForgottenPasswordRequest(ctx context.Context) []error {
	var validationErrs []error
	if !validation.IsPasswordValid(p.NewPassword) {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidPasswordError, InvalidPasswordDescription))
	}
	// 'Email' in the forgotten password request is actually the user id, so we are only checking for presence rather than format
	if p.Email == "" {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidUserIdError, MissingUserIdErrorDescription))
	}
	if p.VerificationToken == "" {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidTokenError, InvalidTokenDescription))
	}
	return validationErrs
}

func (p ChangePassword) BuildConfirmForgotPasswordRequest(clientSecret string, clientId string) *cognitoidentityprovider.ConfirmForgotPasswordInput {
	secretHash := utilities.ComputeSecretHash(clientSecret, p.Email, clientId)

	return &cognitoidentityprovider.ConfirmForgotPasswordInput{
		ClientId:         &clientId,
		Username:         &p.Email,
		ConfirmationCode: &p.VerificationToken,
		SecretHash:       &secretHash,
		Password:         &p.NewPassword,
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

// BuildListUserGroupsRequest build the require input for cognito query to obtain the groups for given user
func (p UserParams) BuildListUserGroupsRequest(userPoolId string, nextToken string) *cognitoidentityprovider.AdminListGroupsForUserInput {

	if nextToken != "" {
		return &cognitoidentityprovider.AdminListGroupsForUserInput{
			UserPoolId: &userPoolId,
			Username:   &p.ID,
			NextToken:  &nextToken,
		}
	}

	return &cognitoidentityprovider.AdminListGroupsForUserInput{
		UserPoolId: &userPoolId,
		Username:   &p.ID}

}

//BuildListUserGroupsSuccessfulJsonResponse
// formats the output to comply with current standards and to json , adds the count of groups returned and
func (p *ListUserGroups) BuildListUserGroupsSuccessfulJsonResponse(ctx context.Context, result *cognitoidentityprovider.AdminListGroupsForUserOutput) ([]byte, error) {

	if result == nil {
		return nil, NewValidationError(ctx, InternalError, UnrecognisedCognitoResponseDescription)
	}

	for _, tmpGroup := range result.Groups {

		newGroup := ListUserGroupType{
			CreationDate:     tmpGroup.CreationDate,
			Description:      tmpGroup.Description,
			GroupName:        tmpGroup.GroupName,
			LastModifiedDate: tmpGroup.LastModifiedDate,
			Precedence:       tmpGroup.Precedence,
			RoleArn:          tmpGroup.RoleArn,
			UserPoolId:       tmpGroup.UserPoolId,
		}

		p.Groups = append(p.Groups, &newGroup)
	}

	p.NextToken = result.NextToken
	p.Count = 0
	if p.Groups != nil {
		p.Count = len(result.Groups)
	}

	jsonResponse, err := json.Marshal(p)
	if err != nil {
		return nil, NewError(ctx, err, JSONMarshalError, ErrorMarshalFailedDescription)
	}
	return jsonResponse, nil
}
