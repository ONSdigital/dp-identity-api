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

// ListUserGroupType output structure from cognitoidentityprovider.AdminListGroupsForUserOutput but changing the
// json output
type ListUserGroupType struct {
	CreationDate     *time.Time `type:"timestamp" json:"creation_date"`
	Name             *string    `type:"string" json:"name"`
	ID               *string    `min:"1" type:"string" json:"id"`
	LastModifiedDate *time.Time `type:"timestamp" json:"last_modified_date"`
	Precedence       *int64     `type:"integer" json:"precedence"`
	RoleArn          *string    `min:"20" type:"string" json:"role_arn"`
	UserPoolId       *string    `min:"1" type:"string" json:"user_pool_id"`
}

// ListUserGroups list of groups for user output structure from cognitoidentityprovider.AdminListGroupsForUserOutput
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
	Name                string   `json:"name"`
	GivenName           string   `json:"given_name"`
	FamilyName          string   `json:"family_name"`
	MiddleName          string   `json:"middle_name"`
	Nickname            string   `json:"nickname"`
	PreferredUsername   string   `json:"preferred_username"`
	Profile             string   `json:"profile"`
	Picture             string   `json:"picture"`
	Website             string   `json:"website"`
	Email               string   `json:"email"`
	EmailVerified       string   `json:"email_verified"`
	Gender              string   `json:"gender"`
	Birthdate           string   `json:"birthdate"`
	ZoneInfo            string   `json:"zoneinfo"`
	Locale              string   `json:"locale"`
	PhoneNumber         string   `json:"phone_number"`
	PhoneNumberVerified string   `json:"phone_number_verified"`
	Address             string   `json:"address"`
	UpdatedAt           string   `json:"updated_at"`
	CognitoMFAEnabled   string   `json:"cognito:mfa_enabled"`
	Username            string   `json:"cognito:username"`
	Password            string   `json:"-"`
	Groups              []string `json:"groups"`
	StatusNotes         string   `json:"status_notes"`
	//ID                  string   `json:"id"`
	Active bool   `json:"active"`
	Status string `json:"status"`
	//Lastname            string   `json:"lastname"`
	//Forename            string   `json:"forename"`
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
	if p.Name == "" {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidForenameError, InvalidForenameErrorDescription))
	}

	if !validation.IsAllowedEmailDomain(p.Email, allowedDomains) {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidEmailError, InvalidEmailDescription))
	}
	return validationErrs
}

