package models_test

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	. "github.com/smartystreets/goconvey/convey"
)

const userPoolId = "euwest-99-aabbcc"

func TestNewAdminRoleGroup(t *testing.T) {
	Convey("builds a Group instance with admin group details", t, func() {
		adminGroup := models.NewAdminRoleGroup()

		So(adminGroup.ID, ShouldEqual, models.AdminRoleGroup)
		So(adminGroup.Name, ShouldEqual, models.AdminRoleGroupHumanReadable)
		So(adminGroup.Precedence, ShouldEqual, models.AdminRoleGroupPrecedence)
	})
}

func TestNewPublisherRoleGroup(t *testing.T) {
	Convey("builds a Group instance with publisher group details", t, func() {
		publisherGroup := models.NewPublisherRoleGroup()

		So(publisherGroup.ID, ShouldEqual, models.PublisherRoleGroup)
		So(publisherGroup.Name, ShouldEqual, models.PublisherRoleGroupHumanReadable)
		So(publisherGroup.Precedence, ShouldEqual, models.PublisherRoleGroupPrecedence)
	})
}

func TestGroup_ValidateAddUser(t *testing.T) {
	var ctx = context.Background()

	Convey("returns InvalidUserId error if no user id is submitted", t, func() {
		group := models.Group{
			ID: "test-group",
		}
		userId := ""

		errs := group.ValidateAddRemoveUser(ctx, userId)

		So(errs, ShouldNotBeNil)
		So(len(errs), ShouldEqual, 1)
		castErr := errs[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidUserIdError)
		So(castErr.Description, ShouldEqual, models.MissingUserIdErrorDescription)
	})

	Convey("returns InvalidGroupID error if no group ID is set", t, func() {
		group := models.Group{}
		userId := "zzzz-9999"

		errs := group.ValidateAddRemoveUser(ctx, userId)

		So(errs, ShouldNotBeNil)
		So(len(errs), ShouldEqual, 1)
		castErr := errs[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidGroupIDError)
		So(castErr.Description, ShouldEqual, models.MissingGroupIDErrorDescription)
	})

	Convey("returns InvalidUserId and InvalidGroupID errors if no user id submitted and group ID are set", t, func() {
		group := models.Group{}
		userId := ""

		errs := group.ValidateAddRemoveUser(ctx, userId)

		So(errs, ShouldNotBeNil)
		So(len(errs), ShouldEqual, 2)
		castErr := errs[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidGroupIDError)
		So(castErr.Description, ShouldEqual, models.MissingGroupIDErrorDescription)
		castErr = errs[1].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidUserIdError)
		So(castErr.Description, ShouldEqual, models.MissingUserIdErrorDescription)
	})

	Convey("returns nil if user id is present", t, func() {
		group := models.Group{
			ID: "test-group",
		}
		userId := "zzzz-9999"

		errs := group.ValidateAddRemoveUser(ctx, userId)

		So(errs, ShouldBeNil)
	})
}

func TestGroup_BuildGetGroupRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito GetGroup request body", t, func() {
		group := models.Group{
			ID: "role-admin",
		}

		response := group.BuildGetGroupRequest(userPoolId)

		So(reflect.TypeOf(*response), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.GetGroupInput{}))
		So(*response.UserPoolId, ShouldEqual, userPoolId)
		So(*response.GroupName, ShouldEqual, group.ID)
	})
}

func TestGroup_BuildAddUserToGroupRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito AdminAddUserToGroup request body", t, func() {
		group := models.Group{
			ID: "role-test",
		}

		userId := "zzzz-9999"

		response := group.BuildAddUserToGroupRequest(userPoolId, userId)

		So(reflect.TypeOf(*response), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.AdminAddUserToGroupInput{}))
		So(*response.UserPoolId, ShouldEqual, userPoolId)
		So(*response.GroupName, ShouldEqual, group.ID)
		So(*response.Username, ShouldEqual, userId)
	})
}

func TestGroup_BuildRemoveUserFromGroupRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito AdminRemoveUserFromGroup request body", t, func() {
		group := models.Group{
			ID: "role-test",
		}

		userId := "zzzz-9999"

		response := group.BuildRemoveUserFromGroupRequest(userPoolId, userId)

		So(reflect.TypeOf(*response), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.AdminRemoveUserFromGroupInput{}))
		So(*response.UserPoolId, ShouldEqual, userPoolId)
		So(*response.GroupName, ShouldEqual, group.ID)
		So(*response.Username, ShouldEqual, userId)
	})
}

