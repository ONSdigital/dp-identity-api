package models

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"time"
)

const (
	AdminRoleGroup     = "role-admin"
	PublisherRoleGroup = "role-publisher"
)

//Type to map for the Cognito GroupType object
type Group struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Precedence  int64        `json:"precedence"`
	Created     time.Time    `json:"created"`
	Members     []UserParams `json:"members"`
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

// BuildListUsersInGroupRequest builds a correctly populated ListUsersInGroupInput object
func (g *Group) BuildListUsersInGroupRequest(userPoolId string) *cognitoidentityprovider.ListUsersInGroupInput {
	return &cognitoidentityprovider.ListUsersInGroupInput{
		GroupName:  &g.Name,
		UserPoolId: &userPoolId,
	}
}

// MapCognitoDetails maps the group details returned from GetGroup requests
func (g *Group) MapCognitoDetails(groupDetails *cognitoidentityprovider.GroupType) {
	g.Name = *groupDetails.GroupName
	g.Precedence = *groupDetails.Precedence
	g.Description = *groupDetails.Description
	g.Created = *groupDetails.CreationDate
}

// MapMembers maps Cognito user details to the internal UserParams model from ListUserInGroup requests
func (g *Group) MapMembers(membersList *[]*cognitoidentityprovider.UserType) {
	g.Members = []UserParams{}
	for _, member := range *membersList {
		g.Members = append(g.Members, UserParams{}.MapCognitoDetails(member))
	}
}

//BuildSuccessfulJsonResponse builds the Group response json for client responses
func (g *Group) BuildSuccessfulJsonResponse(ctx context.Context) ([]byte, error) {
	jsonResponse, err := json.Marshal(g)
	if err != nil {
		return nil, NewError(ctx, err, JSONMarshalError, ErrorMarshalFailedDescription)
	}
	return jsonResponse, nil
}
