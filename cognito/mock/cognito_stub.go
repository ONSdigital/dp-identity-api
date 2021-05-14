package mock

import (
	"errors"
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
	status := "FORCE_CHANGE_PASSWORD"
	name := "smileons"

	if (*input.Username == "smileons") { // 201 - created successfully
		user := &models.CreateUserOutput{
			UserOutput: &cognitoidentityprovider.AdminCreateUserOutput{
				User: &cognitoidentityprovider.UserType{
					Username:   &name,
					UserStatus: &status,
				},
			},
		}
		return user.UserOutput, nil
	} else if (*input.Username == "bob") { // 400 - already exists
		return nil, errors.New("UsernameExistsException: User account already exists")
	}
	return nil, errors.New("InternalErrorException") // 500 - internal exception error
}

func (m *CognitoIdentityProviderClientStub) InitiateAuth(input *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error) {

	accessToken := "accessToken"
	var expiration int64 = 123
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
			return nil, errors.New("NotAuthorizedException: Incorrect username or password.")
		} else {
			return nil, errors.New("NotAuthorizedException: Password attempts exceeded")
		}
	}

	if *input.AuthParameters["PASSWORD"] == "internalerrorException" {
		return nil, errors.New("InternalErrorException")
	}
	return nil, errors.New("InvalidParameterException")
}

func (m *CognitoIdentityProviderClientStub) GlobalSignOut(signOutInput *cognitoidentityprovider.GlobalSignOutInput) (*cognitoidentityprovider.GlobalSignOutOutput, error) {
	if *signOutInput.AccessToken == "InternalError" {
		return nil, errors.New("InternalErrorException: Something went wrong")
	}
	for _, session := range m.Sessions {
		if session.AccessToken == *signOutInput.AccessToken {
			return &cognitoidentityprovider.GlobalSignOutOutput{}, nil
		}
	}
	return nil, errors.New("NotAuthorizedException: Access Token has been revoked")
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
	// default - email already doesn't exist in user pool
	users := &models.ListUsersOutput{
		ListUsersOutput: &cognitoidentityprovider.ListUsersOutput{
			Users: []*cognitoidentityprovider.UserType{},
		},
	}
	return users.ListUsersOutput, nil
}
