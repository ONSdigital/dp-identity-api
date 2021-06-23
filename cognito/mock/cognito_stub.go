package mock

import (
	"errors"
	"regexp"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
)

type CognitoIdentityProviderClientStub struct {
	cognitoidentityprovideriface.CognitoIdentityProviderAPI
	UserPools []string
	Users     []User
	Sessions  []Session
}

func (m *CognitoIdentityProviderClientStub) DescribeUserPool(poolInputData *cognitoidentityprovider.DescribeUserPoolInput) (*cognitoidentityprovider.DescribeUserPoolOutput, error) {
	for _, v := range m.UserPools {
		if v == *poolInputData.UserPoolId {
			return nil, nil
		}
	}
	return nil, errors.New("failed to load user pool data")
}

func (m *CognitoIdentityProviderClientStub) AdminCreateUser(input *cognitoidentityprovider.AdminCreateUserInput) (*cognitoidentityprovider.AdminCreateUserOutput, error) {
	var (
		status, subjectAttrName, forenameAttrName, surnameAttrName, emailAttrName, username, subUUID, forename, surname, email string = "FORCE_CHANGE_PASSWORD", "sub", "name", "family_name", "email", "123e4567-e89b-12d3-a456-426614174000", "f0cf8dd9-755c-4caf-884d-b0c56e7d0704", "smileons", "bobbings", "emailx@ons.gov.uk"
	)

	if *input.UserAttributes[0].Value == "smileons" { // 201 - created successfully
		user := &models.CreateUserOutput{
			UserOutput: &cognitoidentityprovider.AdminCreateUserOutput{
				User: &cognitoidentityprovider.UserType{
					Attributes: []*cognitoidentityprovider.AttributeType{
						{
							Name:  &subjectAttrName,
							Value: &subUUID,
						},
						{
							Name:  &forenameAttrName,
							Value: &forename,
						},
						{
							Name:  &surnameAttrName,
							Value: &surname,
						},
						{
							Name:  &emailAttrName,
							Value: &email,
						},
					},
					Username:   &username,
					UserStatus: &status,
				},
			},
		}
		return user.UserOutput, nil
	}
	return nil, awserr.New(cognitoidentityprovider.ErrCodeInternalErrorException, "Failed to create new user in user pool", nil) // 500 - internal exception error
}

func (m *CognitoIdentityProviderClientStub) InitiateAuth(input *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	var expiration int64 = 123

	if *input.AuthFlow == "USER_PASSWORD_AUTH" {
		// non-verified response - ChallengName = "NEW_PASSWORD_REQUIRED"
		var (
			challengeName, sessionID string = "NEW_PASSWORD_REQUIRED", "AYABeBBsY5be-this-is-a-test-session-id-string-123456789iuerhcfdisieo-end"
		)
		initiateAuthOutputChallenge := &cognitoidentityprovider.InitiateAuthOutput{
			AuthenticationResult: nil,
			ChallengeName:        &challengeName,
			Session:              &sessionID,
		}

		// verified response - ChallengName = ""
		accessToken := "accessToken"
		idToken := "idToken"
		refreshToken := "refreshToken"
		initiateAuthOutput := &cognitoidentityprovider.InitiateAuthOutput{
			AuthenticationResult: &cognitoidentityprovider.AuthenticationResultType{
				AccessToken:  &accessToken,
				ExpiresIn:    &expiration,
				IdToken:      &idToken,
				RefreshToken: &refreshToken,
			},
		}

		for _, user := range m.Users {
			if (user.email == *input.AuthParameters["USERNAME"]) && (user.password == *input.AuthParameters["PASSWORD"]) {
				// non-challenge response
				if user.Attributes != nil {
					return initiateAuthOutput, nil
				} else {
					return initiateAuthOutputChallenge, nil
				}
			} else if user.email != *input.AuthParameters["USERNAME"] {
				return nil, awserr.New(cognitoidentityprovider.ErrCodeNotAuthorizedException, "Incorrect username or password.", nil)
			} else {
				return nil, awserr.New(cognitoidentityprovider.ErrCodeNotAuthorizedException, "Password attempts exceeded", nil)
			}
		}

		if *input.AuthParameters["PASSWORD"] == "internalerrorException" {
			return nil, awserr.New(cognitoidentityprovider.ErrCodeInternalErrorException, "Something went wrong", nil)
		}
		return nil, awserr.New(cognitoidentityprovider.ErrCodeInvalidParameterException, "A parameter was invalid", nil)
	} else if *input.AuthFlow == "REFRESH_TOKEN_AUTH" {
		if *input.AuthParameters["REFRESH_TOKEN"] == "InternalError" {
			return nil, awserr.New(cognitoidentityprovider.ErrCodeInternalErrorException, "Something went wrong", nil)
		} else if *input.AuthParameters["REFRESH_TOKEN"] == "ExpiredToken" {
			return nil, awserr.New(cognitoidentityprovider.ErrCodeNotAuthorizedException, "Refresh Token has expired", nil)
		} else {
			accessToken := "llll.mmmm.nnnn"
			idToken := "zzzz.yyyy.xxxx"
			initiateAuthOutput := &cognitoidentityprovider.InitiateAuthOutput{
				AuthenticationResult: &cognitoidentityprovider.AuthenticationResultType{
					AccessToken: &accessToken,
					ExpiresIn:   &expiration,
					IdToken:     &idToken,
				},
			}
			return initiateAuthOutput, nil
		}
	} else {
		return nil, errors.New("InvalidParameterException: Unknown Auth Flow")
	}
}

