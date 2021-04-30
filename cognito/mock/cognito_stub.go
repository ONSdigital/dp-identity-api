package mock

import (
	"errors"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
)

type CognitoIdentityProviderClientStub struct {
	cognitoidentityprovideriface.CognitoIdentityProviderAPI
	UserPools []string
	Users []User
}

func (m *CognitoIdentityProviderClientStub) DescribeUserPool(poolInputData *cognitoidentityprovider.DescribeUserPoolInput) (*cognitoidentityprovider.DescribeUserPoolOutput, error) {
	for _, v := range m.UserPools {
		if v == *poolInputData.UserPoolId {
			return nil, nil
		}
	}
	return nil, errors.New("Failed to load user pool data")
}

func (m *CognitoIdentityProviderClientStub) AdminGetUser(userInputData *cognitoidentityprovider.AdminGetUserInput) (*cognitoidentityprovider.AdminGetUserOutput, error) {
	for _, user := range m.Users {
		if user.username == *userInputData.Username {
			emailAttributeName := "email"
			givenNameAttributeName := "given_name"
			familyNameAttributeName := "family_name"
			getUserOutput := cognitoidentityprovider.AdminGetUserOutput{
				Enabled: &user.enabled,
				MFAOptions: []*cognitoidentityprovider.MFAOptionType{},
				PreferredMfaSetting: nil,
				UserAttributes: []*cognitoidentityprovider.AttributeType{
					{
						Name: &emailAttributeName,
						Value: &user.email,
					},
					{
						Name: &givenNameAttributeName,
						Value: &user.givenName,
					},
					{
						Name: &familyNameAttributeName,
						Value: &user.familyName,
					},
				},
				UserCreateDate: &user.created,
				UserLastModifiedDate: &user.created,
				UserMFASettingList: []*string{},
				UserStatus: &user.status,
				Username: &user.username,
			}
			return &getUserOutput, nil
		}
	}
	return nil, errors.New("Failed to load user data")
}
