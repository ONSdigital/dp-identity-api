package mock

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"strings"

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
		status, subjectAttrName, forenameAttrName, surnameAttrName, emailAttrName, username, subUUID, forename, surname, email string = "FORCE_CHANGE_PASSWORD", "sub", "name", "family_name", "email", "123e4567-e89b-12d3-a456-426614174000", "f0cf8dd9-755c-4caf-884d-b0c56e7d0704", "smileons", "bobbings", "email@ons.gov.uk"
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
				return initiateAuthOutput, nil
			} else if user.email != *input.AuthParameters["USERNAME"] {
				return nil, awserr.New(cognitoidentityprovider.ErrCodeNotAuthorizedException, "Incorrect username or password", nil)
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
	name := "bob"
	if strings.Contains(*input.Filter, "email@ext.ons.gov.uk") {
		users := &models.ListUsersOutput{
			ListUsersOutput: &cognitoidentityprovider.ListUsersOutput{
				Users: []*cognitoidentityprovider.UserType{
					{
						Username: &name,
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
