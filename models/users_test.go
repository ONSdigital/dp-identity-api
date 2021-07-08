package models_test

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/ONSdigital/dp-identity-api/api"

	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUsersList_BuildListUserRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito ListUsers request body", t, func() {

		user := models.UserParams{
			Email:    "email.email@ons.gov.uk",
			Forename: "Stan",
			Lastname: "Smith",
		}

		filterString := "email = \"" + user.Email + "\""
		requiredAttribute := "email"
		limit := int64(1)
		userPoolId := "euwest-99-aabbcc"

		response := models.UsersList{}.BuildListUserRequest(filterString, requiredAttribute, limit, &userPoolId)

		So(reflect.TypeOf(*response), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.ListUsersInput{}))
		So(*response.UserPoolId, ShouldEqual, userPoolId)
		So(*response.Limit, ShouldEqual, limit)
		So(*response.Filter, ShouldEqual, filterString)
		So(*response.AttributesToGet[0], ShouldEqual, requiredAttribute)
	})
}

func TestUsersList_MapCognitoUsers(t *testing.T) {
	Convey("adds the returned users to the users attribute and sets the count", t, func() {
		cognitoResponse := cognitoidentityprovider.ListUsersOutput{
			Users: []*cognitoidentityprovider.UserType{
				{
					UserStatus: aws.String("CONFIRMED"),
					Username:   aws.String("user-1"),
				},
				{
					UserStatus: aws.String("CONFIRMED"),
					Username:   aws.String("user-2"),
				},
			},
		}
		userList := models.UsersList{}
		userList.MapCognitoUsers(&cognitoResponse)

		So(len(userList.Users), ShouldEqual, len(cognitoResponse.Users))
		So(userList.Count, ShouldEqual, len(cognitoResponse.Users))
	})
}

func TestUsersList_BuildSuccessfulJsonResponse(t *testing.T) {
	Convey("returns a byte array of the response JSON", t, func() {
		ctx := context.Background()
		name, status := "abcd-efgh-ijkl-mnop", "UNCONFIRMED"
		user := models.UserParams{
			Status: status,
			ID:     name,
		}
		userList := models.UsersList{
			Users: []models.UserParams{
				user,
			},
			Count: 1,
		}

		response, err := userList.BuildSuccessfulJsonResponse(ctx)

		So(err, ShouldBeNil)
		So(reflect.TypeOf(response), ShouldEqual, reflect.TypeOf([]byte{}))
		var body map[string]interface{}
		err = json.Unmarshal(response, &body)
		So(err, ShouldBeNil)
		So(body["count"], ShouldEqual, 1)
		usersJson := body["users"].([]interface{})
		userJson := usersJson[0].(map[string]interface{})
		So(userJson["id"], ShouldEqual, name)
		So(userJson["status"], ShouldEqual, status)
	})
}

func TestUserParams_GeneratePassword(t *testing.T) {
	Convey("adds a password to the UserParams object", t, func() {
		ctx := context.Background()

		user := models.UserParams{}

		err := user.GeneratePassword(ctx)

		So(err, ShouldBeNil)
		So(user.Password, ShouldNotBeNil)
		So(user.Password, ShouldNotEqual, "")
	})
}

func TestUserParams_ValidateRegistration(t *testing.T) {
	ctx := context.Background()

	Convey("returns an InvalidForename error if an invalid forename is submitted", t, func() {
		user := models.UserParams{
			Email:    "email.email@ons.gov.uk",
			Forename: "",
			Lastname: "Smith",
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
			Lastname: "",
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
			Lastname: "Smith",
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
			Lastname: "Smith",
		}

		errs := user.ValidateRegistration(ctx)

		So(len(errs), ShouldEqual, 1)
		castErr := errs[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidEmailError)
		So(castErr.Description, ShouldEqual, models.InvalidEmailDescription)
	})
}

