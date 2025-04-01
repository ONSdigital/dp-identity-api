package mock

import (
	"context"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	//"github.com/aws/aws-sdk-go-v2/service/cognito/cognitoiface"
)

type MockCognitoIdentityProviderClient struct { //nolint:revive // Mock here is used to identify this as the mock version during implementation
	//cognitoiface.CognitoIdentityProviderAPI
	cognito.Client
	//AddUserToGroupfunc            func(ctx context.Context, group models.Group, userId string) (*models.UsersList, *models.ErrorResponse)

	AdminAddUserToGroupFunc func(ctx context.Context, input *cognito.AdminAddUserToGroupInput, optFns ...func(*cognito.Options)) (*cognito.AdminAddUserToGroupOutput, error)

	AdminConfirmSignUpFunc func(ctx context.Context, input *cognito.AdminConfirmSignUpInput, optFns ...func(*cognito.Options)) (*cognito.AdminConfirmSignUpOutput, error)
	AdminCreateUserFunc    func(ctx context.Context, userInput *cognito.AdminCreateUserInput, optFns ...func(*cognito.Options)) (*cognito.AdminCreateUserOutput, error)
	//AdminDeleteUserFunc           func(input *cognito.AdminDeleteUserInput) (*cognito.AdminDeleteUserOutput, error)
	AdminDeleteUserFunc func(ctx context.Context, input *cognito.AdminDeleteUserInput, optFns ...func(*cognito.Options)) (*cognito.AdminDeleteUserOutput, error)
	//AdminDisableUserFunc          func(input *cognito.AdminDisableUserInput) (*cognito.AdminDisableUserOutput, error)
	AdminDisableUserFunc func(ctx context.Context, input *cognito.AdminDisableUserInput, optFns ...func(*cognito.Options)) (*cognito.AdminDisableUserOutput, error)
	//AdminEnableUserFunc           func(input *cognito.AdminEnableUserInput) (*cognito.AdminEnableUserOutput, error)
	AdminEnableUserFunc func(ctx context.Context, input *cognito.AdminEnableUserInput, optFns ...func(*cognito.Options)) (*cognito.AdminEnableUserOutput, error)
	//AdminGetUserFunc              func(input *cognito.AdminGetUserInput) (*cognito.AdminGetUserOutput, error)
	AdminGetUserFunc func(ctx context.Context, params *cognito.AdminGetUserInput, optFns ...func(*cognito.Options)) (*cognito.AdminGetUserOutput, error)
	//AdminRemoveUserFromGroupFunc  func(input *cognito.AdminRemoveUserFromGroupInput) (*cognito.AdminRemoveUserFromGroupOutput, error)
	ListGroupsForUserFunc        func(ctx context.Context, input *cognito.AdminListGroupsForUserInput, optFns ...func(*cognito.Options)) (*cognito.AdminListGroupsForUserOutput, error)
	AdminRemoveUserFromGroupFunc func(ctx context.Context, input *cognito.AdminRemoveUserFromGroupInput, optFns ...func(*cognito.Options)) (*cognito.AdminRemoveUserFromGroupOutput, error)
	//AdminSetUserPasswordFunc      func(input *cognito.AdminSetUserPasswordInput) (*cognito.AdminSetUserPasswordOutput, error)
	AdminSetUserPasswordFunc func(ctx context.Context, input *cognito.AdminSetUserPasswordInput, optFns ...func(*cognito.Options)) (*cognito.AdminSetUserPasswordOutput, error)
	//AdminUpdateUserAttributesFunc func(input *cognito.AdminUpdateUserAttributesInput) (*cognito.AdminUpdateUserAttributesOutput, error)
	AdminUpdateUserAttributesFunc func(ctx context.Context, params *cognito.AdminUpdateUserAttributesInput, optFns ...func(*cognito.Options)) (*cognito.AdminUpdateUserAttributesOutput, error)
	//AdminUserGlobalSignOutFunc    func(adminUserGlobalSignOutInput *cognito.AdminUserGlobalSignOutInput) (*cognito.AdminUserGlobalSignOutOutput, error)
	AdminUserGlobalSignOutFunc func(ctx context.Context, params *cognito.AdminUserGlobalSignOutInput, optFns ...func(*cognito.Options)) (*cognito.AdminUserGlobalSignOutOutput, error)
	//ConfirmForgotPasswordFunc     func(input *cognito.ConfirmForgotPasswordInput) (*cognito.ConfirmForgotPasswordOutput, error)
	ConfirmForgotPasswordFunc func(ctx context.Context, params *cognito.ConfirmForgotPasswordInput, optFns ...func(*cognito.Options)) (*cognito.ConfirmForgotPasswordOutput, error)
	CreateGroupFunc           func(ctx context.Context, input *cognito.CreateGroupInput, optFns ...func(*cognito.Options)) (*cognito.CreateGroupOutput, error)
	//DeleteGroupFunc            func(input *cognito.DeleteGroupInput) (*cognito.DeleteGroupOutput, error)
	DeleteGroupFunc func(ctx context.Context, input *cognito.DeleteGroupInput, optFns ...func(*cognito.Options)) (*cognito.DeleteGroupOutput, error)
	//DescribeUserPoolClientFunc func(input *cognito.DescribeUserPoolClientInput) (*cognito.DescribeUserPoolClientOutput, error)
	DescribeUserPoolClientFunc func(ctx context.Context, input *cognito.DescribeUserPoolClientInput, optFns ...func(*cognito.Options)) (*cognito.DescribeUserPoolClientOutput, error)
	//DescribeUserPoolFunc       func(poolInputData *cognito.DescribeUserPoolInput) (*cognito.DescribeUserPoolOutput, error)
	DescribeUserPoolFunc func(ctx context.Context, poolInputData *cognito.DescribeUserPoolInput, optFns ...func(*cognito.Options)) (*cognito.DescribeUserPoolOutput, error)
	//ForgotPasswordFunc         func(input *cognito.ForgotPasswordInput) (*cognito.ForgotPasswordOutput, error)
	ForgotPasswordFunc func(ctx context.Context, params *cognito.ForgotPasswordInput, optFns ...func(*cognito.Options)) (*cognito.ForgotPasswordOutput, error)
	//GetGroupFunc               func(input *cognito.GetGroupInput) (*cognito.GetGroupOutput, error)
	GetGroupFunc func(ctx context.Context, params *cognito.GetGroupInput, optFns ...func(*cognito.Options)) (*cognito.GetGroupOutput, error)
	//GlobalSignOutFunc          func(signOutInput *cognito.GlobalSignOutInput) (*cognito.GlobalSignOutOutput, error)
	GlobalSignOutFunc func(ctx context.Context, signOutInput *cognito.GlobalSignOutInput, optFns ...func(*cognito.Options)) (*cognito.GlobalSignOutOutput, error)
	//InitiateAuthFunc           func(authInput *cognito.InitiateAuthInput) (*cognito.InitiateAuthOutput, error)
	InitiateAuthFunc func(ctx context.Context, params *cognito.InitiateAuthInput, optFns ...func(*cognito.Options)) (*cognito.InitiateAuthOutput, error)
	//ListGroupsForUserFunc      func(input *cognito.AdminListGroupsForUserInput) (*cognito.AdminListGroupsForUserOutput, error)
	//ListGroupsFunc             func(input *cognito.ListGroupsInput) (*cognito.ListGroupsOutput, error)
	ListGroupsFunc func(ctx context.Context, input *cognito.ListGroupsInput, optFns ...func(*cognito.Options)) (*cognito.ListGroupsOutput, error)
	//ListUsersFunc              func(usersInput *cognito.ListUsersInput) (*cognito.ListUsersOutput, error)
	ListUsersFunc func(ctx context.Context, usersInput *cognito.ListUsersInput, optFns ...func(*cognito.Options)) (*cognito.ListUsersOutput, error)
	//ListUsersInGroupFunc       func(input *cognito.ListUsersInGroupInput) (*cognito.ListUsersInGroupOutput, error)
	ListUsersInGroupFunc func(ctx context.Context, input *cognito.ListUsersInGroupInput, optFns ...func(*cognito.Options)) (*cognito.ListUsersInGroupOutput, error)
	//RemoveUserFromGroupfunc    func(ctx context.Context, group models.Group, userId string) (*models.UsersList, *models.ErrorResponse)
	//RespondToAuthChallengeFunc func(input *cognito.RespondToAuthChallengeInput) (*cognito.RespondToAuthChallengeOutput, error)
	RespondToAuthChallengeFunc func(ctx context.Context, params *cognito.RespondToAuthChallengeInput, optFns ...func(*cognito.Options)) (*cognito.RespondToAuthChallengeOutput, error)
	//SetGroupUsersfunc          func(ctx context.Context, group models.Group, users models.UsersList) (*models.UsersList, *models.ErrorResponse)
	//UpdateGroupFunc            func(input *cognito.UpdateGroupInput) (*cognito.UpdateGroupOutput, error)
	UpdateGroupFunc func(ctx context.Context, input *cognito.UpdateGroupInput, optFns ...func(*cognito.Options)) (*cognito.UpdateGroupOutput, error)
	//ValidateAddRemoveUserFunc  func(ctx context.Context, userId string) error
}

