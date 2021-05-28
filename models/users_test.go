package models_test

import (
	"context"
	"encoding/json"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	"testing"
)

func TestUserSignIn_ValidateCredentials(t *testing.T) {
	ctx := context.Background()

	Convey("no errors are returned if a valid email address and password are provided", t, func() {
		signIn := models.UserSignIn{
			Password: "password",
			Email:    "email.email@ons.gov.uk",
		}

		validationErrors := signIn.ValidateCredentials(ctx)
		So(validationErrors, ShouldBeNil)
	})

	Convey("an InvalidPassword error is returned if there isn't a password field in the body", t, func() {
		signIn := models.UserSignIn{
			Email: "email.email@ons.gov.uk",
		}

		validationErrors := *signIn.ValidateCredentials(ctx)
		So(len(validationErrors), ShouldEqual, 1)
		castErr := validationErrors[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidPasswordError)
		So(castErr.Description, ShouldEqual, models.InvalidPasswordDescription)
	})

	Convey("an InvalidPassword error is returned if there is an empty password field in the body", t, func() {
		signIn := models.UserSignIn{
			Password: "",
			Email:    "email.email@ons.gov.uk",
		}

		validationErrors := *signIn.ValidateCredentials(ctx)
		So(len(validationErrors), ShouldEqual, 1)
		castErr := validationErrors[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidPasswordError)
		So(castErr.Description, ShouldEqual, models.InvalidPasswordDescription)
	})

	Convey("an InvalidEmail error is returned if there isn't an email field in the body", t, func() {
		signIn := models.UserSignIn{
			Password: "password",
		}

		validationErrors := *signIn.ValidateCredentials(ctx)
		So(len(validationErrors), ShouldEqual, 1)
		castErr := validationErrors[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidEmailError)
		So(castErr.Description, ShouldEqual, models.InvalidEmailDescription)
	})

	Convey("an InvalidEmail error is returned if there is an empty email field in the body", t, func() {
		signIn := models.UserSignIn{
			Password: "password",
			Email:    "",
		}

		validationErrors := *signIn.ValidateCredentials(ctx)
		So(len(validationErrors), ShouldEqual, 1)
		castErr := validationErrors[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidEmailError)
		So(castErr.Description, ShouldEqual, models.InvalidEmailDescription)
	})

	Convey("an InvalidEmail error is returned if the email doesn't conform to the expected format", t, func() {
		signIn := models.UserSignIn{
			Password: "password",
			Email:    "email",
		}

		validationErrors := *signIn.ValidateCredentials(ctx)
		So(len(validationErrors), ShouldEqual, 1)
		castErr := validationErrors[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidEmailError)
		So(castErr.Description, ShouldEqual, models.InvalidEmailDescription)
	})

	Convey("an InvalidPassword and InvalidEmail error are returned if invalid email address and password are provided", t, func() {
		signIn := models.UserSignIn{
			Password: "",
			Email:    "",
		}

		validationErrors := *signIn.ValidateCredentials(ctx)
		So(len(validationErrors), ShouldEqual, 2)
		castErr := validationErrors[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidPasswordError)
		So(castErr.Description, ShouldEqual, models.InvalidPasswordDescription)
		castErr = validationErrors[1].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidEmailError)
		So(castErr.Description, ShouldEqual, models.InvalidEmailDescription)
	})
}

func TestUserSignIn_BuildOldSessionTerminationRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito AdminUserGlobalSignOutInput request body", t, func() {

		signIn := models.UserSignIn{
			Email:    "email.email@ons.gov.uk",
			Password: "password",
		}

		userPoolId := "eu-west-99-asegrh"

		response := signIn.BuildOldSessionTerminationRequest(userPoolId)

		So(*response.Username, ShouldEqual, signIn.Email)
		So(*response.UserPoolId, ShouldEqual, userPoolId)
	})
}

func TestUserSignIn_BuildCognitoRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito InitiateAuthInput request body", t, func() {

		signIn := models.UserSignIn{
			Email:    "email.email@ons.gov.uk",
			Password: "password",
		}

		clientId := "awsclientid"
		clientSecret := "awsSectret"
		clientAuthFlow := "authflow"

		response := signIn.BuildCognitoRequest(clientId, clientSecret, clientAuthFlow)

		So(*response.AuthParameters["USERNAME"], ShouldEqual, signIn.Email)
		So(*response.AuthParameters["PASSWORD"], ShouldEqual, signIn.Password)
		So(*response.AuthParameters["SECRET_HASH"], ShouldNotBeEmpty)
		So(*response.AuthFlow, ShouldResemble, "authflow")
		So(*response.ClientId, ShouldResemble, "awsclientid")
	})
}

func TestUserSignIn_BuildSuccessfulJsonResponse(t *testing.T) {
	ctx := context.Background()

	Convey("returns an InternalServerError if the Cognito response does not meet expected format", t, func() {
		signIn := models.UserSignIn{}
		result := cognitoidentityprovider.InitiateAuthOutput{}

		response, err := signIn.BuildSuccessfulJsonResponse(ctx, &result)

		So(response, ShouldBeNil)
		castErr := err.(*models.Error)
		So(castErr.Code, ShouldEqual, models.InternalError)
		So(castErr.Description, ShouldEqual, models.UnrecognisedCognitoResponseDescription)
	})

	Convey("returns a byte array of the response JSON", t, func() {
		var expirationLength int64 = 300
		signIn := models.UserSignIn{}
		result := cognitoidentityprovider.InitiateAuthOutput{
			AuthenticationResult: &cognitoidentityprovider.AuthenticationResultType{
				ExpiresIn: &expirationLength,
			},
		}

		response, err := signIn.BuildSuccessfulJsonResponse(ctx, &result)

		So(err, ShouldBeNil)
		So(reflect.TypeOf(response), ShouldEqual, reflect.TypeOf([]byte{}))
		var body map[string]interface{}
		err = json.Unmarshal(response, &body)
		So(body["expirationTime"], ShouldNotBeNil)
	})
}
