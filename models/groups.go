package models

import (
	"context"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

const (
	AdminRoleGroup     = "role-admin"
	PublisherRoleGroup = "role-publisher"
)

//Type to map for the Cognito GroupType object
type Group struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Precedence  int64  `json:"precedence"`
}

// Constructor for a new instance of the admin role group
func NewAdminRoleGroup() Group {
	return Group{
		Name:        AdminRoleGroup,
		Description: "The publishing admins",
		Precedence:  1,
	}
}

// Constructor for a new instance of the publisher role group
func NewPublisherRoleGroup() Group {
	return Group{
		Name:        PublisherRoleGroup,
		Description: "The publishers",
		Precedence:  2,
	}
}

// ValidateAddUser validates the required fields for adding a user to a group, returns validation errors for anything that fails
func (g *Group) ValidateAddUser(ctx context.Context, userId string) []error {
	var validationErrs []error
	if g.Name == "" {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidGroupNameError, MissingGroupNameErrorDescription))
	}

	if userId == "" {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidUserIdError, MissingUserIdErrorDescription))
	}
	return validationErrs
}

// BuildCreateGroupRequest builds a correctly populated CreateGroupInput object using the Groups values
func (g *Group) BuildCreateGroupRequest(userPoolId string) *cognitoidentityprovider.CreateGroupInput {
	return &cognitoidentityprovider.CreateGroupInput{
		GroupName:   &g.Name,
		Description: &g.Description,
		Precedence:  &g.Precedence,
		UserPoolId:  &userPoolId,
	}
}

// BuildCreateGroupRequest builds a correctly populated GetGroupInput object using the Groups values
func (g *Group) BuildGetGroupRequest(userPoolId string) *cognitoidentityprovider.GetGroupInput {
	return &cognitoidentityprovider.GetGroupInput{
		GroupName:  &g.Name,
		UserPoolId: &userPoolId,
	}
}

// BuildAddUserToGroupRequest builds a correctly populated AdminAddUserToGroupInput object
func (g *Group) BuildAddUserToGroupRequest(userPoolId, userId string) *cognitoidentityprovider.AdminAddUserToGroupInput {
	return &cognitoidentityprovider.AdminAddUserToGroupInput{
		GroupName:  &g.Name,
		UserPoolId: &userPoolId,
		Username:   &userId,
	}
}
