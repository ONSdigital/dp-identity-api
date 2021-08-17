package models_test

import (
	"testing"

	"github.com/ONSdigital/dp-identity-api/cognito/mock"
	"github.com/ONSdigital/dp-identity-api/models"
	. "github.com/smartystreets/goconvey/convey"
)


func TestBuildSignOutUserRequest(t *testing.T) {
	Convey("builds a signout request array of data for AdminUserGlobalSignout", t, func() {
		userPoolId := "eu-test-11_hdsahj9hjxsZ"
		usersList := models.UsersList{}
		users := mock.BulkGenerateUsers(5, nil)
		usersList.MapCognitoUsers(&users.Users)
		g := models.GlobalSignOut{}
		userSignOutRequestData := g.BuildSignOutUserRequest(&usersList.Users, &userPoolId)
		So(len(userSignOutRequestData), ShouldEqual, 5)
	})
}