func (m *CognitoIdentityProviderClientStub) GlobalSignOut(signOutInput *cognitoidentityprovider.GlobalSignOutInput) (*cognitoidentityprovider.GlobalSignOutOutput, error) {
	if *signOutInput.AccessToken == "InternalError" {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeInternalErrorException, "Something went wrong", nil)
	}
	for _, session := range m.Sessions {
		if session.AccessToken == *signOutInput.AccessToken {
			return &cognitoidentityprovider.GlobalSignOutOutput{}, nil
		}
	}
	return nil, awserr.New(cognitoidentityprovider.ErrCodeNotAuthorizedException, "Access Token has been revoked", nil)
}

func (m *CognitoIdentityProviderClientStub) AdminUserGlobalSignOut(adminUserGlobalSignOutInput *cognitoidentityprovider.AdminUserGlobalSignOutInput) (*cognitoidentityprovider.AdminUserGlobalSignOutOutput, error) {
	if *adminUserGlobalSignOutInput.Username == "internalservererror@ons.gov.uk" {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeInternalErrorException, "Something went wrong", nil)
	} else if *adminUserGlobalSignOutInput.Username == "clienterror@ons.gov.uk" {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeNotAuthorizedException, "Something went wrong", nil)
	}
	return &cognitoidentityprovider.AdminUserGlobalSignOutOutput{}, nil
}

func (m *CognitoIdentityProviderClientStub) ListUsers(input *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error) {
	var (
		attribute_name, attribute_value string = "email_verified", "true"
	)

	getEmailFromFilter, _ := regexp.Compile(`^email\s\=\s(\D+.*)$`)
	email := getEmailFromFilter.ReplaceAllString(*input.Filter, `$1`)

	var emailRegex = regexp.MustCompile(`^\"email(\d)?@(ext\.)?ons.gov.uk\"`)
	if emailRegex.MatchString(email) {
		users := &models.ListUsersOutput{
			ListUsersOutput: &cognitoidentityprovider.ListUsersOutput{
				Users: []*cognitoidentityprovider.UserType{
					{
						Attributes: []*cognitoidentityprovider.AttributeType{
							{
								Name:  &attribute_name,
								Value: &attribute_value,
							},
						},
						Username: &email,
					},
				},
			},
		}
		return users.ListUsersOutput, nil
	}
	// default - email doesn't exist in user pool
	users := &models.ListUsersOutput{
		ListUsersOutput: &cognitoidentityprovider.ListUsersOutput{
			Users: []*cognitoidentityprovider.UserType{},
		},
	}
	return users.ListUsersOutput, nil
}

func (m *CognitoIdentityProviderClientStub) RespondToAuthChallenge(input *cognitoidentityprovider.RespondToAuthChallengeInput) (*cognitoidentityprovider.RespondToAuthChallengeOutput, error) {
	var expiration int64 = 123

	if *input.ChallengeName == "NEW_PASSWORD_REQUIRED" {
		accessToken := "accessToken"
		idToken := "idToken"
		refreshToken := "refreshToken"
		challengeResponseOutput := &cognitoidentityprovider.RespondToAuthChallengeOutput{
			AuthenticationResult: &cognitoidentityprovider.AuthenticationResultType{
				AccessToken:  &accessToken,
				ExpiresIn:    &expiration,
				IdToken:      &idToken,
				RefreshToken: &refreshToken,
			},
		}

		if *input.ChallengeResponses["NEW_PASSWORD"] == "internalerrorException" {
			return nil, awserr.New(cognitoidentityprovider.ErrCodeInternalErrorException, "Something went wrong", nil)
		} else if *input.ChallengeResponses["NEW_PASSWORD"] == "invalidpassword" {
			return nil, awserr.New(cognitoidentityprovider.ErrCodeInvalidPasswordException, "password does not meet requirements", nil)
		}

		for _, user := range m.Users {
			if user.email == *input.ChallengeResponses["USERNAME"] {
				return challengeResponseOutput, nil
			}
		}
		return nil, awserr.New(cognitoidentityprovider.ErrCodeUserNotFoundException, "user not found", nil)
	} else {
		return nil, errors.New("InvalidParameterException: Unknown Auth Flow")
	}
}

func (m *CognitoIdentityProviderClientStub) ForgotPassword(input *cognitoidentityprovider.ForgotPasswordInput) (*cognitoidentityprovider.ForgotPasswordOutput, error) {
	forgotPasswordOutput := &cognitoidentityprovider.ForgotPasswordOutput{
		CodeDeliveryDetails: &cognitoidentityprovider.CodeDeliveryDetailsType{},
	}

	if *input.Username == "internal.error@ons.gov.uk" {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeInternalErrorException, "Something went wrong", nil)
	}
	if *input.Username == "too.many@ons.gov.uk" {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeTooManyRequestsException, "Slow down", nil)
	}

	for _, user := range m.Users {
		if user.email == *input.Username {
			return forgotPasswordOutput, nil
		}
	}
	return nil, awserr.New(cognitoidentityprovider.ErrCodeUserNotFoundException, "user not found", nil)
}
