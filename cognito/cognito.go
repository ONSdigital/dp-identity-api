package cognito

import (
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

// Client defines an interface for interaction with aws cognitoidentityprovider.
type Client interface {
	AdminAddUserToGroup(input *cognito.AdminAddUserToGroupInput) (*cognito.AdminAddUserToGroupOutput, error)
	AdminConfirmSignUp(input *cognito.AdminConfirmSignUpInput) (*cognito.AdminConfirmSignUpOutput, error)
	AdminCreateUser(input *cognito.AdminCreateUserInput) (*cognito.AdminCreateUserOutput, error)
	AdminDeleteUser(input *cognito.AdminDeleteUserInput) (*cognito.AdminDeleteUserOutput, error)
	AdminDisableUser(input *cognito.AdminDisableUserInput) (*cognito.AdminDisableUserOutput, error)
	AdminEnableUser(input *cognito.AdminEnableUserInput) (*cognito.AdminEnableUserOutput, error)
	AdminGetUser(input *cognito.AdminGetUserInput) (*cognito.AdminGetUserOutput, error)
	AdminListGroupsForUser(input *cognito.AdminListGroupsForUserInput) (*cognito.AdminListGroupsForUserOutput, error)
	AdminRemoveUserFromGroup(input *cognito.AdminRemoveUserFromGroupInput) (*cognito.AdminRemoveUserFromGroupOutput, error)
	AdminSetUserPassword(input *cognito.AdminSetUserPasswordInput) (*cognito.AdminSetUserPasswordOutput, error)
	AdminUpdateUserAttributes(input *cognito.AdminUpdateUserAttributesInput) (*cognito.AdminUpdateUserAttributesOutput, error)
	AdminUserGlobalSignOut(input *cognito.AdminUserGlobalSignOutInput) (*cognito.AdminUserGlobalSignOutOutput, error)
	ConfirmForgotPassword(input *cognito.ConfirmForgotPasswordInput) (*cognito.ConfirmForgotPasswordOutput, error)
	CreateGroup(input *cognito.CreateGroupInput) (*cognito.CreateGroupOutput, error)
	DeleteGroup(input *cognito.DeleteGroupInput) (*cognito.DeleteGroupOutput, error)
	DescribeUserPool(*cognito.DescribeUserPoolInput) (*cognito.DescribeUserPoolOutput, error)
	DescribeUserPoolClient(input *cognito.DescribeUserPoolClientInput) (*cognito.DescribeUserPoolClientOutput, error)
	ForgotPassword(input *cognito.ForgotPasswordInput) (*cognito.ForgotPasswordOutput, error)
	GetGroup(input *cognito.GetGroupInput) (*cognito.GetGroupOutput, error)
	GlobalSignOut(input *cognito.GlobalSignOutInput) (*cognito.GlobalSignOutOutput, error)
	InitiateAuth(input *cognito.InitiateAuthInput) (*cognito.InitiateAuthOutput, error)
	ListGroups(input *cognito.ListGroupsInput) (*cognito.ListGroupsOutput, error)
	ListUsers(input *cognito.ListUsersInput) (*cognito.ListUsersOutput, error)
	ListUsersInGroup(input *cognito.ListUsersInGroupInput) (*cognito.ListUsersInGroupOutput, error)
	RespondToAuthChallenge(input *cognito.RespondToAuthChallengeInput) (*cognito.RespondToAuthChallengeOutput, error)
	UpdateGroup(input *cognito.UpdateGroupInput) (*cognito.UpdateGroupOutput, error)
}
