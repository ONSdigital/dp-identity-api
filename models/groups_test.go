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

func TestNewAdminRoleGroup(t *testing.T) {
	Convey("builds a Group instance with admin group details", t, func() {
		adminGroup := models.NewAdminRoleGroup()

		So(adminGroup.Name, ShouldEqual, models.AdminRoleGroup)
		So(adminGroup.Description, ShouldEqual, models.AdminRoleGroupHumanReadable)
		So(adminGroup.Precedence, ShouldEqual, models.AdminRoleGroupPrecedence)
	})
}

func TestNewPublisherRoleGroup(t *testing.T) {
	Convey("builds a Group instance with publisher group details", t, func() {
		publisherGroup := models.NewPublisherRoleGroup()

		So(publisherGroup.Name, ShouldEqual, models.PublisherRoleGroup)
		So(publisherGroup.Description, ShouldEqual, models.PublisherRoleGroupHumanReadable)
		So(publisherGroup.Precedence, ShouldEqual, models.PublisherRoleGroupPrecedence)
	})
}

func TestGroup_ValidateAddUser(t *testing.T) {
	var ctx = context.Background()

	Convey("returns InvalidUserId error if no user id is submitted", t, func() {
		group := models.Group{
			Name: "test-group",
		}
		userId := ""

		errs := group.ValidateAddRemoveUser(ctx, userId)

		So(errs, ShouldNotBeNil)
		So(len(errs), ShouldEqual, 1)
		castErr := errs[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidUserIdError)
		So(castErr.Description, ShouldEqual, models.MissingUserIdErrorDescription)
	})

	Convey("returns InvalidGroupName error if no group name is set", t, func() {
		group := models.Group{}
		userId := "zzzz-9999"

		errs := group.ValidateAddRemoveUser(ctx, userId)

		So(errs, ShouldNotBeNil)
		So(len(errs), ShouldEqual, 1)
		castErr := errs[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidGroupNameError)
		So(castErr.Description, ShouldEqual, models.MissingGroupNameErrorDescription)
	})

	Convey("returns InvalidUserId and InvalidGroupName errors if no user id submitted and group name set", t, func() {
		group := models.Group{}
		userId := ""

		errs := group.ValidateAddRemoveUser(ctx, userId)

		So(errs, ShouldNotBeNil)
		So(len(errs), ShouldEqual, 2)
		castErr := errs[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidGroupNameError)
		So(castErr.Description, ShouldEqual, models.MissingGroupNameErrorDescription)
		castErr = errs[1].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidUserIdError)
		So(castErr.Description, ShouldEqual, models.MissingUserIdErrorDescription)
	})

	Convey("returns nil if user id is present", t, func() {
		group := models.Group{
			Name: "test-group",
		}
		userId := "zzzz-9999"

		errs := group.ValidateAddRemoveUser(ctx, userId)

		So(errs, ShouldBeNil)
	})
}

func TestGroup_BuildGetGroupRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito GetGroup request body", t, func() {

		group := models.Group{
			Name: "role-admin",
		}

		userPoolId := "euwest-99-aabbcc"

		response := group.BuildGetGroupRequest(userPoolId)

		So(reflect.TypeOf(*response), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.GetGroupInput{}))
		So(*response.UserPoolId, ShouldEqual, userPoolId)
		So(*response.GroupName, ShouldEqual, group.Name)
	})
}

func TestGroup_BuildAddUserToGroupRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito AdminAddUserToGroup request body", t, func() {
		group := models.Group{
			Name: "role-test",
		}
		userPoolId := "euwest-99-aabbcc"
		userId := "zzzz-9999"

		response := group.BuildAddUserToGroupRequest(userPoolId, userId)

		So(reflect.TypeOf(*response), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.AdminAddUserToGroupInput{}))
		So(*response.UserPoolId, ShouldEqual, userPoolId)
		So(*response.GroupName, ShouldEqual, group.Name)
		So(*response.Username, ShouldEqual, userId)
	})
}

func TestGroup_BuildRemoveUserFromGroupRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito AdminRemoveUserFromGroup request body", t, func() {
		group := models.Group{
			Name: "role-test",
		}
		userPoolId := "euwest-99-aabbcc"
		userId := "zzzz-9999"

		response := group.BuildRemoveUserFromGroupRequest(userPoolId, userId)

		So(reflect.TypeOf(*response), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.AdminRemoveUserFromGroupInput{}))
		So(*response.UserPoolId, ShouldEqual, userPoolId)
		So(*response.GroupName, ShouldEqual, group.Name)
		So(*response.Username, ShouldEqual, userId)
	})
}

func TestGroup_BuildListUsersInGroupRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito ListUsersInGroup request body", t, func() {
		group := models.Group{
			Name: "role-test",
		}
		userPoolId := "euwest-99-aabbcc"

		response := group.BuildListUsersInGroupRequest(userPoolId)

		So(reflect.TypeOf(*response), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.ListUsersInGroupInput{}))
		So(*response.UserPoolId, ShouldEqual, userPoolId)
		So(*response.GroupName, ShouldEqual, group.Name)
	})
}

func TestGroup_MapMembers(t *testing.T) {
	Convey("adds the returned users to the members attribute", t, func() {
		cognitoResponse := cognitoidentityprovider.ListUsersOutput{
			Users: []*cognitoidentityprovider.UserType{
				{
					UserStatus: aws.String("CONFIRMED"),
					Username:   aws.String("user-1"),
					Enabled:    aws.Bool(true),
				},
				{
					UserStatus: aws.String("CONFIRMED"),
					Username:   aws.String("user-2"),
					Enabled:    aws.Bool(true),
				},
			},
		}
		group := models.Group{}
		group.MapMembers(&cognitoResponse.Users)

		So(len(group.Members), ShouldEqual, len(cognitoResponse.Users))
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

		So(group.Description, ShouldEqual, *response.Description)
		So(group.Name, ShouldEqual, *response.GroupName)
		So(group.Precedence, ShouldEqual, *response.Precedence)
	})
}

func TestGroup_BuildSuccessfulJsonResponse(t *testing.T) {
	Convey("returns a byte array of the response JSON", t, func() {
		ctx := context.Background()
		name, description := "test-group", "a test group"
		precedence := int64(100)
		group := models.Group{
			Name:        name,
			Description: description,
			Precedence:  precedence,
		}

		response, err := group.BuildSuccessfulJsonResponse(ctx)

		So(err, ShouldBeNil)
		So(reflect.TypeOf(response), ShouldEqual, reflect.TypeOf([]byte{}))
		var userJson map[string]interface{}
		err = json.Unmarshal(response, &userJson)
		So(err, ShouldBeNil)
		So(userJson["name"], ShouldEqual, name)
		So(userJson["description"], ShouldEqual, description)
		So(userJson["precedence"], ShouldEqual, precedence)
	})
}

