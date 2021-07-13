package models_test

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewAdminRoleGroup(t *testing.T) {
	Convey("builds a Group instance with admin group details", t, func() {
		description, precedence := "The publishing admins", 1

		adminGroup := models.NewAdminRoleGroup()

		So(adminGroup.Name, ShouldEqual, models.AdminRoleGroup)
		So(adminGroup.Description, ShouldEqual, description)
		So(adminGroup.Precedence, ShouldEqual, precedence)
	})
}

func TestNewPublisherRoleGroup(t *testing.T) {
	Convey("builds a Group instance with publisher group details", t, func() {
		description, precedence := "The publishers", 2

		adminGroup := models.NewPublisherRoleGroup()

		So(adminGroup.Name, ShouldEqual, models.PublisherRoleGroup)
		So(adminGroup.Description, ShouldEqual, description)
		So(adminGroup.Precedence, ShouldEqual, precedence)
	})
}

func TestGroup_ValidateAddUser(t *testing.T) {
	var ctx = context.Background()

	Convey("returns InvalidUserId error if no user id is submitted", t, func() {
		group := models.Group{
			Name: "test-group",
		}
		userId := ""

		errs := group.ValidateAddUser(ctx, userId)

		So(errs, ShouldNotBeNil)
		So(len(errs), ShouldEqual, 1)
		castErr := errs[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidUserIdError)
		So(castErr.Description, ShouldEqual, models.MissingUserIdErrorDescription)
	})

	Convey("returns InvalidGroupName error if no group name is set", t, func() {
		group := models.Group{}
		userId := "zzzz-9999"

		errs := group.ValidateAddUser(ctx, userId)

		So(errs, ShouldNotBeNil)
		So(len(errs), ShouldEqual, 1)
		castErr := errs[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidGroupNameError)
		So(castErr.Description, ShouldEqual, models.MissingGroupNameErrorDescription)
	})

	Convey("returns InvalidUserId and InvalidGroupName errors if no user id submitted and group name set", t, func() {
		group := models.Group{}
		userId := ""

		errs := group.ValidateAddUser(ctx, userId)

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

		errs := group.ValidateAddUser(ctx, userId)

		So(errs, ShouldBeNil)
	})
}

func TestGroups_BuildCreateGroupRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito CreateGroup request body", t, func() {

		group := models.Group{
			Name:        "role-admin",
			Description: "Test admin role group",
			Precedence:  1,
		}

		userPoolId := "euwest-99-aabbcc"

		response := group.BuildCreateGroupRequest(userPoolId)

		So(reflect.TypeOf(*response), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.CreateGroupInput{}))
		So(*response.UserPoolId, ShouldEqual, userPoolId)
		So(*response.GroupName, ShouldEqual, group.Name)
		So(*response.Description, ShouldEqual, group.Description)
		So(*response.Precedence, ShouldEqual, group.Precedence)
	})
}

func TestGroups_BuildGetGroupRequest(t *testing.T) {
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

func TestGroups_BuildAddUserToGroupRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito GetGroup request body", t, func() {
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
