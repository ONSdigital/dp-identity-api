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
}
