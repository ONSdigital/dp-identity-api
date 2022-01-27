package models

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

const (
	AdminRoleGroup                  = "role-admin"
	AdminRoleGroupPrecedence        = 2
	AdminRoleGroupHumanReadable     = "Administrators"
	PublisherRoleGroup              = "role-publisher"
	PublisherRoleGroupPrecedence    = 3
	PublisherRoleGroupHumanReadable = "Publishing Officers"
)

var (
	groupNameSpecialChars = `]£\s^\\\$\*\.\]\[\}\(\)\?\"\!\@\#\%\&\/\,\>\<\'\:\;\|\_\~\-`
	groupPrecedenceMin    = int64(10)
	groupPrecedenceMax    = int64(100)
)

//Group is a type to map for the Cognito GroupType object
type Group struct {
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	Precedence int64        `json:"precedence"`
	Created    time.Time    `json:"created"`
	Members    []UserParams `json:"members"`
}

// NewAdminRoleGroup is a constructor for a new instance of the admin role group
func NewAdminRoleGroup() Group {
	return Group{
		ID:         AdminRoleGroup,
		Name:       AdminRoleGroupHumanReadable,
		Precedence: AdminRoleGroupPrecedence,
	}
}

// NewPublisherRoleGroup is a constructor for a new instance of the publisher role group
func NewPublisherRoleGroup() Group {
	return Group{
		ID:         PublisherRoleGroup,
		Name:       PublisherRoleGroupHumanReadable,
		Precedence: PublisherRoleGroupPrecedence,
	}
}

// ValidateAddRemoveUser validates the required fields for adding a user to a group, returns validation errors for anything that fails
func (g *Group) ValidateAddRemoveUser(ctx context.Context, userId string) []error {
	var validationErrs []error
	if g.ID == "" {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidGroupIDError, MissingGroupIDErrorDescription))
	}

	if userId == "" {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidUserIdError, MissingUserIdErrorDescription))
	}
	return validationErrs
}

// BuildCreateGroupRequest builds a correctly populated CreateGroupInput object using the Groups values
func (g *Group) BuildCreateGroupRequest(userPoolId string) *cognitoidentityprovider.CreateGroupInput {
	return &cognitoidentityprovider.CreateGroupInput{
		GroupName:   &g.ID,
		Description: &g.Name,
		Precedence:  &g.Precedence,
		UserPoolId:  &userPoolId,
	}
}

// BuildGetGroupRequest builds a correctly populated GetGroupInput object using the Groups values
func (g *Group) BuildGetGroupRequest(userPoolId string) *cognitoidentityprovider.GetGroupInput {
	return &cognitoidentityprovider.GetGroupInput{
		GroupName:  &g.ID,
		UserPoolId: &userPoolId,
	}
}

// BuildDeleteGroupRequest builds a correctly populated DeleteGroupInput object using the Groups values
func (g *Group) BuildDeleteGroupRequest(userPoolId string) *cognitoidentityprovider.DeleteGroupInput {
	return &cognitoidentityprovider.DeleteGroupInput{
		GroupName:  &g.ID,
		UserPoolId: &userPoolId,
	}
}

// BuildAddUserToGroupRequest builds a correctly populated AdminAddUserToGroupInput object
func (g *Group) BuildAddUserToGroupRequest(userPoolId, userId string) *cognitoidentityprovider.AdminAddUserToGroupInput {
	return &cognitoidentityprovider.AdminAddUserToGroupInput{
		GroupName:  &g.ID,
		UserPoolId: &userPoolId,
		Username:   &userId,
	}
}

// BuildRemoveUserFromGroupRequest builds a correctly populated AdminRemoveUserFromGroupInput object
func (g *Group) BuildRemoveUserFromGroupRequest(userPoolId, userId string) *cognitoidentityprovider.AdminRemoveUserFromGroupInput {
	return &cognitoidentityprovider.AdminRemoveUserFromGroupInput{
		GroupName:  &g.ID,
		UserPoolId: &userPoolId,
		Username:   &userId,
	}
}

// BuildListUsersInGroupRequest builds a correctly populated ListUsersInGroupInput object
func (g *Group) BuildListUsersInGroupRequest(userPoolId string) *cognitoidentityprovider.ListUsersInGroupInput {
	return &cognitoidentityprovider.ListUsersInGroupInput{
		GroupName:  &g.ID,
		UserPoolId: &userPoolId,
	}
}

// BuildListUsersInGroupRequestWithNextToken builds a correctly populated ListUsersInGroupInput object with Next Token
func (g *Group) BuildListUsersInGroupRequestWithNextToken(userPoolId string, nextToken string) *cognitoidentityprovider.ListUsersInGroupInput {
	if nextToken == "" {
		return &cognitoidentityprovider.ListUsersInGroupInput{
			GroupName:  &g.ID,
			UserPoolId: &userPoolId,
		}
	}
	return &cognitoidentityprovider.ListUsersInGroupInput{
		GroupName:  &g.ID,
		UserPoolId: &userPoolId,
		NextToken:  &nextToken,
	}
}

