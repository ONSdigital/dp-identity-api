package models_test

import (
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

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
