package mock

import (
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
)

type MockCognitoIdentityProviderClient struct {
	cognitoidentityprovideriface.CognitoIdentityProviderAPI
	DescribeUserPoolFunc func(poolInputData *cognitoidentityprovider.DescribeUserPoolInput) (*cognitoidentityprovider.DescribeUserPoolOutput, error)
	AdminCreateUserFunc  func(userInput *cognitoidentityprovider.AdminCreateUserInput) (*cognitoidentityprovider.AdminCreateUserOutput, error)
}

func (m *MockCognitoIdentityProviderClient) DescribeUserPool(poolInputData *cognitoidentityprovider.DescribeUserPoolInput) (*cognitoidentityprovider.DescribeUserPoolOutput, error) {
	return m.DescribeUserPoolFunc(poolInputData)
}

// AdminCreateUser function
func (m *MockCognitoIdentityProviderClient) AdminCreateUser(userInput *cognitoidentityprovider.AdminCreateUserInput) (*cognitoidentityprovider.AdminCreateUserOutput, error) {
	return m.AdminCreateUserFunc(userInput)
}
