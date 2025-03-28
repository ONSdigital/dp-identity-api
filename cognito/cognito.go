package cognito

import (
	"context"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

// Client defines an interface for interaction with aws cognitoidentityprovider.
// NB. For a full list of the Client interface functions go to https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider#Client
type Client interface {
	AdminAddUserToGroup(ctx context.Context, params *cognito.AdminAddUserToGroupInput, optFns ...func(*cognito.Options)) (*cognito.AdminAddUserToGroupOutput, error)
	//AdminConfirmSignUp(ctx context.Context, params *cognito.AdminConfirmSignUpInput, optFns ...func(*cognito.Options)) (*cognito.AdminConfirmSignUpOutput, error)
	DescribeUserPool(ctx context.Context, params *cognito.DescribeUserPoolInput, optFns ...func(*cognito.Options)) (*cognito.DescribeUserPoolOutput, error)
	GetGroup(ctx context.Context, params *cognito.GetGroupInput, optFns ...func(*cognito.Options)) (*cognito.GetGroupOutput, error)
	CreateGroup(ctx context.Context, params *cognito.CreateGroupInput, optFns ...func(*cognito.Options)) (*cognito.CreateGroupOutput, error)
	UpdateGroup(ctx context.Context, params *cognito.UpdateGroupInput, optFns ...func(*cognito.Options)) (*cognito.UpdateGroupOutput, error)
	ListUsersInGroup(ctx context.Context, params *cognito.ListUsersInGroupInput, optFns ...func(*cognito.Options)) (*cognito.ListUsersInGroupOutput, error)
	ListGroups(ctx context.Context, params *cognito.ListGroupsInput, optFns ...func(*cognito.Options)) (*cognito.ListGroupsOutput, error)
	DeleteGroup(ctx context.Context, params *cognito.DeleteGroupInput, optFns ...func(*cognito.Options)) (*cognito.DeleteGroupOutput, error)
	AdminRemoveUserFromGroup(ctx context.Context, params *cognito.AdminRemoveUserFromGroupInput, optFns ...func(*cognito.Options)) (*cognito.AdminRemoveUserFromGroupOutput, error)
	InitiateAuth(ctx context.Context, params *cognito.InitiateAuthInput, optFns ...func(*cognito.Options)) (*cognito.InitiateAuthOutput, error)
	DescribeUserPoolClient(ctx context.Context, params *cognito.DescribeUserPoolClientInput, optFns ...func(*cognito.Options)) (*cognito.DescribeUserPoolClientOutput, error)
	GlobalSignOut(ctx context.Context, params *cognito.GlobalSignOutInput, optFns ...func(*cognito.Options)) (*cognito.GlobalSignOutOutput, error)
	AdminUserGlobalSignOut(ctx context.Context, params *cognito.AdminUserGlobalSignOutInput, optFns ...func(*cognito.Options)) (*cognito.AdminUserGlobalSignOutOutput, error)
	ListUsers(ctx context.Context, params *cognito.ListUsersInput, optFns ...func(*cognito.Options)) (*cognito.ListUsersOutput, error)
	AdminCreateUser(ctx context.Context, params *cognito.AdminCreateUserInput, optFns ...func(*cognito.Options)) (*cognito.AdminCreateUserOutput, error)
	AdminGetUser(ctx context.Context, params *cognito.AdminGetUserInput, optFns ...func(*cognito.Options)) (*cognito.AdminGetUserOutput, error)
	AdminEnableUser(ctx context.Context, params *cognito.AdminEnableUserInput, optFns ...func(*cognito.Options)) (*cognito.AdminEnableUserOutput, error)
	AdminDisableUser(ctx context.Context, params *cognito.AdminDisableUserInput, optFns ...func(*cognito.Options)) (*cognito.AdminDisableUserOutput, error)
	AdminUpdateUserAttributes(ctx context.Context, params *cognito.AdminUpdateUserAttributesInput, optFns ...func(*cognito.Options)) (*cognito.AdminUpdateUserAttributesOutput, error)
	RespondToAuthChallenge(ctx context.Context, params *cognito.RespondToAuthChallengeInput, optFns ...func(*cognito.Options)) (*cognito.RespondToAuthChallengeOutput, error)
	ConfirmForgotPassword(ctx context.Context, params *cognito.ConfirmForgotPasswordInput, optFns ...func(*cognito.Options)) (*cognito.ConfirmForgotPasswordOutput, error)
	ForgotPassword(ctx context.Context, params *cognito.ForgotPasswordInput, optFns ...func(*cognito.Options)) (*cognito.ForgotPasswordOutput, error)
	AdminListGroupsForUser(ctx context.Context, params *cognito.AdminListGroupsForUserInput, optFns ...func(*cognito.Options)) (*cognito.AdminListGroupsForUserOutput, error)
	//AdminConfirmSignUp(input *cognito.AdminConfirmSignUpInput) (*cognito.AdminConfirmSignUpOutput, error)
	//AdminDeleteUser(input *cognito.AdminDeleteUserInput) (*cognito.AdminDeleteUserOutput, error)
	//AdminSetUserPassword(input *cognito.AdminSetUserPasswordInput) (*cognito.AdminSetUserPasswordOutput, error)
	//DescribeUserPool(*cognito.DescribeUserPoolInput) (*cognito.DescribeUserPoolOutput, error)
	//GetGroup(input *cognito.GetGroupInput) (cognito.GetGroupOutput, error)
}
