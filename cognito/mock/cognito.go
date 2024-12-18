package mock

import (
	"context"

	"github.com/ONSdigital/dp-identity-api/v2/models"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
)

type MockCognitoIdentityProviderClient struct { //nolint:revive // Mock here is used to identify this as the mock version during implementation
	cognitoidentityprovideriface.CognitoIdentityProviderAPI
	AddUserToGroupfunc            func(ctx context.Context, group models.Group, userId string) (*models.UsersList, *models.ErrorResponse)
	AdminAddUserToGroupFunc       func(input *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error)
	AdminConfirmSignUpFunc        func(input *cognitoidentityprovider.AdminConfirmSignUpInput) (*cognitoidentityprovider.AdminConfirmSignUpOutput, error)
	AdminCreateUserFunc           func(userInput *cognitoidentityprovider.AdminCreateUserInput) (*cognitoidentityprovider.AdminCreateUserOutput, error)
	AdminDeleteUserFunc           func(input *cognitoidentityprovider.AdminDeleteUserInput) (*cognitoidentityprovider.AdminDeleteUserOutput, error)
	AdminDisableUserFunc          func(input *cognitoidentityprovider.AdminDisableUserInput) (*cognitoidentityprovider.AdminDisableUserOutput, error)
	AdminEnableUserFunc           func(input *cognitoidentityprovider.AdminEnableUserInput) (*cognitoidentityprovider.AdminEnableUserOutput, error)
	AdminGetUserFunc              func(input *cognitoidentityprovider.AdminGetUserInput) (*cognitoidentityprovider.AdminGetUserOutput, error)
	AdminRemoveUserFromGroupFunc  func(input *cognitoidentityprovider.AdminRemoveUserFromGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error)
	AdminSetUserPasswordFunc      func(input *cognitoidentityprovider.AdminSetUserPasswordInput) (*cognitoidentityprovider.AdminSetUserPasswordOutput, error)
	AdminUpdateUserAttributesFunc func(input *cognitoidentityprovider.AdminUpdateUserAttributesInput) (*cognitoidentityprovider.AdminUpdateUserAttributesOutput, error)
	AdminUserGlobalSignOutFunc    func(adminUserGlobalSignOutInput *cognitoidentityprovider.AdminUserGlobalSignOutInput) (*cognitoidentityprovider.AdminUserGlobalSignOutOutput, error)
	ConfirmForgotPasswordFunc     func(input *cognitoidentityprovider.ConfirmForgotPasswordInput) (*cognitoidentityprovider.ConfirmForgotPasswordOutput, error)
	CreateGroupFunc               func(input *cognitoidentityprovider.CreateGroupInput) (*cognitoidentityprovider.CreateGroupOutput, error)
	DeleteGroupFunc               func(input *cognitoidentityprovider.DeleteGroupInput) (*cognitoidentityprovider.DeleteGroupOutput, error)
	DescribeUserPoolClientFunc    func(input *cognitoidentityprovider.DescribeUserPoolClientInput) (*cognitoidentityprovider.DescribeUserPoolClientOutput, error)
	DescribeUserPoolFunc          func(poolInputData *cognitoidentityprovider.DescribeUserPoolInput) (*cognitoidentityprovider.DescribeUserPoolOutput, error)
	ForgotPasswordFunc            func(input *cognitoidentityprovider.ForgotPasswordInput) (*cognitoidentityprovider.ForgotPasswordOutput, error)
	GetGroupFunc                  func(input *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error)
	GlobalSignOutFunc             func(signOutInput *cognitoidentityprovider.GlobalSignOutInput) (*cognitoidentityprovider.GlobalSignOutOutput, error)
	InitiateAuthFunc              func(authInput *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error)
	ListGroupsForUserFunc         func(input *cognitoidentityprovider.AdminListGroupsForUserInput) (*cognitoidentityprovider.AdminListGroupsForUserOutput, error)
	ListGroupsFunc                func(input *cognitoidentityprovider.ListGroupsInput) (*cognitoidentityprovider.ListGroupsOutput, error)
	ListUsersFunc                 func(usersInput *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error)
	ListUsersInGroupFunc          func(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error)
	RemoveUserFromGroupfunc       func(ctx context.Context, group models.Group, userId string) (*models.UsersList, *models.ErrorResponse)
	RespondToAuthChallengeFunc    func(input *cognitoidentityprovider.RespondToAuthChallengeInput) (*cognitoidentityprovider.RespondToAuthChallengeOutput, error)
	SetGroupUsersfunc             func(ctx context.Context, group models.Group, users models.UsersList) (*models.UsersList, *models.ErrorResponse)
	UpdateGroupFunc               func(input *cognitoidentityprovider.UpdateGroupInput) (*cognitoidentityprovider.UpdateGroupOutput, error)
	ValidateAddRemoveUserFunc     func(ctx context.Context, userId string) error
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

func (m *MockCognitoIdentityProviderClient) AdminSetUserPassword(input *cognitoidentityprovider.AdminSetUserPasswordInput) (*cognitoidentityprovider.AdminSetUserPasswordOutput, error) {
	return m.AdminSetUserPasswordFunc(input)
}

func (m *MockCognitoIdentityProviderClient) AdminListGroupsForUser(input *cognitoidentityprovider.AdminListGroupsForUserInput) (*cognitoidentityprovider.AdminListGroupsForUserOutput, error) {
	return m.ListGroupsForUserFunc(input)
}

func (m *MockCognitoIdentityProviderClient) DescribeUserPoolClient(input *cognitoidentityprovider.DescribeUserPoolClientInput) (*cognitoidentityprovider.DescribeUserPoolClientOutput, error) {
	return m.DescribeUserPoolClientFunc(input)
}

//	func (m *MockCognitoIdentityProviderClient) ListGroupsUsers(input *cognitoidentityprovider.ListGroupsOutput) (*[]models.ListGroupUsersType, error) {
//		return m.ListGroupsUsersFunc(input)
//	}
func (m *MockCognitoIdentityProviderClient) ListGroups(input *cognitoidentityprovider.ListGroupsInput) (*cognitoidentityprovider.ListGroupsOutput, error) {
	return m.ListGroupsFunc(input)
}

func (m *MockCognitoIdentityProviderClient) UpdateGroup(input *cognitoidentityprovider.UpdateGroupInput) (*cognitoidentityprovider.UpdateGroupOutput, error) {
	return m.UpdateGroupFunc(input)
}

func (m *MockCognitoIdentityProviderClient) SetGroupUsers(ctx context.Context, group models.Group, users models.UsersList) (*models.UsersList, *models.ErrorResponse) {
	return m.SetGroupUsersfunc(ctx, group, users)
}

func (m *MockCognitoIdentityProviderClient) AddUserToGroup(ctx context.Context, group models.Group, userID string) (*models.UsersList, *models.ErrorResponse) {
	return m.AddUserToGroupfunc(ctx, group, userID)
}

func (m *MockCognitoIdentityProviderClient) RemoveUserFromGroup(ctx context.Context, group models.Group, userID string) (*models.UsersList, *models.ErrorResponse) {
	return m.RemoveUserFromGroupfunc(ctx, group, userID)
}
