package models_test

import (
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
		description, precedence := "The publishers", 1

		adminGroup := models.NewPublisherRoleGroup()

		So(adminGroup.Name, ShouldEqual, models.PublisherRoleGroup)
		So(adminGroup.Description, ShouldEqual, description)
		So(adminGroup.Precedence, ShouldEqual, precedence)
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