func TestGroup_BuildListUsersInGroupRequestWithNextToken(t *testing.T) {
	Convey("builds a correctly populated Cognito ListUsersInGroup request body without a nextToken", t, func() {
		group := models.Group{
			Name: "role-test",
		}
		userPoolId := "euwest-99-aabbcc"
		nextToken := ""

		response := group.BuildListUsersInGroupRequestWithNextToken(userPoolId, nextToken)

		So(reflect.TypeOf(*response), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.ListUsersInGroupInput{}))
		So(*response.UserPoolId, ShouldEqual, userPoolId)
		So(*response.GroupName, ShouldEqual, group.Name)
		So(response.NextToken, ShouldBeNil)
	})

	Convey("builds a correctly populated Cognito ListUsersInGroup request body with a nextToken", t, func() {
		group := models.Group{
			Name: "role-test",
		}
		userPoolId := "euwest-99-aabbcc"
		nextToken := "abcd"

		response := group.BuildListUsersInGroupRequestWithNextToken(userPoolId, nextToken)

		So(reflect.TypeOf(*response), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.ListUsersInGroupInput{}))
		So(*response.UserPoolId, ShouldEqual, userPoolId)
		So(*response.GroupName, ShouldEqual, group.Name)
		So(*response.NextToken, ShouldEqual, nextToken)
	})
}
func TestGroup_BuildListGroupsRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito ListGroups request body", t, func() {
		group := models.ListUserGroupType{}
		userPoolId := "euwest-99-aabbcc"
		nextToken := "Next-Token"

		response := group.BuildListGroupsRequest(userPoolId, nextToken)

		So(reflect.TypeOf(*response), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.ListGroupsInput{}))
		So(*response.UserPoolId, ShouldEqual, userPoolId)
		So(*response.NextToken, ShouldEqual, nextToken)
	})

	Convey("builds a correctly populated Cognito ListGroups request body", t, func() {
		group := models.ListUserGroupType{}
		userPoolId := "euwest-99-aabbcc"
		nextToken := ""

		response := group.BuildListGroupsRequest(userPoolId, nextToken)

		So(reflect.TypeOf(*response), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.ListGroupsInput{}))
		So(*response.UserPoolId, ShouldEqual, userPoolId)
		So(response.NextToken, ShouldBeNil)
	})

	Convey("builds a nill Cognito ListGroups request body", t, func() {
		group := models.ListUserGroupType{}
		userPoolId := "euwest-99-aabbcc"
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
		var groupsJson map[string]interface{}
		err = json.Unmarshal(response, &groupsJson)
		So(err, ShouldBeNil)
		So(groupsJson["next_token"], ShouldBeEmpty)
		So(groupsJson["count"], ShouldEqual, 2)
		So(groupsJson["groups"], ShouldNotBeNil)
		jsonGroups := groupsJson["groups"].([]interface{})
		So(len(jsonGroups), ShouldEqual, 2)
		for _, testgroup := range jsonGroups {
			jsonGroup := testgroup.(map[string]interface{})
			So(jsonGroup["description"], ShouldEqual, description)
			So(jsonGroup["precedence"], ShouldEqual, precedence)
			So(jsonGroup["group_name"], ShouldEqual, name)
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

func TestGroup_ValidateCreateGroupRequest(t *testing.T) {
	var(
		ctx = context.Background()
		name = "This^& is a £Tes\t GRoup n%$ame"
		nameWithRole = "role_This^& is a £Tes\t GRoup n%$ame"
		precedence = int64(100)
		lowPrecedence = int64(1)
	)

	Convey("No errors generated", t, func() {
		createGroupTests := []struct {
			description string
			createGroup models.CreateGroup
			expectedResponse map[string]interface{}
			expectedErrors []string
		}{
			{
				"No errors generated",
				models.CreateGroup{
					Description: &name,
					Precedence:  &precedence,
				},
				map[string]interface{}{
					"name": "thisisatestgroupname",
					"precedence":  22,
				},
				nil,
			},
			{
				"Invalid group name error generated",
				models.CreateGroup{
					Precedence:  &precedence,
				},
				nil,
				[]string{
					models.InvalidGroupName,
				},
			},
			{
				"Invalid group pattern error generated",
				models.CreateGroup{
					Description: &nameWithRole,
					Precedence:  &precedence,
				},
				nil,
				[]string{
					models.InvalidGroupName,
				},
			},
			{
				"Invalid group precedence error generated",
				models.CreateGroup{
					Description: &name,
				},
				nil,
				[]string{
					models.InvalidGroupPrecedence,
				},
			},
			{
				"Group precedence incorrect error generated",
				models.CreateGroup{
					Description: &name,
					Precedence: &lowPrecedence,
				},
				nil,
				[]string{
					models.InvalidGroupPrecedence,
				},
			},
			{
				"No group name and precedence in request body error generated",
				models.CreateGroup{},
				nil,
				[]string{
					models.InvalidGroupName,
					models.InvalidGroupPrecedence,
				},
			},
		}
		
		for _, tt := range createGroupTests {
			Convey(tt.description, func() {
				validationErrs := tt.createGroup.ValidateCreateGroupRequest(ctx)
				if tt.expectedErrors != nil {
					for i, err := range tt.expectedErrors {
						So(len(validationErrs), ShouldEqual, len(tt.expectedErrors))
						So(validationErrs[i].Error(), ShouldEqual, err)
					}
				} else {
					So(validationErrs, ShouldEqual, nil)
				}
			})
		}
	})
}

func TestGroup_BuildCreateGroupRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito CreateGroup request body", t, func() {
		var(
			name = "This^& is a £Tes\t GRoup n%$ame"
			precedence = int64(100)
			groupName = "123e4567-e89b-12d3-a456-426614174000"
		)

		group := models.CreateGroup{
			Description: &name,
			Precedence:  &precedence,
			GroupName: groupName,
		}

		userPoolId := "euwest-99-aabbcc"

		response := group.BuildCreateGroupInput(&userPoolId)

		So(reflect.TypeOf(*response), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.CreateGroupInput{}))
		So(*response.UserPoolId, ShouldEqual, userPoolId)
		So(*response.GroupName, ShouldEqual, group.GroupName)
		So(*response.Description, ShouldEqual, *group.Description)
		So(*response.Precedence, ShouldEqual, *group.Precedence)
	})
}