func TestGroup_BuildListUsersInGroupRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito ListUsersInGroup request body", t, func() {
		group := models.Group{
			ID: "role-test",
		}

		response := group.BuildListUsersInGroupRequest(userPoolId)

		So(reflect.TypeOf(*response), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.ListUsersInGroupInput{}))
		So(*response.UserPoolId, ShouldEqual, userPoolId)
		So(*response.GroupName, ShouldEqual, group.ID)
	})
}

func TestGroup_MapCognitoDetails(t *testing.T) {
	Convey("correctly maps values from Cognito GroupType", t, func() {
		group := models.Group{}

		timestamp := time.Now()
		response := &cognitoidentityprovider.GroupType{
			Description:  aws.String("A test group"),
			GroupName:    aws.String("test-group"),
			Precedence:   aws.Int64(1),
			CreationDate: &timestamp,
		}

		group.MapCognitoDetails(response)

		So(group.Name, ShouldEqual, *response.Description)
		So(group.ID, ShouldEqual, *response.GroupName)
		So(group.Precedence, ShouldEqual, *response.Precedence)
	})
}

func TestGroup_BuildSuccessfulJsonResponse(t *testing.T) {
	Convey("returns a byte array of the response JSON", t, func() {
		ctx := context.Background()
		id, name := "test-group", "a test group"
		precedence := int64(100)
		group := models.Group{
			ID:         id,
			Name:       name,
			Precedence: precedence,
		}

		response, err := group.BuildSuccessfulJsonResponse(ctx)

		So(err, ShouldBeNil)
		So(reflect.TypeOf(response), ShouldEqual, reflect.TypeOf([]byte{}))
		var userJSON map[string]interface{}
		err = json.Unmarshal(response, &userJSON)
		So(err, ShouldBeNil)
		So(userJSON["id"], ShouldEqual, id)
		So(userJSON["name"], ShouldEqual, name)
		So(userJSON["precedence"], ShouldEqual, precedence)
	})
}

func TestGroup_BuildListUsersInGroupRequestWithNextToken(t *testing.T) {
	Convey("builds a correctly populated Cognito ListUsersInGroup request body without a nextToken", t, func() {
		group := models.Group{
			ID: "role-test",
		}

		nextToken := ""

		response := group.BuildListUsersInGroupRequestWithNextToken(userPoolId, nextToken)

		So(reflect.TypeOf(*response), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.ListUsersInGroupInput{}))
		So(*response.UserPoolId, ShouldEqual, userPoolId)
		So(*response.GroupName, ShouldEqual, group.ID)
		So(response.NextToken, ShouldBeNil)
	})

	Convey("builds a correctly populated Cognito ListUsersInGroup request body with a nextToken", t, func() {
		group := models.Group{
			ID: "role-test",
		}

		nextToken := "abcd"

		response := group.BuildListUsersInGroupRequestWithNextToken(userPoolId, nextToken)

		So(reflect.TypeOf(*response), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.ListUsersInGroupInput{}))
		So(*response.UserPoolId, ShouldEqual, userPoolId)
		So(*response.GroupName, ShouldEqual, group.ID)
		So(*response.NextToken, ShouldEqual, nextToken)
	})
}
func TestGroup_BuildListGroupsRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito ListGroups request body", t, func() {
		group := models.ListUserGroupType{}
		nextToken := "Next-Token"

		response := group.BuildListGroupsRequest(userPoolId, nextToken)

		So(reflect.TypeOf(*response), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.ListGroupsInput{}))
		So(*response.UserPoolId, ShouldEqual, userPoolId)
		So(*response.NextToken, ShouldEqual, nextToken)
	})

	Convey("builds a correctly populated Cognito ListGroups request body", t, func() {
		group := models.ListUserGroupType{}
		nextToken := ""

		response := group.BuildListGroupsRequest(userPoolId, nextToken)

		So(reflect.TypeOf(*response), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.ListGroupsInput{}))
		So(*response.UserPoolId, ShouldEqual, userPoolId)
		So(response.NextToken, ShouldBeNil)
	})

	Convey("builds a nill Cognito ListGroups request body", t, func() {
		group := models.ListUserGroupType{}
		nextToken := ""

		response := group.BuildListGroupsRequest(userPoolId, nextToken)

		So(reflect.TypeOf(*response), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.ListGroupsInput{}))
		So(*response.UserPoolId, ShouldEqual, userPoolId)
		So(response.NextToken, ShouldBeNil)
	})
}