func TestUserParams_ValidateUpdate(t *testing.T) {
	ctx := context.Background()

	Convey("returns an InvalidForename error if an invalid forename is submitted", t, func() {
		user := models.UserParams{
			Forename: "",
			Lastname: "Smith",
		}

		errs := user.ValidateUpdate(ctx)

		So(len(errs), ShouldEqual, 1)
		castErr := errs[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidForenameError)
		So(castErr.Description, ShouldEqual, models.InvalidForenameErrorDescription)
	})

	Convey("returns an InvalidSurname error if an invalid surname is submitted", t, func() {
		user := models.UserParams{
			Forename: "Stan",
			Lastname: "",
		}

		errs := user.ValidateUpdate(ctx)

		So(len(errs), ShouldEqual, 1)
		castErr := errs[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidSurnameError)
		So(castErr.Description, ShouldEqual, models.InvalidSurnameErrorDescription)
	})

	Convey("returns an InvalidForename and InvalidSurname errors if no forename or lastname is submitted", t, func() {
		user := models.UserParams{
			Forename: "",
			Lastname: "",
		}

		errs := user.ValidateUpdate(ctx)

		So(len(errs), ShouldEqual, 2)
		castErr := errs[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidForenameError)
		So(castErr.Description, ShouldEqual, models.InvalidForenameErrorDescription)
		castErr = errs[1].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidSurnameError)
		So(castErr.Description, ShouldEqual, models.InvalidSurnameErrorDescription)
	})
}

func TestUserParams_CheckForDuplicateEmail(t *testing.T) {
	ctx := context.Background()

	Convey("returns nothing if there is no user returned from the ListUser request", t, func() {
		user := models.UserParams{
			Email:    "email.email@ons.gov.uk",
			Forename: "Stan",
			Lastname: "Smith",
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
			Lastname: "Smith",
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

func TestUserParams_BuildCreateUserRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito AdminUserCreateInput request body", t, func() {

		user := models.UserParams{
			Email:    "email.email@ons.gov.uk",
			Forename: "Stan",
			Lastname: "Smith",
		}

		userId := uuid.NewString()
		userPoolId := "euwest-99-aabbcc"

		response := user.BuildCreateUserRequest(userId, userPoolId)

		So(reflect.TypeOf(*response), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.AdminCreateUserInput{}))
		So(*response.Username, ShouldEqual, userId)
		So(*response.UserPoolId, ShouldEqual, userPoolId)
		So(*response.UserAttributes[0].Value, ShouldEqual, user.Forename)
		So(*response.UserAttributes[1].Value, ShouldEqual, user.Lastname)
		So(*response.UserAttributes[2].Value, ShouldEqual, user.Email)
	})
}

func TestUserParams_BuildUpdateUserRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito AdminUpdateUserAttributeInput request body", t, func() {

		user := models.UserParams{
			ID:          "abcd1234",
			Forename:    "Stan",
			Lastname:    "Smith",
			StatusNotes: "user suspended",
		}

		userPoolId := "euwest-99-aabbcc"

		response := user.BuildUpdateUserRequest(userPoolId)

		So(reflect.TypeOf(*response), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.AdminUpdateUserAttributesInput{}))
		So(*response.Username, ShouldEqual, user.ID)
		So(*response.UserPoolId, ShouldEqual, userPoolId)
		So(*response.UserAttributes[0].Value, ShouldEqual, user.Forename)
		So(*response.UserAttributes[1].Value, ShouldEqual, user.Lastname)
		So(*response.UserAttributes[2].Value, ShouldEqual, user.StatusNotes)
	})
}

func TestUserParams_BuildSuccessfulJsonResponse(t *testing.T) {
	Convey("returns a byte array of the response JSON", t, func() {
		ctx := context.Background()
		name, status := "abcd-efgh-ijkl-mnop", "UNCONFIRMED"
		createdUser := models.UserParams{
			Status: status,
			ID:     name,
		}

		response, err := createdUser.BuildSuccessfulJsonResponse(ctx)

		So(err, ShouldBeNil)
		So(reflect.TypeOf(response), ShouldEqual, reflect.TypeOf([]byte{}))
		var userJson map[string]interface{}
		err = json.Unmarshal(response, &userJson)
		So(err, ShouldBeNil)
		So(userJson["id"], ShouldEqual, name)
		So(userJson["status"], ShouldEqual, status)
	})
}

func TestUserParams_BuildAdminGetUserRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito AdminGetUserInput request body", t, func() {
		userId := "abcd1234"
		user := models.UserParams{
			ID: userId,
		}

		userPoolId := "euwest-99-aabbcc"

		request := user.BuildAdminGetUserRequest(userPoolId)

		So(reflect.TypeOf(*request), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.AdminGetUserInput{}))
		So(*request.Username, ShouldEqual, userId)
		So(*request.UserPoolId, ShouldEqual, userPoolId)
	})
}