func TestGroup_BuildCreateGroupSuccessfulJsonResponse(t *testing.T) {
	Convey("returns a byte array of the response JSON", t, func() {

		var(
			ctx = context.Background()
			name = "This^& is a £Tes\t GRoup n%$ame"
			precedence = int64(100)
			groupName = "123e4567-e89b-12d3-a456-426614174000"
		)

		group := models.CreateGroup{
			GroupName:   groupName,
			Description: &name,
			Precedence:  &precedence,
		}

		response, err := group.BuildSuccessfulJsonResponse(ctx)

		So(err, ShouldBeNil)
		So(reflect.TypeOf(response), ShouldEqual, reflect.TypeOf([]byte{}))
		var groupJson map[string]interface{}
		err = json.Unmarshal(response, &groupJson)
		So(err, ShouldBeNil)
		So(groupJson["name"], ShouldEqual, name)
		So(groupJson["precedence"], ShouldEqual, precedence)
	})
}

func TestGroup_CreateGroupCleanGroupDescription(t *testing.T) {
	Convey("return a cleaned group name from description", t, func() {

		var(
			name = "This^& is a £Tes\\t GRoup n%$ame"
			cleanGroupName = "thisisatestgroupname"
			precedence = int64(100)
			groupName = "123e4567-e89b-12d3-a456-426614174000"
		)

		group := models.CreateGroup{
			GroupName:   groupName,
			Description: &name,
			Precedence:  &precedence,
		}

		So(*group.Description, ShouldEqual, name)

		group.CleanGroupDescription()

		So(*group.Description, ShouldEqual, cleanGroupName)
	})
}

func TestGroup_CreateGroupNewSuccessResponse(t *testing.T) {
	Convey("builds correctly populated api response for successful CreateGroup request", t, func() {
		var(
			ctx = context.Background()
			name = "thisisatestgroupname"
			precedence = int64(100)
			groupName = "123e4567-e89b-12d3-a456-426614174000"
		)

		group := models.CreateGroup{
			GroupName:   groupName,
			Description: &name,
			Precedence:  &precedence,
		}

		response, _ := group.BuildSuccessfulJsonResponse(ctx)
		successResponse := group.NewSuccessResponse(response, http.StatusCreated, nil)

		So(reflect.TypeOf(*successResponse), ShouldEqual, reflect.TypeOf(models.SuccessResponse{}))

		createGroupResponse := make(map[string]interface{})
		_ = json.Unmarshal(successResponse.Body, &createGroupResponse)

		So(createGroupResponse["name"].(string), ShouldEqual, name)
		So(int64(createGroupResponse["precedence"].(float64)), ShouldEqual, precedence)

	})
}
