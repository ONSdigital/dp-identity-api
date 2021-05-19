package models_test

import (
	"github.com/ONSdigital/dp-identity-api/cognito/mock"
	"github.com/ONSdigital/dp-identity-api/models"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIdToken_ParseWithoutValidating(t *testing.T) {

	Convey("successfully parse valid token", t, func() {
		testEmailAddress := "test@ons.gov.uk"

		idTokenString := mock.GenerateMockIDToken(testEmailAddress)
		idToken := models.IdToken{}

		err := idToken.ParseWithoutValidating(idTokenString)

		So(err, ShouldBeNil)
		So(idToken.Claims.Email, ShouldEqual, testEmailAddress)
	})

	Convey("error returned when passed an invalid token", t, func() {
		idTokenString := "aaaa.bbbb.cccc"
		idToken := models.IdToken{}

		err := idToken.ParseWithoutValidating(idTokenString)
		So(err, ShouldNotBeNil)
	})
}