func (m *MockCognitoIdentityProviderClient) DescribeUserPool(ctx context.Context, poolInputData *cognito.DescribeUserPoolInput, _ ...func(*cognito.Options)) (*cognito.DescribeUserPoolOutput, error) {
	return m.DescribeUserPoolFunc(ctx, poolInputData, nil)
}

// AdminCreateUser function
func (m *MockCognitoIdentityProviderClient) AdminCreateUser(ctx context.Context, userInput *cognito.AdminCreateUserInput, _ ...func(*cognito.Options)) (*cognito.AdminCreateUserOutput, error) {
	return m.AdminCreateUserFunc(ctx, userInput, nil)
}

func (m *MockCognitoIdentityProviderClient) GlobalSignOut(ctx context.Context, signOutInput *cognito.GlobalSignOutInput, _ ...func(*cognito.Options)) (*cognito.GlobalSignOutOutput, error) {
	return m.GlobalSignOutFunc(ctx, signOutInput, nil)
}

func (m *MockCognitoIdentityProviderClient) ListUsers(ctx context.Context, usersInput *cognito.ListUsersInput, _ ...func(*cognito.Options)) (*cognito.ListUsersOutput, error) {
	return m.ListUsersFunc(ctx, usersInput, nil)
}

func (m *MockCognitoIdentityProviderClient) InitiateAuth(ctx context.Context, authInput *cognito.InitiateAuthInput, _ ...func(*cognito.Options)) (*cognito.InitiateAuthOutput, error) {
	return m.InitiateAuthFunc(ctx, authInput, nil)
}