func TestUserParams_BuildEnableUserRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito AdminEnableUserInput request body", t, func() {
		userId := "abcd1234"
		user := models.UserParams{
			ID: userId,
		}

		userPoolId := "euwest-99-aabbcc"

		request := user.BuildEnableUserRequest(userPoolId)

		So(reflect.TypeOf(*request), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.AdminEnableUserInput{}))
		So(*request.Username, ShouldEqual, userId)
		So(*request.UserPoolId, ShouldEqual, userPoolId)
	})
}

func TestUserParams_BuildDisableUserRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito AdminDisableUserInput request body", t, func() {
		userId := "abcd1234"
		user := models.UserParams{
			ID: userId,
		}

		userPoolId := "euwest-99-aabbcc"

		request := user.BuildDisableUserRequest(userPoolId)

		So(reflect.TypeOf(*request), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.AdminDisableUserInput{}))
		So(*request.Username, ShouldEqual, userId)
		So(*request.UserPoolId, ShouldEqual, userPoolId)
	})
}

func TestUserParams_MapCognitoDetails(t *testing.T) {
	Convey("maps the returned user details to the UserParam attributes", t, func() {
		var forename, surname, email, status, id string = "Bob", "Smith", "email@ons.gov.uk", "CONFIRMED", "user-1"
		cognitoUser := cognitoidentityprovider.UserType{
			Attributes: []*cognitoidentityprovider.AttributeType{
				{
					Name:  aws.String("given_name"),
					Value: &forename,
				},
				{
					Name:  aws.String("family_name"),
					Value: &surname,
				},
				{
					Name:  aws.String("email"),
					Value: &email,
				},
			},
			UserStatus: &status,
			Username:   &id,
		}
		user := models.UserParams{}.MapCognitoDetails(&cognitoUser)

		So(user.Forename, ShouldEqual, forename)
		So(user.Lastname, ShouldEqual, surname)
		So(user.Email, ShouldEqual, email)
		So(user.Status, ShouldEqual, status)
		So(user.ID, ShouldEqual, id)
	})
}

