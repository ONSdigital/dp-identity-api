package models_test

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/apierrorsdeprecated"
	"github.com/ONSdigital/dp-identity-api/cognito/mock"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/ONSdigital/dp-identity-api/utilities"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAccessToken_Validate(t *testing.T) {
	var ctx = context.Background()

	Convey("returns InvalidToken error if no auth header string is set", t, func() {
		accessToken := models.AccessToken{}

		err := accessToken.Validate(ctx)

		castErr, ok := err.(*models.Error)

		So(ok, ShouldEqual, true)
		So(castErr.Code, ShouldEqual, models.InvalidTokenError)
		So(castErr.Description, ShouldEqual, models.MissingAuthorizationTokenDescription)
	})

	Convey("returns InvalidToken error if auth header string is set as empty string", t, func() {
		accessToken := models.AccessToken{AuthHeader: ""}

		err := accessToken.Validate(ctx)

		castErr, ok := err.(*models.Error)

		So(ok, ShouldEqual, true)
		So(castErr.Code, ShouldEqual, models.InvalidTokenError)
		So(castErr.Description, ShouldEqual, models.MissingAuthorizationTokenDescription)
	})

	Convey("returns InvalidToken error if auth header string cannot be split", t, func() {
		accessToken := models.AccessToken{AuthHeader: "Beareraaaa.bbbb.cccc"}

		err := accessToken.Validate(ctx)

		castErr, ok := err.(*models.Error)

		So(ok, ShouldEqual, true)
		So(castErr.Code, ShouldEqual, models.InvalidTokenError)
		So(castErr.Description, ShouldEqual, models.MalformedAuthorizationTokenDescription)
	})

	Convey("returns nil if auth header string is valid", t, func() {
		accessToken := models.AccessToken{AuthHeader: "Bearer aaaa.bbbb.cccc"}

		err := accessToken.Validate(ctx)

		So(err, ShouldBeNil)
	})
}

func TestAccessToken_GenerateSignOutRequest(t *testing.T) {
	Convey("returns a sign out request input with the access token set", t, func() {
		accessTokenString := "aaaa.bbbb.cccc"
		accessToken := models.AccessToken{TokenString: accessTokenString}

		requestInput := accessToken.GenerateSignOutRequest()

		So(*requestInput.AccessToken, ShouldEqual, accessTokenString)
	})
}

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

func TestIdToken_Validate(t *testing.T) {
	var ctx = context.Background()

	Convey("adds a missing token error if no token string is set", t, func() {
		var errorList []apierrorsdeprecated.Error
		idToken := models.IdToken{}

		errorList = idToken.Validate(ctx, errorList)

		So(len(errorList), ShouldEqual, 1)
		So(errorList[0].Description, ShouldEqual, apierrorsdeprecated.MissingIDTokenMessage)
	})

	Convey("adds a missing token error if token string is set as empty string", t, func() {
		var errorList []apierrorsdeprecated.Error
		idToken := models.IdToken{TokenString: ""}

		errorList = idToken.Validate(ctx, errorList)

		So(len(errorList), ShouldEqual, 1)
		So(errorList[0].Description, ShouldEqual, apierrorsdeprecated.MissingIDTokenMessage)
	})

	Convey("adds a malformed token error if token string is not parsable", t, func() {
		var errorList []apierrorsdeprecated.Error
		idToken := models.IdToken{TokenString: "aaaa.bbbb.cccc"}

		errorList = idToken.Validate(ctx, errorList)

		So(len(errorList), ShouldEqual, 1)
		So(errorList[0].Description, ShouldEqual, apierrorsdeprecated.MalformedIDTokenMessage)
	})

	Convey("does not add any errors and sets claims if token string valid", t, func() {
		var errorList []apierrorsdeprecated.Error
		testEmailAddress := "test@ons.gov.uk"

		idTokenString := mock.GenerateMockIDToken(testEmailAddress)
		idToken := models.IdToken{TokenString: idTokenString}

		errorList = idToken.Validate(ctx, errorList)

		So(len(errorList), ShouldEqual, 0)
		So(idToken.Claims.Email, ShouldEqual, testEmailAddress)
	})
}

func TestRefreshToken_Validate(t *testing.T) {
	var ctx = context.Background()

	Convey("adds a missing token error if no token string is set", t, func() {
		var errorList []apierrorsdeprecated.Error
		refreshToken := models.RefreshToken{}

		errorList = refreshToken.Validate(ctx, errorList)

		So(len(errorList), ShouldEqual, 1)
		So(errorList[0].Description, ShouldEqual, apierrorsdeprecated.MissingRefreshTokenMessage)
	})

	Convey("adds a missing token error if token string is set as empty string", t, func() {
		var errorList []apierrorsdeprecated.Error
		refreshToken := models.RefreshToken{TokenString: ""}

		errorList = refreshToken.Validate(ctx, errorList)

		So(len(errorList), ShouldEqual, 1)
		So(errorList[0].Description, ShouldEqual, apierrorsdeprecated.MissingRefreshTokenMessage)
	})

	Convey("does not add any errors token string is set", t, func() {
		var errorList []apierrorsdeprecated.Error
		refreshToken := models.RefreshToken{TokenString: "aaaa.bbbb.cccc.dddd.eeee"}

		errorList = refreshToken.Validate(ctx, errorList)

		So(len(errorList), ShouldEqual, 0)
	})
}

func TestRefreshToken_GenerateRefreshRequest(t *testing.T) {
	Convey("returns a filled InitiateAuthInput object", t, func() {
		var clientId, clientSecret, username, refreshTokenString string = "abcdefg12345", "hijklmnop67890", "onsTestUser", "zzzz.yyyy.xxxx.wwww.vvvv"
		refreshToken := models.RefreshToken{TokenString: refreshTokenString}

		initiateAuthInput := refreshToken.GenerateRefreshRequest(clientSecret, username, clientId)

		expectedAuthFlow := "REFRESH_TOKEN_AUTH"
		expectedSecretHash := utilities.ComputeSecretHash(clientSecret, username, clientId)
		So(*initiateAuthInput.AuthFlow, ShouldEqual, expectedAuthFlow)
		So(*initiateAuthInput.AuthParameters["REFRESH_TOKEN"], ShouldEqual, refreshTokenString)
		So(*initiateAuthInput.AuthParameters["SECRET_HASH"], ShouldEqual, expectedSecretHash)
		So(*initiateAuthInput.ClientId, ShouldEqual, clientId)
	})
}
