package mock

import (
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
)

type MockCognitoIdentityProviderClient struct {
	cognitoidentityprovideriface.CognitoIdentityProviderAPI
	DescribeUserPoolFunc       func(poolInputData *cognitoidentityprovider.DescribeUserPoolInput) (*cognitoidentityprovider.DescribeUserPoolOutput, error)
	AdminCreateUserFunc        func(userInput *cognitoidentityprovider.AdminCreateUserInput) (*cognitoidentityprovider.AdminCreateUserOutput, error)
	GlobalSignOutFunc          func(signOutInput *cognitoidentityprovider.GlobalSignOutInput) (*cognitoidentityprovider.GlobalSignOutOutput, error)
	ListUsersFunc              func(usersInput *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error)
	InitiateAuthFunc           func(authInput *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error)
	AdminUserGlobalSignOutFunc func(adminUserGlobalSignOutInput *cognitoidentityprovider.AdminUserGlobalSignOutInput) (*cognitoidentityprovider.AdminUserGlobalSignOutOutput, error)
}

func (m *MockCognitoIdentityProviderClient) DescribeUserPool(poolInputData *cognitoidentityprovider.DescribeUserPoolInput) (*cognitoidentityprovider.DescribeUserPoolOutput, error) {
	return m.DescribeUserPoolFunc(poolInputData)
}

// AdminCreateUser function
func (m *MockCognitoIdentityProviderClient) AdminCreateUser(userInput *cognitoidentityprovider.AdminCreateUserInput) (*cognitoidentityprovider.AdminCreateUserOutput, error) {
	return m.AdminCreateUserFunc(userInput)
}

func (m *MockCognitoIdentityProviderClient) GlobalSignOut(signOutInput *cognitoidentityprovider.GlobalSignOutInput) (*cognitoidentityprovider.GlobalSignOutOutput, error) {
	return m.GlobalSignOutFunc(signOutInput)
}

func (m *MockCognitoIdentityProviderClient) ListUsers(usersInput *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error) {
	return m.ListUsersFunc(usersInput)
}

func (m *MockCognitoIdentityProviderClient) InitiateAuth(authInput *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	return m.InitiateAuthFunc(authInput)
}

func (m *MockCognitoIdentityProviderClient) AdminUserGlobalSignOut(adminUserGlobalSignOutInput *cognitoidentityprovider.AdminUserGlobalSignOutInput) (*cognitoidentityprovider.AdminUserGlobalSignOutOutput, error) {
	return m.AdminUserGlobalSignOutFunc(adminUserGlobalSignOutInput)
}