//ValidateUpdate validates the required fields for user update, returning validation errors for any failures
func (p UserParams) ValidateUpdate(ctx context.Context) []error {
	var validationErrs []error
	if p.Name == "" {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidForenameError, InvalidForenameErrorDescription))
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
		deliveryMethod, nameAttrName, givenNameAttrName, familyNameAttrName, middleNameAttrName,
		nicknameAttrName, preferredUsernameAttrName, profileAttrName, pictureAttrName,
		websiteAttrName, emailAttrName, emailVerifiedAttrName, genderAttrName, birthdateAttrName,
		zoneInfoAttrName, localeAttrName, phoneNumberAttrName, phoneNumberVerifiedAttrName,
		addressAttrName, updatedAtAttrName string = "EMAIL", "name", "given_name", "family_name",
			"middle_name", "nickname", "preferred_username", "profile", "picture", "website",
			"email", "emailVerified", "gender", "birthdate", "zoneInfo", "locale", "phoneNumber",
			"phoneNumberVerified", "address", "updatedAt"
	)
	//var (
	//	deliveryMethod, nameAttrName, givenNameAttrName, emailAttrName, emailVerifiedAttrName, emailVerifiedValue, usernameAttrName string = "EMAIL", "name", "family_name", "email", "email_verified", "true", "username"
	//)

	return &cognitoidentityprovider.AdminCreateUserInput{
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  &nameAttrName,
				Value: &p.Name,
			},
			{
				Name:  &givenNameAttrName,
				Value: &p.GivenName,
			},
			{
				Name:  &familyNameAttrName,
				Value: &p.FamilyName,
			},
			{
				Name:  &middleNameAttrName,
				Value: &p.MiddleName,
			},
			{
				Name:  &nicknameAttrName,
				Value: &p.Nickname,
			},
			{
				Name:  &preferredUsernameAttrName,
				Value: &p.PreferredUsername,
			},
			{
				Name:  &profileAttrName,
				Value: &p.Profile,
			},
			{
				Name:  &pictureAttrName,
				Value: &p.Picture,
			},
			{
				Name:  &websiteAttrName,
				Value: &p.Website,
			},
			{
				Name:  &emailAttrName,
				Value: &p.Email,
			},
			{
				Name:  &emailVerifiedAttrName,
				Value: &p.EmailVerified,
			},
			{
				Name:  &genderAttrName,
				Value: &p.Gender,
			},
			{
				Name:  &birthdateAttrName,
				Value: &p.Birthdate,
			},
			{
				Name:  &zoneInfoAttrName,
				Value: &p.ZoneInfo,
			},
			{
				Name:  &localeAttrName,
				Value: &p.Locale,
			},
			{
				Name:  &phoneNumberAttrName,
				Value: &p.PhoneNumber,
			},
			{
				Name:  &phoneNumberVerifiedAttrName,
				Value: &p.PhoneNumberVerified,
			},
			{
				Name:  &addressAttrName,
				Value: &p.Address,
			},
			{
				Name:  &updatedAtAttrName,
				Value: &p.UpdatedAt,
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
		nameAttrName, givenNameAttrName, familyNameAttrName, middleNameAttrName,
		nicknameAttrName, preferredUsernameAttrName, profileAttrName, pictureAttrName,
		websiteAttrName, genderAttrName, birthdateAttrName,
		zoneInfoAttrName, localeAttrName, phoneNumberAttrName, phoneNumberVerifiedAttrName,
		addressAttrName string = "name", "given_name", "family_name", "middle_name", "nickname",
			"preferred_username", "profile", "picture", "website", "gender", "birthdate", "zoneInfo",
			"locale", "phoneNumber", "phoneNumberVerified", "address"
	)

	return &cognitoidentityprovider.AdminUpdateUserAttributesInput{
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  &nameAttrName,
				Value: &p.Name,
			},
			{
				Name:  &givenNameAttrName,
				Value: &p.GivenName,
			},
			{
				Name:  &familyNameAttrName,
				Value: &p.FamilyName,
			},
			{
				Name:  &middleNameAttrName,
				Value: &p.MiddleName,
			},
			{
				Name:  &nicknameAttrName,
				Value: &p.Nickname,
			},
			{
				Name:  &preferredUsernameAttrName,
				Value: &p.PreferredUsername,
			},
			{
				Name:  &profileAttrName,
				Value: &p.Profile,
			},
			{
				Name:  &pictureAttrName,
				Value: &p.Picture,
			},
			{
				Name:  &websiteAttrName,
				Value: &p.Website,
			},
			{
				Name:  &genderAttrName,
				Value: &p.Gender,
			},
			{
				Name:  &birthdateAttrName,
				Value: &p.Birthdate,
			},
			{
				Name:  &zoneInfoAttrName,
				Value: &p.ZoneInfo,
			},
			{
				Name:  &localeAttrName,
				Value: &p.Locale,
			},
			{
				Name:  &phoneNumberAttrName,
				Value: &p.PhoneNumber,
			},
			{
				Name:  &phoneNumberVerifiedAttrName,
				Value: &p.PhoneNumberVerified,
			},
			{
				Name:  &addressAttrName,
				Value: &p.Address,
			},
		},
		UserPoolId: &userPoolId,
		Username:   &p.Username,
	}
}

//BuildEnableUserRequest generates a AdminEnableUserInput for Cognito
func (p UserParams) BuildEnableUserRequest(userPoolId string) *cognitoidentityprovider.AdminEnableUserInput {
	return &cognitoidentityprovider.AdminEnableUserInput{
		UserPoolId: &userPoolId,
		Username:   &p.Username,
	}
}

