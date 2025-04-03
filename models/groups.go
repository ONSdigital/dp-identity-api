package models

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
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
	groupPrecedenceMin    = int32(10)
	groupPrecedenceMax    = int32(100)
)

// ListGroupUsersType list of groups and the membership for user report group-report
type ListGroupUsersType struct {
	GroupName string `type:"string" json:"group"`
	UserEmail string `type:"string" json:"user"`
}

// Group is a type for the identity API representation of a group's details
type Group struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Precedence int32     `json:"precedence"`
	Created    time.Time `json:"created"`
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
func (g *Group) ValidateAddRemoveUser(ctx context.Context, userID string) []error {
	var validationErrs []error
	if g.ID == "" {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidGroupIDError, MissingGroupIDErrorDescription))
	}

	if userID == "" {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidUserIDError, MissingUserIDErrorDescription))
	}
	return validationErrs
}

// BuildCreateGroupRequest builds a correctly populated CreateGroupInput object using the Groups values
func (g *Group) BuildCreateGroupRequest(userPoolID string) *cognitoidentityprovider.CreateGroupInput {
	return &cognitoidentityprovider.CreateGroupInput{
		GroupName:   &g.ID,
		Description: &g.Name,
		Precedence:  &g.Precedence,
		UserPoolId:  &userPoolID,
	}
}

// BuildGetGroupRequest builds a correctly populated GetGroupInput object using the Groups values
func (g *Group) BuildGetGroupRequest(userPoolID string) *cognitoidentityprovider.GetGroupInput {
	return &cognitoidentityprovider.GetGroupInput{
		GroupName:  &g.ID,
		UserPoolId: &userPoolID,
	}
}

// BuildDeleteGroupRequest builds a correctly populated DeleteGroupInput object using the Groups values
func (g *Group) BuildDeleteGroupRequest(userPoolID string) *cognitoidentityprovider.DeleteGroupInput {
	return &cognitoidentityprovider.DeleteGroupInput{
		GroupName:  &g.ID,
		UserPoolId: &userPoolID,
	}
}

// BuildAddUserToGroupRequest builds a correctly populated AdminAddUserToGroupInput object
func (g *Group) BuildAddUserToGroupRequest(userPoolID, userID string) *cognitoidentityprovider.AdminAddUserToGroupInput {
	return &cognitoidentityprovider.AdminAddUserToGroupInput{
		GroupName:  &g.ID,
		UserPoolId: &userPoolID,
		Username:   &userID,
	}
}

// BuildRemoveUserFromGroupRequest builds a correctly populated AdminRemoveUserFromGroupInput object
func (g *Group) BuildRemoveUserFromGroupRequest(userPoolID, userID string) *cognitoidentityprovider.AdminRemoveUserFromGroupInput {
	return &cognitoidentityprovider.AdminRemoveUserFromGroupInput{
		GroupName:  &g.ID,
		UserPoolId: &userPoolID,
		Username:   &userID,
	}
}

// BuildListUsersInGroupRequest builds a correctly populated ListUsersInGroupInput object
func (g *Group) BuildListUsersInGroupRequest(userPoolID string) *cognitoidentityprovider.ListUsersInGroupInput {
	return &cognitoidentityprovider.ListUsersInGroupInput{
		GroupName:  &g.ID,
		UserPoolId: &userPoolID,
	}
}

// BuildListUsersInGroupRequestWithNextToken builds a correctly populated ListUsersInGroupInput object with Next Token
func (g *Group) BuildListUsersInGroupRequestWithNextToken(userPoolID, nextToken string) *cognitoidentityprovider.ListUsersInGroupInput {
	if nextToken == "" {
		return &cognitoidentityprovider.ListUsersInGroupInput{
			GroupName:  &g.ID,
			UserPoolId: &userPoolID,
		}
	}
	return &cognitoidentityprovider.ListUsersInGroupInput{
		GroupName:  &g.ID,
		UserPoolId: &userPoolID,
		NextToken:  &nextToken,
	}
}

// MapCognitoDetails maps the group details returned from GetGroup requests
func (g *Group) MapCognitoDetails(groupDetails types.GroupType) {
	g.ID = *groupDetails.GroupName
	g.Precedence = *groupDetails.Precedence
	g.Name = *groupDetails.Description
	g.Created = *groupDetails.CreationDate
}

// BuildSuccessfulJSONResponse builds the Group response json for client responses
func (g *Group) BuildSuccessfulJSONResponse(ctx context.Context) ([]byte, error) {
	jsonResponse, err := json.Marshal(g)
	if err != nil {
		return nil, NewError(ctx, err, JSONMarshalError, ErrorMarshalFailedDescription)
	}
	return jsonResponse, nil
}

type CreateUpdateGroup struct {
	Name       *string `json:"name"`
	Precedence *int32  `json:"precedence"`
	ID         *string `json:"id"`
	GroupsList *cognitoidentityprovider.ListGroupsOutput
}