func TestUserParams_MapCognitoGetResponse(t *testing.T) {
	Convey("maps the returned user details to the UserParam attributes", t, func() {
		var forename, surname, email, status, id string = "Bob", "Smith", "email@ons.gov.uk", "CONFIRMED", "user-1"
		cognitoUser := cognitoidentityprovider.AdminGetUserOutput{
			UserAttributes: []*cognitoidentityprovider.AttributeType{
				{
					Name:  aws.String("given_name"),
					Value: &forename,
				},
				{
					Name:  aws.String("family_name"),
					Value: &surname,
				},
				{
					Name:  aws.String("email"),
					Value: &email,
				},
			},
			UserStatus: &status,
			Username:   &id,
		}
		user := models.UserParams{ID: id}
		user.MapCognitoGetResponse(&cognitoUser)

		So(user.Forename, ShouldEqual, forename)
		So(user.Lastname, ShouldEqual, surname)
		So(user.Email, ShouldEqual, email)
		So(user.Status, ShouldEqual, status)
		So(user.ID, ShouldEqual, id)
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

func TestChangePassword_ValidateNewPasswordRequiredRequest(t *testing.T) {
	ctx := context.Background()

	Convey("returns validation errors if required parameters are missing", t, func() {
		missingParamsTests := []struct {
			Session        string
			Email          string
			Password       string
			ExpectedErrors []string
		}{
			{
				// missing session
				"",
				"email@gmail.com",
				"Password2",
				[]string{models.InvalidChallengeSessionError},
			},
			{
				// missing email
				"auth-challenge-session",
				"",
				"Password2",
				[]string{models.InvalidEmailError},
			},
			{
				// missing password
				"auth-challenge-session",
				"email@gmail.com",
				"",
				[]string{models.InvalidPasswordError},
			},
			{
				// missing session and email
				"",
				"",
				"Password2",
				[]string{models.InvalidEmailError, models.InvalidChallengeSessionError},
			},
			{
				// missing session and password
				"",
				"email@gmail.com",
				"",
				[]string{models.InvalidPasswordError, models.InvalidChallengeSessionError},
			},
			{
				// missing email and password
				"auth-challenge-session",
				"",
				"",
				[]string{models.InvalidPasswordError, models.InvalidEmailError},
			},
			{
				// missing session, email and password
				"",
				"",
				"",
				[]string{models.InvalidPasswordError, models.InvalidEmailError, models.InvalidChallengeSessionError},
			},
		}
		for _, tt := range missingParamsTests {
			passwordChangeParams := models.ChangePassword{
				ChangeType:  models.NewPasswordRequiredType,
				Session:     tt.Session,
				Email:       tt.Email,
				NewPassword: tt.Password,
			}

			validationErrs := passwordChangeParams.ValidateNewPasswordRequiredRequest(ctx)

			So(len(validationErrs), ShouldEqual, len(tt.ExpectedErrors))
			for i, expectedErrCode := range tt.ExpectedErrors {
				castErr := validationErrs[i].(*models.Error)
				So(castErr.Code, ShouldEqual, expectedErrCode)
			}
		}
	})

	Convey("returns an empty slice if there are no validation failures", t, func() {
		passwordChangeParams := models.ChangePassword{
			ChangeType:  models.NewPasswordRequiredType,
			Session:     "auth-challenge-session",
			Email:       "email@gmail.com",
			NewPassword: "Password2",
		}
		validationErrs := passwordChangeParams.ValidateNewPasswordRequiredRequest(ctx)

		So(len(validationErrs), ShouldEqual, 0)
	})
}

func TestChangePassword_BuildAuthChallengeResponseRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito RespondToAuthChallengeInput request body", t, func() {

		passwordChangeParams := models.ChangePassword{
			ChangeType:  models.NewPasswordRequiredType,
			Session:     "auth-challenge-session",
			Email:       "email@gmail.com",
			NewPassword: "Password2",
		}

		clientId := "awsclientid"
		clientSecret := "awsSectret"

		response := passwordChangeParams.BuildAuthChallengeResponseRequest(clientSecret, clientId, api.NewPasswordChallenge)

		So(*response.ChallengeResponses["USERNAME"], ShouldEqual, passwordChangeParams.Email)
		So(*response.ChallengeResponses["NEW_PASSWORD"], ShouldEqual, passwordChangeParams.NewPassword)
		So(*response.ChallengeResponses["SECRET_HASH"], ShouldNotBeEmpty)
		So(*response.ChallengeName, ShouldEqual, api.NewPasswordChallenge)
		So(*response.Session, ShouldEqual, passwordChangeParams.Session)
		So(*response.ClientId, ShouldResemble, clientId)
	})
}

func TestChangePassword_BuildAuthChallengeSuccessfulJsonResponse(t *testing.T) {
	ctx := context.Background()

	Convey("returns an InternalServerError if the Cognito response does not meet expected format", t, func() {
		passwordChangeParams := models.ChangePassword{}
		result := cognitoidentityprovider.RespondToAuthChallengeOutput{}

		response, err := passwordChangeParams.BuildAuthChallengeSuccessfulJsonResponse(ctx, &result)

		So(response, ShouldBeNil)
		castErr := err.(*models.Error)
		So(castErr.Code, ShouldEqual, models.InternalError)
		So(castErr.Description, ShouldEqual, models.UnrecognisedCognitoResponseDescription)
	})

	Convey("returns a byte array of the response JSON", t, func() {
		var expirationLength int64 = 300
		passwordChangeParams := models.ChangePassword{}
		result := cognitoidentityprovider.RespondToAuthChallengeOutput{
			AuthenticationResult: &cognitoidentityprovider.AuthenticationResultType{
				ExpiresIn: &expirationLength,
			},
		}

		response, err := passwordChangeParams.BuildAuthChallengeSuccessfulJsonResponse(ctx, &result)

		So(err, ShouldBeNil)
		So(reflect.TypeOf(response), ShouldEqual, reflect.TypeOf([]byte{}))
		var body map[string]interface{}
		err = json.Unmarshal(response, &body)
		So(body["expirationTime"], ShouldNotBeNil)
	})
}