func (m *MockCognitoIdentityProviderClient) AdminUserGlobalSignOut(ctx context.Context, adminUserGlobalSignOutInput *cognito.AdminUserGlobalSignOutInput, _ ...func(*cognito.Options)) (*cognito.AdminUserGlobalSignOutOutput, error) {
	return m.AdminUserGlobalSignOutFunc(ctx, adminUserGlobalSignOutInput, nil)
}

func (m *MockCognitoIdentityProviderClient) RespondToAuthChallenge(ctx context.Context, input *cognito.RespondToAuthChallengeInput, _ ...func(*cognito.Options)) (*cognito.RespondToAuthChallengeOutput, error) {
	return m.RespondToAuthChallengeFunc(ctx, input, nil)
}

func (m *MockCognitoIdentityProviderClient) ConfirmForgotPassword(ctx context.Context, input *cognito.ConfirmForgotPasswordInput, _ ...func(*cognito.Options)) (*cognito.ConfirmForgotPasswordOutput, error) {
	return m.ConfirmForgotPasswordFunc(ctx, input, nil)
}

func (m *MockCognitoIdentityProviderClient) ForgotPassword(ctx context.Context, input *cognito.ForgotPasswordInput, _ ...func(*cognito.Options)) (*cognito.ForgotPasswordOutput, error) {
	return m.ForgotPasswordFunc(ctx, input, nil)
}

func (m *MockCognitoIdentityProviderClient) AdminGetUser(ctx context.Context, input *cognito.AdminGetUserInput, _ ...func(*cognito.Options)) (*cognito.AdminGetUserOutput, error) {
	return m.AdminGetUserFunc(ctx, input, nil)
}

func (m *MockCognitoIdentityProviderClient) CreateGroup(ctx context.Context, input *cognito.CreateGroupInput, _ ...func(*cognito.Options)) (*cognito.CreateGroupOutput, error) {
	return m.CreateGroupFunc(ctx, input, nil)
}

func (m *MockCognitoIdentityProviderClient) GetGroup(ctx context.Context, input *cognito.GetGroupInput, _ ...func(*cognito.Options)) (*cognito.GetGroupOutput, error) {
	return m.GetGroupFunc(ctx, input, nil)
}

func (m *MockCognitoIdentityProviderClient) AdminUpdateUserAttributes(ctx context.Context, input *cognito.AdminUpdateUserAttributesInput, _ ...func(*cognito.Options)) (*cognito.AdminUpdateUserAttributesOutput, error) {
	return m.AdminUpdateUserAttributesFunc(ctx, input, nil)
}

func (m *MockCognitoIdentityProviderClient) AdminEnableUser(ctx context.Context, input *cognito.AdminEnableUserInput, _ ...func(*cognito.Options)) (*cognito.AdminEnableUserOutput, error) {
	return m.AdminEnableUserFunc(ctx, input, nil)
}

