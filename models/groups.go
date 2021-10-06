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
	groupPrecedenceMin    = int64(3)
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
		Description: AdminRoleGroupHumanReadable,
		Precedence:  AdminRoleGroupPrecedence,
	}
}

// Constructor for a new instance of the publisher role group
func NewPublisherRoleGroup() Group {
	return Group{
		Name:        PublisherRoleGroup,
		Description: PublisherRoleGroupHumanReadable,
		Precedence:  PublisherRoleGroupPrecedence,
	}
}

// ValidateAddRemoveUser validates the required fields for adding a user to a group, returns validation errors for anything that fails
func (g *Group) ValidateAddRemoveUser(ctx context.Context, userId string) []error {
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

// BuildRemoveUserFromGroupRequest builds a correctly populated AdminRemoveUserFromGroupInput object
func (g *Group) BuildRemoveUserFromGroupRequest(userPoolId, userId string) *cognitoidentityprovider.AdminRemoveUserFromGroupInput {
	return &cognitoidentityprovider.AdminRemoveUserFromGroupInput{
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

// BuildListUsersInGroupRequest builds a correctly populated ListUsersInGroupInput object with Next Token
func (g *Group) BuildListUsersInGroupRequestWithNextToken(userPoolId string, nextToken string) *cognitoidentityprovider.ListUsersInGroupInput {
	if nextToken == "" {
		return &cognitoidentityprovider.ListUsersInGroupInput{
			GroupName:  &g.Name,
			UserPoolId: &userPoolId,
		}
	}
	return &cognitoidentityprovider.ListUsersInGroupInput{
		GroupName:  &g.Name,
		UserPoolId: &userPoolId,
		NextToken:  &nextToken,
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

type CreateGroup struct {
	Description *string `json:"name"`
	Precedence  *int64  `json:"precedence"`
	GroupName   string
}

func (g *CreateGroup) ValidateCreateGroupRequest(ctx context.Context) []error {
	var validationErrs []error

	if g.Description == nil {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidGroupName, MissingGroupName))
	} else if m, _ := regexp.MatchString("(?i)^role_.*", *g.Description); m {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidGroupName, IncorrectPatternInGroupName))
	}
	if g.Precedence == nil {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidGroupPrecedence, MissingGroupPrecedence))
	} else if *g.Precedence < groupPrecedenceMin {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidGroupPrecedence, GroupPrecedenceIncorrect))
	}

	return validationErrs
}

func (c *CreateGroup) BuildCreateGroupInput(userPoolId *string) *cognitoidentityprovider.CreateGroupInput {
	return &cognitoidentityprovider.CreateGroupInput{
		Description: c.Description,
		GroupName:   &c.GroupName,
		Precedence:  c.Precedence,
		UserPoolId:  userPoolId,
	}
}

func (c *CreateGroup) BuildSuccessfulJsonResponse(ctx context.Context) ([]byte, error) {
	jsonResponse, err := json.Marshal(c)
	if err != nil {
		e := NewError(ctx, err, JSONMarshalError, ErrorMarshalFailedDescription)
		return nil, e
	}
	return jsonResponse, nil
}

// CleanGroupDescription - New group name to be created from group description
//                        (minus special characters, trimmed and to lowercase)
func (c *CreateGroup) CleanGroupDescription() {
	// strip special chars out of group description string and trim
	// special chars groupset => []£\s^\\$*.]}()?"!@#%&/,><':;|_~-]
	regExp := regexp.MustCompile("[" + groupNameSpecialChars + "]+")
	*c.Description = strings.TrimSpace(
		strings.ToLower(
			regExp.ReplaceAllString(*c.Description, ""),
		),
	)
}

// NewSuccessResponse - returns a custom response where group description is returned as group name
func (c *CreateGroup) NewSuccessResponse(jsonBody []byte, statusCode int, headers map[string]string) *SuccessResponse {
	// unmarshall response and transform: API_Req:name -> Cognito:Description -> API_Resp:name
	var cg = CreateGroup{}
	_ = json.Unmarshal(jsonBody, &cg)

	jsonResponse, _ := json.Marshal(
		map[string]interface{}{
			"name": cg.Description,
			"precedence": cg.Precedence,
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
