package models_test

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ONSdigital/dp-identity-api/v2/cognito/mock"
	"github.com/ONSdigital/dp-identity-api/v2/models"
	"github.com/ONSdigital/dp-identity-api/v2/utilities"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

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
		idToken := models.IDToken{}

		err := idToken.ParseWithoutValidating(ctx, idTokenString)

		So(err, ShouldBeNil)
		So(idToken.Claims.Email, ShouldEqual, testEmailAddress)
	})

	Convey("error returned when passed an invalid token", t, func() {
		idTokenString := "aaaa.bbbb.cccc"
		idToken := models.IDToken{}

		err := idToken.ParseWithoutValidating(ctx, idTokenString)
		So(err, ShouldNotBeNil)
	})
}

func TestIdToken_Validate(t *testing.T) {
	var ctx = context.Background()

	Convey("returns an InvalidToken error if no token string is set", t, func() {
		idToken := models.IDToken{}

		err := idToken.Validate(ctx)

		So(err, ShouldNotBeNil)
		So(err.Code, ShouldEqual, models.InvalidTokenError)
		So(err.Description, ShouldEqual, models.MissingIDTokenDescription)
	})

	Convey("returns an InvalidToken error if token string is set as empty string", t, func() {
		idToken := models.IDToken{TokenString: ""}

		err := idToken.Validate(ctx)

		So(err, ShouldNotBeNil)
		So(err.Code, ShouldEqual, models.InvalidTokenError)
		So(err.Description, ShouldEqual, models.MissingIDTokenDescription)
	})

	Convey("returns an InvalidToken error if token string is not parsable", t, func() {
		idToken := models.IDToken{TokenString: "aaaa.bbbb.cccc"}

		err := idToken.Validate(ctx)

		So(err, ShouldNotBeNil)
		So(err.Code, ShouldEqual, models.InvalidTokenError)
		So(err.Description, ShouldEqual, models.MalformedIDTokenDescription)
	})

	Convey("does not return any errors and sets claims if token string valid", t, func() {
		testEmailAddress := "test@ons.gov.uk"

		idTokenString := mock.GenerateMockIDToken(testEmailAddress)
		idToken := models.IDToken{TokenString: idTokenString}

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
		var clientID, clientSecret, username, refreshTokenString = "abcdefg12345", "hijklmnop67890", "onsTestUser", "zzzz.yyyy.xxxx.wwww.vvvv"
		refreshToken := models.RefreshToken{TokenString: refreshTokenString}

		initiateAuthInput := refreshToken.GenerateRefreshRequest(clientSecret, username, clientID)

		expectedAuthFlow := "REFRESH_TOKEN_AUTH"
		expectedSecretHash := utilities.ComputeSecretHash(clientSecret, username, clientID)
		So(*initiateAuthInput.AuthFlow, ShouldEqual, expectedAuthFlow)
		So(*initiateAuthInput.AuthParameters["REFRESH_TOKEN"], ShouldEqual, refreshTokenString)
		So(*initiateAuthInput.AuthParameters["SECRET_HASH"], ShouldEqual, expectedSecretHash)
		So(*initiateAuthInput.ClientId, ShouldEqual, clientID)
	})
}

func TestRefreshToken_BuildSuccessfulJsonResponse(t *testing.T) {
	ctx := context.Background()

	Convey("returns an InternalServerError if the Cognito response does not meet expected format", t, func() {
		refreshToken := models.RefreshToken{}
		result := cognitoidentityprovider.InitiateAuthOutput{}

		response, err := refreshToken.BuildSuccessfulJSONResponse(ctx, &result)

		So(response, ShouldBeNil)
		castErr := err.(*models.Error)
		So(castErr.Code, ShouldEqual, models.InternalError)
		So(castErr.Description, ShouldEqual, models.UnrecognisedCognitoResponseDescription)
	})

	Convey("returns a byte array of the response JSON", t, func() {
		var expirationLength int64 = 300
		refreshToken := models.RefreshToken{}
		result := cognitoidentityprovider.InitiateAuthOutput{
			AuthenticationResult: &cognitoidentityprovider.AuthenticationResultType{
				ExpiresIn: &expirationLength,
			},
		}

		response, err := refreshToken.BuildSuccessfulJSONResponse(ctx, &result)

		So(err, ShouldBeNil)
		So(reflect.TypeOf(response), ShouldEqual, reflect.TypeOf([]byte{}))
		var body map[string]interface{}
		err = json.Unmarshal(response, &body)
		So(err, ShouldBeNil)
		So(body["expirationTime"], ShouldNotBeNil)
	})
}