func TestGroup_BuildListGroupsSuccessfulJsonResponse(t *testing.T) {
	Convey("returns a byte array of the response JSON", t, func() {
		ctx := context.Background()
		name, description := "test-group", "a test group"
		precedence := int64(100)
		group := models.ListUserGroups{}
		results := &cognitoidentityprovider.ListGroupsOutput{
			Groups: []*cognitoidentityprovider.GroupType{
				{
					GroupName:   &name,
					Description: &description,
					Precedence:  &precedence,
				},
				{
					GroupName:   &name,
					Description: &description,
					Precedence:  &precedence}},
			NextToken: new(string),
		}

		response, err := group.BuildListGroupsSuccessfulJsonResponse(ctx, results)

		So(err, ShouldBeNil)
		So(reflect.TypeOf(response), ShouldEqual, reflect.TypeOf([]byte{}))
		var groupsJSON map[string]interface{}
		err = json.Unmarshal(response, &groupsJSON)
		So(err, ShouldBeNil)
		So(groupsJSON["next_token"], ShouldBeEmpty)
		So(groupsJSON["count"], ShouldEqual, 2)
		So(groupsJSON["groups"], ShouldNotBeNil)
		jsonGroups := groupsJSON["groups"].([]interface{})
		So(len(jsonGroups), ShouldEqual, 2)
		for _, testgroup := range jsonGroups {
			jsonGroup := testgroup.(map[string]interface{})
			So(jsonGroup["name"], ShouldEqual, description)
			So(jsonGroup["precedence"], ShouldEqual, precedence)
			So(jsonGroup["id"], ShouldEqual, name)
		}
	})

	Convey("nil result", t, func() {
		ctx := context.Background()
		group := models.ListUserGroups{}
		var results *cognitoidentityprovider.ListGroupsOutput = nil
		response, err := group.BuildListGroupsSuccessfulJsonResponse(ctx, results)
		So(response, ShouldBeNil)
		So(err, ShouldNotBeNil)
	},
	)
}

func TestGroup_ValidateCreateUpdateGroupRequest(t *testing.T) {
	var (
		ctx           = context.Background()
		name          = "This^& is a £Tes\\t GRoup n%$ame"
		nameWithRole  = "role-This^& is a £Tes\t GRoup n%$ame"
		precedence    = int64(100)
		lowPrecedence = int64(1)
		d             = "thisisatestgroupname"
		g             = "123e4567-e89b-12d3-a456-426614174000"
		p             = int64(12)
	)

	Convey("No errors generated", t, func() {
		CreateUpdateGroupTests := []struct {
			description       string
			CreateUpdateGroup models.CreateUpdateGroup
			expectedResponse  map[string]interface{}
			expectedErrors    []string
			isCreate          bool
		}{
			{
				"No errors generated",
				models.CreateUpdateGroup{
					Name:       &name,
					Precedence: &precedence,
				},
				map[string]interface{}{
					"name":       "thisisatestgroupname",
					"precedence": 22,
				},
				nil,
				true,
			},
			{
				"No errors generated updating without precedence",
				models.CreateUpdateGroup{
					Name: &name,
				},
				map[string]interface{}{
					"name": "thisisatestUpdateGroupname",
				},
				nil,
				false,
			},
			{
				"Invalid group name error generated",
				models.CreateUpdateGroup{
					Precedence: &precedence,
				},
				nil,
				[]string{
					models.InvalidGroupName,
				},
				true,
			},
			{
				"Invalid group pattern error generated",
				models.CreateUpdateGroup{
					Name:       &nameWithRole,
					Precedence: &precedence,
				},
				nil,
				[]string{
					models.InvalidGroupName,
				},
				true,
			},
			{
				"Invalid group precedence error generated",
				models.CreateUpdateGroup{
					Name: &name,
				},
				nil,
				[]string{
					models.InvalidGroupPrecedence,
				},
				true,
			},
			{
				"Group precedence incorrect error generated",
				models.CreateUpdateGroup{
					Name:       &name,
					Precedence: &lowPrecedence,
				},
				nil,
				[]string{
					models.InvalidGroupPrecedence,
				},
				true,
			},
			{
				"No group name and precedence in request body error generated",
				models.CreateUpdateGroup{},
				nil,
				[]string{
					models.InvalidGroupName,
					models.InvalidGroupPrecedence,
				},
				true,
			},
			{
				"Group name already exists error generated",
				models.CreateUpdateGroup{
					Name:       &name,
					Precedence: &precedence,
					GroupsList: &cognitoidentityprovider.ListGroupsOutput{
						NextToken: nil,
						Groups: []*cognitoidentityprovider.GroupType{
							{
								Description: &d,
								GroupName:   &g,
								Precedence:  &p,
							},
						},
					},
				},
				nil,
				[]string{
					models.GroupExistsError,
				},
				true,
			},
		}

		for _, tt := range CreateUpdateGroupTests {
			Convey(tt.description, func() {
				validationErrs := tt.CreateUpdateGroup.ValidateCreateUpdateGroupRequest(ctx, tt.isCreate)
				if tt.expectedErrors != nil {
					for i, err := range tt.expectedErrors {
						So(len(validationErrs), ShouldEqual, len(tt.expectedErrors))
						So(validationErrs[i].Error(), ShouldEqual, err)
					}
				} else {
					So(validationErrs, ShouldBeNil)
				}
			})
		}
	})
}

func TestGroup_BuildCreateUpdateGroupRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito CreateUpdateGroup request body", t, func() {
		var (
			name       = "This^& is a £Tes\t GRoup n%$ame"
			precedence = int64(100)
			groupName  = "123e4567-e89b-12d3-a456-426614174000"
		)

		group := models.CreateUpdateGroup{
			Name:       &name,
			Precedence: &precedence,
			ID:         &groupName,
		}

		userPoolIDVar := userPoolId
		response := group.BuildCreateGroupInput(&userPoolIDVar)

		So(reflect.TypeOf(*response), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.CreateGroupInput{}))
		So(*response.UserPoolId, ShouldEqual, userPoolId)
		So(*response.GroupName, ShouldEqual, *group.ID)
		So(*response.Description, ShouldEqual, *group.Name)
		So(*response.Precedence, ShouldEqual, *group.Precedence)
	})
}

func TestGroup_BuildCreateUpdateGroupSuccessfulJsonResponse(t *testing.T) {
	Convey("returns a byte array of the response JSON", t, func() {
		var (
			ctx        = context.Background()
			name       = "This^& is a £Tes\t GRoup n%$ame"
			precedence = int64(100)
			groupName  = "123e4567-e89b-12d3-a456-426614174000"
		)

		group := models.CreateUpdateGroup{
			ID:         &groupName,
			Name:       &name,
			Precedence: &precedence,
		}

		response, err := group.BuildSuccessfulJsonResponse(ctx)

		So(err, ShouldBeNil)
		So(reflect.TypeOf(response), ShouldEqual, reflect.TypeOf([]byte{}))
		var groupJSON map[string]interface{}
		err = json.Unmarshal(response, &groupJSON)
		So(err, ShouldBeNil)
		So(groupJSON["name"], ShouldEqual, name)
		So(groupJSON["precedence"], ShouldEqual, precedence)
	})
}

func TestGroup_CreateUpdateGroupCleanGroupName(t *testing.T) {
	Convey("return a cleaned group name from description", t, func() {
		var (
			name       = "This^& is a £Tes\\t GRoup n%$ame"
			precedence = int64(100)
			groupName  = "123e4567-e89b-12d3-a456-426614174000"
		)

		group := models.CreateUpdateGroup{
			ID:         &groupName,
			Name:       &name,
			Precedence: &precedence,
		}

		So(*group.Name, ShouldEqual, name)
	})
}

func TestGroup_CreateUpdateGroupNewSuccessResponse(t *testing.T) {
	Convey("builds correctly populated api response for successful CreateUpdateGroup request", t, func() {
		var (
			ctx        = context.Background()
			name       = "thisisatestgroupname"
			precedence = int64(100)
			groupName  = "123e4567-e89b-12d3-a456-426614174000"
		)

		group := models.CreateUpdateGroup{
			ID:         &groupName,
			Name:       &name,
			Precedence: &precedence,
		}

		response, _ := group.BuildSuccessfulJsonResponse(ctx)
		successResponse := group.NewSuccessResponse(response, http.StatusCreated, nil)

		So(reflect.TypeOf(*successResponse), ShouldEqual, reflect.TypeOf(models.SuccessResponse{}))

		CreateUpdateGroupResponse := make(map[string]interface{})
		_ = json.Unmarshal(successResponse.Body, &CreateUpdateGroupResponse)

		So(CreateUpdateGroupResponse["name"].(string), ShouldEqual, name)
		So(int64(CreateUpdateGroupResponse["precedence"].(float64)), ShouldEqual, precedence)
	})
}