// ValidateCreateUpdateGroupRequest validate the create group request
func (g *CreateUpdateGroup) ValidateCreateUpdateGroupRequest(ctx context.Context, isCreate bool) []error {
	var validationErrs []error

	if g.Name == nil {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidGroupName, MissingGroupName))
	} else if m, _ := regexp.MatchString("(?i)^role-.*", *g.Name); m {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidGroupName, IncorrectPatternInGroupName))
	} else if g.GroupsList != nil {
		// Ensure group name in description doesn't already exist - creation only
		// g.GroupsList not set on updates
		for _, group := range g.GroupsList.Groups {
			if group.Description != nil && CleanString(*group.Description) == CleanString(*g.Name) {
				validationErrs = append(validationErrs, NewValidationError(ctx, GroupExistsError, GroupAlreadyExistsDescription))
				break
			}
		}
	}
	if g.Precedence == nil {
		if isCreate {
			validationErrs = append(validationErrs, NewValidationError(ctx, InvalidGroupPrecedence, MissingGroupPrecedence))
		}
	} else if *g.Precedence < groupPrecedenceMin || *g.Precedence > groupPrecedenceMax {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidGroupPrecedence, GroupPrecedenceIncorrect))
	}

	return validationErrs
}

func (g *CreateUpdateGroup) BuildCreateGroupInput(userPoolID *string) *cognitoidentityprovider.CreateGroupInput {
	return &cognitoidentityprovider.CreateGroupInput{
		Description: g.Name,
		GroupName:   g.ID,
		Precedence:  g.Precedence,
		UserPoolId:  userPoolID,
	}
}

// BuildUpdateGroupInput builds a correctly populated UpdateGroupInput object using Groups values
func (g *CreateUpdateGroup) BuildUpdateGroupInput(userPoolID string) *cognitoidentityprovider.UpdateGroupInput {
	return &cognitoidentityprovider.UpdateGroupInput{
		GroupName:   g.ID,
		Description: g.Name,
		UserPoolId:  &userPoolID,
	}
}

func (g *CreateUpdateGroup) BuildSuccessfulJSONResponse(ctx context.Context) ([]byte, error) {
	jsonResponse, err := json.Marshal(g)
	if err != nil {
		e := NewError(ctx, err, JSONMarshalError, ErrorMarshalFailedDescription)
		return nil, e
	}
	return jsonResponse, nil
}

// NewSuccessResponse - returns a custom response where group description is returned as group name
func (g *CreateUpdateGroup) NewSuccessResponse(jsonBody []byte, statusCode int, headers map[string]string) *SuccessResponse {
	// unmarshall response and transform: API_Req:name -> Cognito:Description -> API_Resp:name
	var cg = CreateUpdateGroup{}
	_ = json.Unmarshal(jsonBody, &cg)

	jsonResponse, _ := json.Marshal(
		map[string]interface{}{
			"name":       cg.Name,
			"precedence": cg.Precedence,
			"id":         cg.ID,
		},
	)

	return &SuccessResponse{
		Body:    jsonResponse,
		Status:  statusCode,
		Headers: headers,
	}
}

// BuildListGroupsSuccessfulJSONResponse
// formats the output to comply with current standards and to json , adds the count of groups returned and
func (p *ListUserGroups) BuildListGroupsSuccessfulJSONResponse(ctx context.Context, result *cognitoidentityprovider.ListGroupsOutput) ([]byte, error) {
	if result == nil {
		return nil, NewValidationError(ctx, InternalError, UnrecognisedCognitoResponseDescription)
	}

	for _, tmpGroup := range result.Groups {
		newGroup := ListUserGroupType{
			CreationDate:     tmpGroup.CreationDate,
			Name:             tmpGroup.Description,
			ID:               tmpGroup.GroupName,
			LastModifiedDate: tmpGroup.LastModifiedDate,
			Precedence:       tmpGroup.Precedence,
			RoleArn:          tmpGroup.RoleArn,
			UserPoolID:       tmpGroup.UserPoolId,
		}

		p.Groups = append(p.Groups, &newGroup)
	}

	p.NextToken = result.NextToken
	p.Count = len(result.Groups)

	jsonResponse, err := json.Marshal(p)
	if err != nil {
		return nil, NewError(ctx, err, JSONMarshalError, ErrorMarshalFailedDescription)
	}
	return jsonResponse, nil
}

// BuildListGroupsRequest build the require input for cognito query to obtain the groups for given user
func (g *ListUserGroupType) BuildListGroupsRequest(userPoolID, nextToken string) *cognitoidentityprovider.ListGroupsInput {
	if nextToken != "" {
		return &cognitoidentityprovider.ListGroupsInput{
			UserPoolId: &userPoolID,
			NextToken:  &nextToken,
		}
	}

	return &cognitoidentityprovider.ListGroupsInput{
		UserPoolId: &userPoolID}
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