//BuildDisableUserRequest generates a AdminDisableUserInput for Cognito
func (p UserParams) BuildDisableUserRequest(userPoolId string) *cognitoidentityprovider.AdminDisableUserInput {
	return &cognitoidentityprovider.AdminDisableUserInput{
		UserPoolId: &userPoolId,
		Username:   &p.Username,
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
		Username:   &p.Username,
	}
}

//MapCognitoDetails maps the details from the Cognito ListUser User model to the UserParams model
func (p UserParams) MapCognitoDetails(userDetails *cognitoidentityprovider.UserType) UserParams {
	var name, given_name, family_name, middle_name, nickname, preferred_username, profile,
		picture, website, email, email_verified, gender, birthdate, zoneinfo, locale, phone_number,
		phone_number_verified, address, updated_at string
	for _, attr := range userDetails.Attributes {
		if *attr.Name == "name" {
			name = *attr.Value
		} else if *attr.Name == "given_name" {
			given_name = *attr.Value
		} else if *attr.Name == "family_name" {
			family_name = *attr.Value
		} else if *attr.Name == "middle_name" {
			middle_name = *attr.Value
		} else if *attr.Name == "nickname" {
			nickname = *attr.Value
		} else if *attr.Name == "preferred_username" {
			preferred_username = *attr.Value
		} else if *attr.Name == "profile" {
			profile = *attr.Value
		} else if *attr.Name == "picture" {
			picture = *attr.Value
		} else if *attr.Name == "website" {
			website = *attr.Value
		} else if *attr.Name == "email" {
			email = *attr.Value
		} else if *attr.Name == "email_verified" {
			email_verified = *attr.Value
		} else if *attr.Name == "gender" {
			gender = *attr.Value
		} else if *attr.Name == "birthdate" {
			birthdate = *attr.Value
		} else if *attr.Name == "zoneinfo" {
			zoneinfo = *attr.Value
		} else if *attr.Name == "locale" {
			locale = *attr.Value
		} else if *attr.Name == "phone_number" {
			phone_number = *attr.Value
		} else if *attr.Name == "phone_number_verified" {
			phone_number_verified = *attr.Value
		} else if *attr.Name == "address" {
			address = *attr.Value
		} else if *attr.Name == "updated_at" {
			updated_at = *attr.Value
		}
	}

	return UserParams{
		Name:                name,
		GivenName:           given_name,
		FamilyName:          family_name,
		MiddleName:          middle_name,
		Nickname:            nickname,
		PreferredUsername:   preferred_username,
		Profile:             profile,
		Picture:             picture,
		Website:             website,
		Email:               email,
		EmailVerified:       email_verified,
		Gender:              gender,
		Birthdate:           birthdate,
		ZoneInfo:            zoneinfo,
		Locale:              locale,
		PhoneNumber:         phone_number,
		PhoneNumberVerified: phone_number_verified,
		Address:             address,
		UpdatedAt:           updated_at,
		Username:            *userDetails.Username,
		Groups:              []string{},
		Active:              *userDetails.Enabled,
	}
}

//MapCognitoGetResponse maps the details from the Cognito GetUser User model to the UserParams model
func (p *UserParams) MapCognitoGetResponse(userDetails *cognitoidentityprovider.AdminGetUserOutput) {
	for _, attr := range userDetails.UserAttributes {
		if *attr.Name == "name" {
			p.Name = *attr.Value
		} else if *attr.Name == "email" {
			p.Email = *attr.Value
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
			Username:   &p.Username,
			NextToken:  &nextToken,
		}
	}

	return &cognitoidentityprovider.AdminListGroupsForUserInput{
		UserPoolId: &userPoolId,
		Username:   &p.Username}

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
			Name:             tmpGroup.Description,
			ID:               tmpGroup.GroupName,
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
