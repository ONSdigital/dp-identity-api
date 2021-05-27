package models_test

import (
	"context"
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

		So(err.Code, ShouldEqual, models.InvalidTokenError)
		So(err.Description, ShouldEqual, models.MissingAuthorizationTokenDescription)
	})

	Convey("returns InvalidToken error if auth header string is set as empty string", t, func() {
		accessToken := models.AccessToken{AuthHeader: ""}

		err := accessToken.Validate(ctx)

		So(err.Code, ShouldEqual, models.InvalidTokenError)
		So(err.Description, ShouldEqual, models.MissingAuthorizationTokenDescription)
	})

	Convey("returns InvalidToken error if auth header string cannot be split", t, func() {
		accessToken := models.AccessToken{AuthHeader: "Beareraaaa.bbbb.cccc"}

		err := accessToken.Validate(ctx)

		So(err.Code, ShouldEqual, models.InvalidTokenError)
		So(err.Description, ShouldEqual, models.MalformedAuthorizationTokenDescription)
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
	ctx := context.Background()

	Convey("successfully parse valid token", t, func() {
		testEmailAddress := "test@ons.gov.uk"

		idTokenString := mock.GenerateMockIDToken(testEmailAddress)
		idToken := models.IdToken{}

		err := idToken.ParseWithoutValidating(ctx, idTokenString)

		So(err, ShouldBeNil)
		So(idToken.Claims.Email, ShouldEqual, testEmailAddress)
	})

	Convey("error returned when passed an invalid token", t, func() {
		idTokenString := "aaaa.bbbb.cccc"
		idToken := models.IdToken{}

		err := idToken.ParseWithoutValidating(ctx, idTokenString)
		So(err, ShouldNotBeNil)
	})
}

func TestIdToken_Validate(t *testing.T) {
	var ctx = context.Background()

	Convey("returns an InvalidToken error if no token string is set", t, func() {
		idToken := models.IdToken{}

		err := idToken.Validate(ctx)

		So(err, ShouldNotBeNil)
		So(err.Code, ShouldEqual, models.InvalidTokenError)
		So(err.Description, ShouldEqual, models.MissingIDTokenDescription)
	})

	Convey("returns an InvalidToken error if token string is set as empty string", t, func() {
		idToken := models.IdToken{TokenString: ""}

		err := idToken.Validate(ctx)

		So(err, ShouldNotBeNil)
		So(err.Code, ShouldEqual, models.InvalidTokenError)
		So(err.Description, ShouldEqual, models.MissingIDTokenDescription)
	})

	Convey("returns an InvalidToken error if token string is not parsable", t, func() {
		idToken := models.IdToken{TokenString: "aaaa.bbbb.cccc"}

		err := idToken.Validate(ctx)

		So(err, ShouldNotBeNil)
		So(err.Code, ShouldEqual, models.InvalidTokenError)
		So(err.Description, ShouldEqual, models.MalformedIDTokenDescription)
	})

	Convey("does not return any errors and sets claims if token string valid", t, func() {
		testEmailAddress := "test@ons.gov.uk"

		idTokenString := mock.GenerateMockIDToken(testEmailAddress)
		idToken := models.IdToken{TokenString: idTokenString}

		err := idToken.Validate(ctx)

		So(err, ShouldBeNil)
		So(idToken.Claims.Email, ShouldEqual, testEmailAddress)
	})
}

func TestRefreshToken_Validate(t *testing.T) {
	var ctx = context.Background()

	Convey("returns an InvalidToken error if no token string is set", t, func() {
		refreshToken := models.RefreshToken{}

		err := refreshToken.Validate(ctx)

		So(err, ShouldNotBeNil)
		So(err.Code, ShouldEqual, models.InvalidTokenError)
		So(err.Description, ShouldEqual, models.MissingRefreshTokenDescription)
	})

	Convey("returns an InvalidToken error if token string is set as empty string", t, func() {
		refreshToken := models.RefreshToken{TokenString: ""}

		err := refreshToken.Validate(ctx)

		So(err, ShouldNotBeNil)
		So(err.Code, ShouldEqual, models.InvalidTokenError)
		So(err.Description, ShouldEqual, models.MissingRefreshTokenDescription)
	})

	Convey("does not return any errors token string is set", t, func() {
		refreshToken := models.RefreshToken{TokenString: "aaaa.bbbb.cccc.dddd.eeee"}

		err := refreshToken.Validate(ctx)

		So(err, ShouldBeNil)
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