func TestChangePassword_ValidateForgottenPasswordRequiredRequest(t *testing.T) {
	ctx := context.Background()

	Convey("returns validation errors if required parameters are missing", t, func() {
		missingParamsTests := []struct {
			VerificationToken string
			Email             string
			Password          string
			ExpectedErrors    []string
		}{
			{
				// missing VerificationToken
				"",
				"email@gmail.com",
				"Password2",
				[]string{models.InvalidTokenError},
			},
			{
				// missing email
				"≈",
				"",
				"Password2",
				[]string{models.InvalidEmailError},
			},
			{
				// missing password
				"verification_token",
				"email@gmail.com",
				"",
				[]string{models.InvalidPasswordError},
			},
			{
				// missing VerificationToken and email
				"",
				"",
				"Password2",
				[]string{models.InvalidEmailError, models.InvalidTokenError},
			},
			{
				// missing VerificationToken and password
				"",
				"email@gmail.com",
				"",
				[]string{models.InvalidPasswordError, models.InvalidTokenError},
			},
			{
				// missing email and password
				"verification_token",
				"",
				"",
				[]string{models.InvalidPasswordError, models.InvalidEmailError},
			},
			{
				// missing VerificationToken, email and password
				"",
				"",
				"",
				[]string{models.InvalidPasswordError, models.InvalidEmailError, models.InvalidTokenError},
			},
		}
		for _, tt := range missingParamsTests {
			passwordChangeParams := models.ChangePassword{
				ChangeType:        models.NewPasswordRequiredType,
				VerificationToken: tt.VerificationToken,
				Email:             tt.Email,
				NewPassword:       tt.Password,
			}

			validationErrs := passwordChangeParams.ValidateForgottenPasswordRequiredRequest(ctx)

			So(len(validationErrs), ShouldEqual, len(tt.ExpectedErrors))
			for i, expectedErrCode := range tt.ExpectedErrors {
				castErr := validationErrs[i].(*models.Error)
				So(castErr.Code, ShouldEqual, expectedErrCode)
			}
		}
	})
}

func TestForgottenPassword_BuildConfirmForgotPasswordRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito RespondToAuthChallengeInput request body", t, func() {

		passwordChangeParams := models.ChangePassword{
			VerificationToken: "verification_token",
			Email:             "email@gmail.com",
			NewPassword:       "Password2",
		}

		clientId := "awsclientid"
		clientSecret := "awsSecret"

		response := passwordChangeParams.BuildConfirmForgotPasswordRequest(clientSecret, clientId)

		So(*response.ClientId, ShouldResemble, clientId)
		So(*response.ConfirmationCode, ShouldResemble, passwordChangeParams.VerificationToken)
		So(*response.Password, ShouldResemble, passwordChangeParams.NewPassword)
		So(*response.Username, ShouldResemble, passwordChangeParams.Email)

	})
}

func TestPasswordReset_Validate(t *testing.T) {
	ctx := context.Background()

	Convey("returns validation errors if required parameters are missing", t, func() {
		missingParamsTests := []struct {
			Email         string
			ExpectedError string
		}{
			{
				// missing email
				"",
				"InvalidEmail",
			},
		}
		for _, tt := range missingParamsTests {
			passwordChangeParams := models.PasswordReset{
				Email: tt.Email,
			}

			validationErr := passwordChangeParams.Validate(ctx)

			castErr := validationErr.(*models.Error)
			So(castErr.Code, ShouldEqual, tt.ExpectedError)
		}
	})

	Convey("returns nil if there are no validation failures", t, func() {
		passwordResetParams := models.PasswordReset{
			Email: "email@gmail.com",
		}
		validationErr := passwordResetParams.Validate(ctx)

		So(validationErr, ShouldBeNil)
	})
}

func TestPasswordReset_BuildCognitoRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito ForgotPasswordInput request body", t, func() {

		passwordResetParams := models.PasswordReset{
			Email: "email@gmail.com",
		}

		clientId := "awsclientid"
		clientSecret := "awsSectret"

		response := passwordResetParams.BuildCognitoRequest(clientSecret, clientId)

		So(*response.Username, ShouldEqual, passwordResetParams.Email)
		So(*response.SecretHash, ShouldNotBeEmpty)
		So(*response.ClientId, ShouldResemble, clientId)
	})
}
