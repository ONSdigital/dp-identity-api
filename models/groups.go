package models

import "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

const (
	AdminRoleGroup     = "role-admin"
	PublisherRoleGroup = "role-publisher"
)

type Group struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Precedence  int64  `json:"precedence"`
}

func NewAdminRoleGroup() Group {
	return Group{
		Name:        AdminRoleGroup,
		Description: "The publishing admins",
		Precedence:  1,
	}
}

func NewPublisherRoleGroup() Group {
	return Group{
		Name:        PublisherRoleGroup,
		Description: "The publishers",
		Precedence:  1,
	}
}

func (g *Group) BuildCreateGroupRequest(userPoolId string) *cognitoidentityprovider.CreateGroupInput {
	return &cognitoidentityprovider.CreateGroupInput{
		GroupName:   &g.Name,
		Description: &g.Description,
		Precedence:  &g.Precedence,
		UserPoolId:  &userPoolId,
	}
}

func (g *Group) BuildGetGroupRequest(userPoolId string) *cognitoidentityprovider.GetGroupInput {
	return &cognitoidentityprovider.GetGroupInput{
		GroupName:  &g.Name,
		UserPoolId: &userPoolId,
	}
}
