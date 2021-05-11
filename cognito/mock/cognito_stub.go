package mock

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
)

type CognitoIdentityProviderClientStub struct {
	cognitoidentityprovideriface.CognitoIdentityProviderAPI
	UserPools []string
	Users     []User
}

func (m *CognitoIdentityProviderClientStub) DescribeUserPool(poolInputData *cognitoidentityprovider.DescribeUserPoolInput) (*cognitoidentityprovider.DescribeUserPoolOutput, error) {
	for _, v := range m.UserPools {
		if v == *poolInputData.UserPoolId {
			return nil, nil
		}
	}
	return nil, errors.New("Failed to load user pool data")
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
		} else if (user.email == *input.AuthParameters["USERNAME"]) && (user.password != *input.AuthParameters["PASSWORD"]) {
			fmt.Println("I came in here")
			return nil, errors.New("NotAuthorizedException")
		} else {
			return nil, errors.New("NotAuthorizedException")
		}
	}
	return nil, errors.New("InternalErrorException")
}
