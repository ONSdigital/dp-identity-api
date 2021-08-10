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
	ListUsersInGroup(input *cognito.ListUsersInGroupInput) (*cognito.ListUsersInGroupOutput, error)
	AdminRemoveUserFromGroup(input *cognito.AdminRemoveUserFromGroupInput) (*cognito.AdminRemoveUserFromGroupOutput, error)
<<<<<<< HEAD
	AdminConfirmSignUp(input *cognito.AdminConfirmSignUpInput) (*cognito.AdminConfirmSignUpOutput, error)
	AdminDeleteUser(input *cognito.AdminDeleteUserInput) (*cognito.AdminDeleteUserOutput, error)
	DeleteGroup(input *cognito.DeleteGroupInput) (*cognito.DeleteGroupOutput, error)
	AdminSetUserPassword(input *cognito.AdminSetUserPasswordInput) (*cognito.AdminSetUserPasswordOutput, error)
=======
>>>>>>> c33b6e55b0c5cd6b7c1acc9f792b7b283c167575
	AdminListGroupsForUser(input *cognito.AdminListGroupsForUserInput) (*cognito.AdminListGroupsForUserOutput, error)
}
