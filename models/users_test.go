package models_test

import (
	"context"
	"encoding/json"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	"testing"
)

func TestUserParams_GeneratePassword(t *testing.T) {
	Convey("adds a password to the UserParams object", t, func() {
		ctx := context.Background()

		user := models.UserParams{}

		err := user.GeneratePassword(ctx)

		So(err, ShouldBeNil)
		So(user.Password, ShouldNotBeNil)
	})
}

func TestUserParams_ValidateRegistration(t *testing.T) {
	ctx := context.Background()

	Convey("returns an InvalidForename error if an invalid forename is submitted", t, func() {
		user := models.UserParams{
			Email:    "email.email@ons.gov.uk",
			Forename: "",
			Surname:  "Smith",
		}

		errs := user.ValidateRegistration(ctx)

		So(len(errs), ShouldEqual, 1)
		castErr := errs[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidForenameError)
		So(castErr.Description, ShouldEqual, models.InvalidForenameErrorDescription)
	})

	Convey("returns an InvalidSurname error if an invalid surname is submitted", t, func() {
		user := models.UserParams{
			Email:    "email.email@ons.gov.uk",
			Forename: "Stan",
			Surname:  "",
		}

		errs := user.ValidateRegistration(ctx)

		So(len(errs), ShouldEqual, 1)
		castErr := errs[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidSurnameError)
		So(castErr.Description, ShouldEqual, models.InvalidSurnameErrorDescription)
	})

	Convey("returns an InvalidEmail error if an invalid email is submitted", t, func() {
		user := models.UserParams{
			Email:    "email",
			Forename: "Stan",
			Surname:  "Smith",
		}

		errs := user.ValidateRegistration(ctx)

		So(len(errs), ShouldEqual, 1)
		castErr := errs[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidEmailError)
		So(castErr.Description, ShouldEqual, models.InvalidEmailDescription)
	})

	Convey("returns an InvalidEmail error if a non ONS email is submitted", t, func() {
		user := models.UserParams{
			Email:    "email@gmail.com",
			Forename: "Stan",
			Surname:  "Smith",
		}

		errs := user.ValidateRegistration(ctx)

		So(len(errs), ShouldEqual, 1)
		castErr := errs[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidEmailError)
		So(castErr.Description, ShouldEqual, models.InvalidEmailDescription)
	})
}

func TestUserParams_CheckForDuplicateEmail(t *testing.T) {
	ctx := context.Background()

	Convey("returns nothing if there is no user returned from the ListUser request", t, func() {
		user := models.UserParams{
			Email:    "email.email@ons.gov.uk",
			Forename: "Stan",
			Surname:  "Smith",
		}

		listUserResponse := cognitoidentityprovider.ListUsersOutput{
			Users: []*cognitoidentityprovider.UserType{},
		}

		err := user.CheckForDuplicateEmail(ctx, &listUserResponse)

		So(err, ShouldBeNil)
	})

	Convey("returns an InvalidEmail error if there is a user returned from the ListUser request", t, func() {
		user := models.UserParams{
			Email:    "email.email@ons.gov.uk",
			Forename: "Stan",
			Surname:  "Smith",
		}

		name, status := "abcd-efgh-ijkl-mnop", "UNCONFIRMED"
		listUserResponse := cognitoidentityprovider.ListUsersOutput{
			Users: []*cognitoidentityprovider.UserType{
				{
					Username:   &name,
					UserStatus: &status,
				},
			},
		}

		err := user.CheckForDuplicateEmail(ctx, &listUserResponse)

		castErr := err.(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidEmailError)
		So(castErr.Description, ShouldEqual, models.DuplicateEmailDescription)
	})
}

func TestUserParams_BuildListUserRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito ListUsers request body", t, func() {

		user := models.UserParams{
			Email:    "email.email@ons.gov.uk",
			Forename: "Stan",
			Surname:  "Smith",
		}

		filterString := "email = \"" + user.Email + "\""
		requiredAttribute := "email"
		limit := int64(1)
		userPoolId := "euwest-99-aabbcc"

		response := user.BuildListUserRequest(filterString, requiredAttribute, limit, &userPoolId)

		So(reflect.TypeOf(*response), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.ListUsersInput{}))
		So(*response.UserPoolId, ShouldEqual, userPoolId)
		So(*response.Limit, ShouldEqual, limit)
		So(*response.Filter, ShouldEqual, filterString)
		So(*response.AttributesToGet[0], ShouldEqual, requiredAttribute)
	})
}

func TestUserParams_BuildCreateUserRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito AdminUserCreateInput request body", t, func() {

		user := models.UserParams{
			Email:    "email.email@ons.gov.uk",
			Forename: "Stan",
			Surname:  "Smith",
		}

		userId := uuid.NewString()
		userPoolId := "euwest-99-aabbcc"

		response := user.BuildCreateUserRequest(userId, userPoolId)

		So(reflect.TypeOf(*response), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.AdminCreateUserInput{}))
		So(*response.Username, ShouldEqual, userId)
		So(*response.UserPoolId, ShouldEqual, userPoolId)
		So(*response.UserAttributes[0].Value, ShouldEqual, user.Forename)
		So(*response.UserAttributes[1].Value, ShouldEqual, user.Surname)
		So(*response.UserAttributes[2].Value, ShouldEqual, user.Email)
	})
}

func TestUserParams_BuildSuccessfulJsonResponse(t *testing.T) {
	Convey("returns a byte array of the response JSON", t, func() {
		ctx := context.Background()
		name, status := "abcd-efgh-ijkl-mnop", "UNCONFIRMED"
		user := models.UserParams{}
		result := cognitoidentityprovider.AdminCreateUserOutput{
			User: &cognitoidentityprovider.UserType{
				Username:   &name,
				UserStatus: &status,
			},
		}

		response, err := user.BuildSuccessfulJsonResponse(ctx, &result)

		So(err, ShouldBeNil)
		So(reflect.TypeOf(response), ShouldEqual, reflect.TypeOf([]byte{}))
		var body map[string]interface{}
		err = json.Unmarshal(response, &body)
		So(err, ShouldBeNil)
		userJson := body["User"].(map[string]interface{})
		So(userJson["Username"], ShouldEqual, name)
		So(userJson["UserStatus"], ShouldEqual, status)
	})
}

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