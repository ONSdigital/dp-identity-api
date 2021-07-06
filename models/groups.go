package models

import "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

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
