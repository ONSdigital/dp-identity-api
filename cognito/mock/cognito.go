package mock

import (
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
)

type MockCognitoIdentityProviderClient struct {
	cognitoidentityprovideriface.CognitoIdentityProviderAPI
	DescribeUserPoolFunc          func(poolInputData *cognitoidentityprovider.DescribeUserPoolInput) (*cognitoidentityprovider.DescribeUserPoolOutput, error)
	AdminCreateUserFunc           func(userInput *cognitoidentityprovider.AdminCreateUserInput) (*cognitoidentityprovider.AdminCreateUserOutput, error)
	GlobalSignOutFunc             func(signOutInput *cognitoidentityprovider.GlobalSignOutInput) (*cognitoidentityprovider.GlobalSignOutOutput, error)
	ListUsersFunc                 func(usersInput *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error)
	InitiateAuthFunc              func(authInput *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error)
	AdminUserGlobalSignOutFunc    func(adminUserGlobalSignOutInput *cognitoidentityprovider.AdminUserGlobalSignOutInput) (*cognitoidentityprovider.AdminUserGlobalSignOutOutput, error)
	RespondToAuthChallengeFunc    func(input *cognitoidentityprovider.RespondToAuthChallengeInput) (*cognitoidentityprovider.RespondToAuthChallengeOutput, error)
	ConfirmForgotPasswordFunc     func(input *cognitoidentityprovider.ConfirmForgotPasswordInput) (*cognitoidentityprovider.ConfirmForgotPasswordOutput, error)
	ForgotPasswordFunc            func(input *cognitoidentityprovider.ForgotPasswordInput) (*cognitoidentityprovider.ForgotPasswordOutput, error)
	AdminGetUserFunc              func(input *cognitoidentityprovider.AdminGetUserInput) (*cognitoidentityprovider.AdminGetUserOutput, error)
	CreateGroupFunc               func(input *cognitoidentityprovider.CreateGroupInput) (*cognitoidentityprovider.CreateGroupOutput, error)
	GetGroupFunc                  func(input *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error)
	AdminUpdateUserAttributesFunc func(input *cognitoidentityprovider.AdminUpdateUserAttributesInput) (*cognitoidentityprovider.AdminUpdateUserAttributesOutput, error)
	AdminEnableUserFunc           func(input *cognitoidentityprovider.AdminEnableUserInput) (*cognitoidentityprovider.AdminEnableUserOutput, error)
	AdminDisableUserFunc          func(input *cognitoidentityprovider.AdminDisableUserInput) (*cognitoidentityprovider.AdminDisableUserOutput, error)
	AdminAddUserToGroupFunc       func(input *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error)
	ListUsersInGroupFunc          func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error)
	AdminRemoveUserFromGroupFunc  func(input *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error)
	AdminConfirmSignUpFunc        func(input *cognitoidentityprovider.AdminConfirmSignUpInput) (*cognitoidentityprovider.AdminConfirmSignUpOutput, error)
	AdminDeleteUserFunc           func(input *cognitoidentityprovider.AdminDeleteUserInput) (*cognitoidentityprovider.AdminDeleteUserOutput, error)
	DeleteGroupFunc               func(input *cognitoidentityprovider.DeleteGroupInput) (*cognitoidentityprovider.DeleteGroupOutput, error)
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

func (m *MockCognitoIdentityProviderClient) RespondToAuthChallenge(input *cognitoidentityprovider.RespondToAuthChallengeInput) (*cognitoidentityprovider.RespondToAuthChallengeOutput, error) {
	return m.RespondToAuthChallengeFunc(input)
}

func (m *MockCognitoIdentityProviderClient) ConfirmForgotPassword(input *cognitoidentityprovider.ConfirmForgotPasswordInput) (*cognitoidentityprovider.ConfirmForgotPasswordOutput, error) {
	return m.ConfirmForgotPasswordFunc(input)
}

func (m *MockCognitoIdentityProviderClient) ForgotPassword(input *cognitoidentityprovider.ForgotPasswordInput) (*cognitoidentityprovider.ForgotPasswordOutput, error) {
	return m.ForgotPasswordFunc(input)
}

func (m *MockCognitoIdentityProviderClient) AdminGetUser(input *cognitoidentityprovider.AdminGetUserInput) (*cognitoidentityprovider.AdminGetUserOutput, error) {
	return m.AdminGetUserFunc(input)
}

func (m *MockCognitoIdentityProviderClient) CreateGroup(input *cognitoidentityprovider.CreateGroupInput) (*cognitoidentityprovider.CreateGroupOutput, error) {
	return m.CreateGroupFunc(input)
}

func (m *MockCognitoIdentityProviderClient) GetGroup(input *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error) {
	return m.GetGroupFunc(input)
}

func (m *MockCognitoIdentityProviderClient) AdminUpdateUserAttributes(input *cognitoidentityprovider.AdminUpdateUserAttributesInput) (*cognitoidentityprovider.AdminUpdateUserAttributesOutput, error) {
	return m.AdminUpdateUserAttributesFunc(input)
}

func (m *MockCognitoIdentityProviderClient) AdminEnableUser(input *cognitoidentityprovider.AdminEnableUserInput) (*cognitoidentityprovider.AdminEnableUserOutput, error) {
	return m.AdminEnableUserFunc(input)
}

func (m *MockCognitoIdentityProviderClient) AdminDisableUser(input *cognitoidentityprovider.AdminDisableUserInput) (*cognitoidentityprovider.AdminDisableUserOutput, error) {
	return m.AdminDisableUserFunc(input)
}

func (m *MockCognitoIdentityProviderClient) AdminAddUserToGroup(input *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
	return m.AdminAddUserToGroupFunc(input)
}

func (m *MockCognitoIdentityProviderClient) ListUsersInGroup(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
	return m.ListUsersInGroupFunc(input)
}

func (m *MockCognitoIdentityProviderClient) AdminRemoveUserFromGroup(input *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error) {
	return m.AdminRemoveUserFromGroupFunc(input)
}

func (m *MockCognitoIdentityProviderClient) AdminConfirmSignUp(input *cognitoidentityprovider.AdminConfirmSignUpInput) (*cognitoidentityprovider.AdminConfirmSignUpOutput, error) {
	return m.AdminConfirmSignUpFunc(input)
}

func (m *MockCognitoIdentityProviderClient) AdminDeleteUser(input *cognitoidentityprovider.AdminDeleteUserInput) (*cognitoidentityprovider.AdminDeleteUserOutput, error) {
	return m.AdminDeleteUserFunc(input)
}

func (m *MockCognitoIdentityProviderClient) DeleteGroup(input *cognitoidentityprovider.DeleteGroupInput) (*cognitoidentityprovider.DeleteGroupOutput, error) {
	return m.DeleteGroupFunc(input)
}