func (m *MockCognitoIdentityProviderClient) AdminDisableUser(ctx context.Context, input *cognito.AdminDisableUserInput, _ ...func(*cognito.Options)) (*cognito.AdminDisableUserOutput, error) {
	return m.AdminDisableUserFunc(ctx, input, nil)
}

func (m *MockCognitoIdentityProviderClient) AdminAddUserToGroup(ctx context.Context, input *cognito.AdminAddUserToGroupInput, _ ...func(*cognito.Options)) (*cognito.AdminAddUserToGroupOutput, error) {
	return m.AdminAddUserToGroupFunc(ctx, input, nil)
}

func (m *MockCognitoIdentityProviderClient) ListUsersInGroup(ctx context.Context, input *cognito.ListUsersInGroupInput, _ ...func(*cognito.Options)) (*cognito.ListUsersInGroupOutput, error) {
	return m.ListUsersInGroupFunc(ctx, input, nil)
}

func (m *MockCognitoIdentityProviderClient) AdminRemoveUserFromGroup(ctx context.Context, input *cognito.AdminRemoveUserFromGroupInput, _ ...func(*cognito.Options)) (*cognito.AdminRemoveUserFromGroupOutput, error) {
	return m.AdminRemoveUserFromGroupFunc(ctx, input, nil)
}

func (m *MockCognitoIdentityProviderClient) AdminConfirmSignUp(ctx context.Context, input *cognito.AdminConfirmSignUpInput) (*cognito.AdminConfirmSignUpOutput, error) {
	return m.AdminConfirmSignUpFunc(ctx, input)
}

func (m *MockCognitoIdentityProviderClient) AdminDeleteUser(ctx context.Context, input *cognito.AdminDeleteUserInput, _ ...func(*cognito.Options)) (*cognito.AdminDeleteUserOutput, error) {
	return m.AdminDeleteUserFunc(ctx, input)
}

func (m *MockCognitoIdentityProviderClient) DeleteGroup(ctx context.Context, input *cognito.DeleteGroupInput, _ ...func(*cognito.Options)) (*cognito.DeleteGroupOutput, error) {
	return m.DeleteGroupFunc(ctx, input, nil)
}

func (m *MockCognitoIdentityProviderClient) AdminSetUserPassword(ctx context.Context, input *cognito.AdminSetUserPasswordInput, _ ...func(*cognito.Options)) (*cognito.AdminSetUserPasswordOutput, error) {
	return m.AdminSetUserPasswordFunc(ctx, input, nil)
}

func (m *MockCognitoIdentityProviderClient) AdminListGroupsForUser(ctx context.Context, input *cognito.AdminListGroupsForUserInput, _ ...func(*cognito.Options)) (*cognito.AdminListGroupsForUserOutput, error) {
	return m.ListGroupsForUserFunc(ctx, input, nil)
}

func (m *MockCognitoIdentityProviderClient) DescribeUserPoolClient(ctx context.Context, input *cognito.DescribeUserPoolClientInput, _ ...func(*cognito.Options)) (*cognito.DescribeUserPoolClientOutput, error) {
	return m.DescribeUserPoolClientFunc(ctx, input, nil)
}

//	func (m *MockCognitoIdentityProviderClient) ListGroupsUsers(input *cognito.ListGroupsOutput) (*[]models.ListGroupUsersType, error) {
//		return m.ListGroupsUsersFunc(input)
//	}
func (m *MockCognitoIdentityProviderClient) ListGroups(ctx context.Context, input *cognito.ListGroupsInput, _ ...func(*cognito.Options)) (*cognito.ListGroupsOutput, error) {
	return m.ListGroupsFunc(ctx, input, nil)
}

func (m *MockCognitoIdentityProviderClient) UpdateGroup(ctx context.Context, input *cognito.UpdateGroupInput, _ ...func(*cognito.Options)) (*cognito.UpdateGroupOutput, error) {
	return m.UpdateGroupFunc(ctx, input, nil)
}

//func (m *MockCognitoIdentityProviderClient) SetGroupUsers(ctx context.Context, group models.Group, users models.UsersList) (*models.UsersList, *models.ErrorResponse) {
//	return m.SetGroupUsersfunc(ctx, group, users)
//}

//func (m *MockCognitoIdentityProviderClient) AddUserToGroup(ctx context.Context, group models.Group, userID string) (*models.UsersList, *models.ErrorResponse) {
//	return m.AddUserToGroupfunc(ctx, group, userID)
//}

//func (m *MockCognitoIdentityProviderClient) RemoveUserFromGroup(ctx context.Context, group models.Group, userID string) (*models.UsersList, *models.ErrorResponse) {
//	return m.RemoveUserFromGroupfunc(ctx, group, userID)
//}