// MapCognitoDetails maps the group details returned from GetGroup requests
func (g *Group) MapCognitoDetails(groupDetails *cognitoidentityprovider.GroupType) {
	g.ID = *groupDetails.GroupName
	g.Precedence = *groupDetails.Precedence
	g.Name = *groupDetails.Description
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

type CreateUpdateGroup struct {
	Name       *string `json:"name"`
	Precedence *int64  `json:"precedence"`
	ID         *string `json:"id"`
	GroupsList *cognitoidentityprovider.ListGroupsOutput
}

// ValidateCreateUpdateGroupRequest validate the create group request
func (g *CreateUpdateGroup) ValidateCreateUpdateGroupRequest(ctx context.Context) []error {
	var validationErrs []error

	if g.Name == nil {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidGroupName, MissingGroupName))
	} else if m, _ := regexp.MatchString("(?i)^role-.*", *g.Name); m {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidGroupName, IncorrectPatternInGroupName))
	} else {
		//ensure group name in description doesn't already exist - creation only - g.GroupsList not set on updates
		if g.GroupsList != nil {
			for _, group := range g.GroupsList.Groups {
				if group.Description != nil && CleanString(*group.Description) == CleanString(*g.Name) {
					validationErrs = append(validationErrs, NewValidationError(ctx, GroupExistsError, GroupAlreadyExistsDescription))
					break
				}
			}
		}
	}
	if g.Precedence == nil {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidGroupPrecedence, MissingGroupPrecedence))
	} else if *g.Precedence < groupPrecedenceMin || *g.Precedence > groupPrecedenceMax {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidGroupPrecedence, GroupPrecedenceIncorrect))
	}

	return validationErrs
}

func (c *CreateUpdateGroup) BuildCreateGroupInput(userPoolId *string) *cognitoidentityprovider.CreateGroupInput {
	return &cognitoidentityprovider.CreateGroupInput{
		Description: c.Name,
		GroupName:   c.ID,
		Precedence:  c.Precedence,
		UserPoolId:  userPoolId,
	}
}

// BuildUpdateGroupInput builds a correctly populated UpdateGroupInput object using Groups values
func (g *CreateUpdateGroup) BuildUpdateGroupInput(userPoolId string) *cognitoidentityprovider.UpdateGroupInput {
	return &cognitoidentityprovider.UpdateGroupInput{
		GroupName:   g.Name,
		Description: g.Name,
		Precedence:  g.Precedence,
		UserPoolId:  &userPoolId,
	}
}

func (c *CreateUpdateGroup) BuildSuccessfulJsonResponse(ctx context.Context) ([]byte, error) {
	jsonResponse, err := json.Marshal(c)
	if err != nil {
		e := NewError(ctx, err, JSONMarshalError, ErrorMarshalFailedDescription)
		return nil, e
	}
	return jsonResponse, nil
}

// NewSuccessResponse - returns a custom response where group description is returned as group name
func (c *CreateUpdateGroup) NewSuccessResponse(jsonBody []byte, statusCode int, headers map[string]string) *SuccessResponse {
	// unmarshall response and transform: API_Req:name -> Cognito:Description -> API_Resp:name
	var cg = CreateUpdateGroup{}
	_ = json.Unmarshal(jsonBody, &cg)

	jsonResponse, _ := json.Marshal(
		map[string]interface{}{
			"name":       cg.Name,
			"precedence": cg.Precedence,
			"id":  cg.ID,
		},
	)

	return &SuccessResponse{
		Body:    jsonResponse,
		Status:  statusCode,
		Headers: headers,
	}
}

//BuildListGroupsSuccessfulJsonResponse
// formats the output to comply with current standards and to json , adds the count of groups returned and
func (g *ListUserGroups) BuildListGroupsSuccessfulJsonResponse(ctx context.Context, result *cognitoidentityprovider.ListGroupsOutput) ([]byte, error) {

	if result == nil {
		return nil, NewValidationError(ctx, InternalError, UnrecognisedCognitoResponseDescription)
	}

	for _, tmpGroup := range result.Groups {

		newGroup := ListUserGroupType{
			CreationDate:     tmpGroup.CreationDate,
			Description:      tmpGroup.Description,
			GroupName:        tmpGroup.GroupName,
			LastModifiedDate: tmpGroup.LastModifiedDate,
			Precedence:       tmpGroup.Precedence,
			RoleArn:          tmpGroup.RoleArn,
			UserPoolId:       tmpGroup.UserPoolId,
		}

		g.Groups = append(g.Groups, &newGroup)
	}

	g.NextToken = result.NextToken
	g.Count = len(result.Groups)

	jsonResponse, err := json.Marshal(g)
	if err != nil {
		return nil, NewError(ctx, err, JSONMarshalError, ErrorMarshalFailedDescription)
	}
	return jsonResponse, nil
}

// BuildListGroupsRequest build the require input for cognito query to obtain the groups for given user
func (g *ListUserGroupType) BuildListGroupsRequest(userPoolId string, nextToken string) *cognitoidentityprovider.ListGroupsInput {

	if nextToken != "" {
		return &cognitoidentityprovider.ListGroupsInput{
			UserPoolId: &userPoolId,
			NextToken:  &nextToken,
		}
	}

	return &cognitoidentityprovider.ListGroupsInput{
		UserPoolId: &userPoolId}

}

// CleanString - strip special chars out of incoming string and trim
func CleanString(description string) string {
	// strip special chars out of group description string and trim
	// special chars groupset => []£\s^\\$*.]}()?"!@#%&/,><':;|_~-]
	regExp := regexp.MustCompile("[" + groupNameSpecialChars + "]+")
	return strings.TrimSpace(
		strings.ToLower(
			regExp.ReplaceAllString(description, ""),
		),
	)
}
