package cognito

import (
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

// Client defines an interface for interaction with aws cognitoidentityprovider.
type Client interface {
	DescribeUserPool(*cognito.DescribeUserPoolInput) (*cognito.DescribeUserPoolOutput, error)
	AdminCreateUser(input *cognito.AdminCreateUserInput) (*cognito.AdminCreateUserOutput, error)
	InitiateAuth(input *cognito.InitiateAuthInput) (*cognito.InitiateAuthOutput, error)
	GlobalSignOut(input *cognito.GlobalSignOutInput) (*cognito.GlobalSignOutOutput, error)
	ListUsers(input *cognito.ListUsersInput) (*cognito.ListUsersOutput, error)
	AdminUserGlobalSignOut(input *cognito.AdminUserGlobalSignOutInput) (*cognito.AdminUserGlobalSignOutOutput, error)
	RespondToAuthChallenge(input *cognito.RespondToAuthChallengeInput) (*cognito.RespondToAuthChallengeOutput, error)
	ConfirmForgotPassword(input *cognito.ConfirmForgotPasswordInput) (*cognito.ConfirmForgotPasswordOutput, error)
	ForgotPassword(input *cognito.ForgotPasswordInput) (*cognito.ForgotPasswordOutput, error)
	AdminGetUser(input *cognito.AdminGetUserInput) (*cognito.AdminGetUserOutput, error)
	CreateGroup(input *cognito.CreateGroupInput) (*cognito.CreateGroupOutput, error)
	GetGroup(input *cognito.GetGroupInput) (*cognito.GetGroupOutput, error)
	AdminUpdateUserAttributes(input *cognito.AdminUpdateUserAttributesInput) (*cognito.AdminUpdateUserAttributesOutput, error)
	AdminEnableUser(input *cognito.AdminEnableUserInput) (*cognito.AdminEnableUserOutput, error)
	AdminDisableUser(input *cognito.AdminDisableUserInput) (*cognito.AdminDisableUserOutput, error)
	AdminAddUserToGroup(input *cognito.AdminAddUserToGroupInput) (*cognito.AdminAddUserToGroupOutput, error)
}
