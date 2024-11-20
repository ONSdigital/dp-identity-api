package models_test

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/ONSdigital/dp-identity-api/v2/api"
	"github.com/ONSdigital/dp-identity-api/v2/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	userID       = "abcd1234"
	userPoolID   = "euwest-99-aabbcc"
	clientID     = "awsclientid"
	clientSecret = "awsSecret"
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

		userPoolIDVar := userPoolID
		response := models.UsersList{}.BuildListUserRequest(filterString, requiredAttribute, limit, nil, &userPoolIDVar)

		So(reflect.TypeOf(*response), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.ListUsersInput{}))
		So(*response.UserPoolId, ShouldEqual, userPoolID)
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
					Enabled:    aws.Bool(true),
					UserStatus: aws.String("CONFIRMED"),
					Username:   aws.String("user-1"),
				},
				{
					Enabled:    aws.Bool(true),
					UserStatus: aws.String("CONFIRMED"),
					Username:   aws.String("user-2"),
				},
			},
		}
		userList := models.UsersList{}
		userList.MapCognitoUsers(&cognitoResponse.Users)

		So(len(userList.Users), ShouldEqual, len(cognitoResponse.Users))
		So(userList.Count, ShouldEqual, len(cognitoResponse.Users))
	})
}

func TestUsersList_SetUsers(t *testing.T) {
	Convey("adds the supplied users to the users attribute and sets the count", t, func() {
		listOfUsers := []models.UserParams{
			{
				Forename:    "Jane",
				Lastname:    "Doe",
				Email:       "jane.doe@ons.gov.uk",
				Status:      "Confirmed",
				Active:      true,
				ID:          "user-1",
				StatusNotes: "",
			},
			{
				Forename:    "John",
				Lastname:    "Doe",
				Email:       "john.doe@ons.gov.uk",
				Status:      "Confirmed",
				Active:      true,
				ID:          "user-2",
				StatusNotes: "",
			},
		}
		userList := models.UsersList{}
		userList.SetUsers(&listOfUsers)

		So(len(userList.Users), ShouldEqual, len(listOfUsers))
		So(userList.Count, ShouldEqual, len(listOfUsers))
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

		response, err := userList.BuildSuccessfulJSONResponse(ctx)

		So(err, ShouldBeNil)
		So(reflect.TypeOf(response), ShouldEqual, reflect.TypeOf([]byte{}))
		var body map[string]interface{}
		err = json.Unmarshal(response, &body)
		So(err, ShouldBeNil)
		So(body["count"], ShouldEqual, 1)
		usersJSON := body["users"].([]interface{})
		userJSON := usersJSON[0].(map[string]interface{})
		So(userJSON["id"], ShouldEqual, name)
		So(userJSON["status"], ShouldEqual, status)
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
	allowedDomains := []string{"@ons.gov.uk", "@ext.ons.gov.uk"}

	Convey("returns an InvalidForename error if an invalid forename is submitted", t, func() {
		user := models.UserParams{
			Email:    "email.email@ons.gov.uk",
			Forename: "",
			Lastname: "Smith",
		}

		errs := user.ValidateRegistration(ctx, allowedDomains)

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

		errs := user.ValidateRegistration(ctx, allowedDomains)

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

		errs := user.ValidateRegistration(ctx, allowedDomains)

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

		errs := user.ValidateRegistration(ctx, allowedDomains)

		So(len(errs), ShouldEqual, 1)
		castErr := errs[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidEmailError)
		So(castErr.Description, ShouldEqual, models.InvalidEmailDescription)
	})
}

func TestUserParams_ValidateUpdate(t *testing.T) {
	ctx := context.Background()

	invalidStatusNotes := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras eu turpis libero. Sed convallis pharetra mollis. Mauris ex nisi, finibus in mi quis, tincidunt pulvinar risus. Ut iaculis lobortis nisl. Suspendisse venenatis ante congue erat posuere, eget mattis massa facilisis. Vivamus bibendum pharetra suscipit. Integer laoreet molestie velit, vitae euismod ligula dictum eu. Phasellus a fermentum metus, nec dignissim ex. Sed dolor lectus, sollicitudin sit amet imperdiet eget, fringilla nec felis. Morbi commodo diam massa, sed interdum tellus sit"

	Convey("returns an InvalidForename error if an invalid forename is submitted", t, func() {
		user := models.UserParams{
			Forename:    "",
			Lastname:    "Smith",
			StatusNotes: "",
		}

		errs := user.ValidateUpdate(ctx)

		So(len(errs), ShouldEqual, 1)
		castErr := errs[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidForenameError)
		So(castErr.Description, ShouldEqual, models.InvalidForenameErrorDescription)
	})

	Convey("returns an InvalidSurname error if an invalid surname is submitted", t, func() {
		user := models.UserParams{
			Forename:    "Stan",
			Lastname:    "",
			StatusNotes: "",
		}

		errs := user.ValidateUpdate(ctx)

		So(len(errs), ShouldEqual, 1)
		castErr := errs[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidSurnameError)
		So(castErr.Description, ShouldEqual, models.InvalidSurnameErrorDescription)
	})

	Convey("returns an InvalidStatusNotes error if an invalid status notes is submitted", t, func() {
		user := models.UserParams{
			Forename:    "Stan",
			Lastname:    "Smith",
			StatusNotes: invalidStatusNotes,
		}

		errs := user.ValidateUpdate(ctx)

		So(len(errs), ShouldEqual, 1)
		castErr := errs[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidStatusNotesError)
		So(castErr.Description, ShouldEqual, models.TooLongStatusNotesDescription)
	})

	Convey("returns an InvalidForename and InvalidSurname errors if no forename or lastname are submitted", t, func() {
		user := models.UserParams{
			Forename:    "",
			Lastname:    "",
			StatusNotes: "",
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

	Convey("returns an InvalidForename and InvalidStatusNotes errors if no forename and invalid notes are submitted", t, func() {
		user := models.UserParams{
			Forename:    "",
			Lastname:    "Smith",
			StatusNotes: invalidStatusNotes,
		}

		errs := user.ValidateUpdate(ctx)

		So(len(errs), ShouldEqual, 2)
		castErr := errs[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidForenameError)
		So(castErr.Description, ShouldEqual, models.InvalidForenameErrorDescription)
		castErr = errs[1].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidStatusNotesError)
		So(castErr.Description, ShouldEqual, models.TooLongStatusNotesDescription)
	})

	Convey("returns an InvalidSurname and InvalidStatusNotes errors if no surname and invalid notes are submitted", t, func() {
		user := models.UserParams{
			Forename:    "Stan",
			Lastname:    "",
			StatusNotes: invalidStatusNotes,
		}

		errs := user.ValidateUpdate(ctx)

		So(len(errs), ShouldEqual, 2)
		castErr := errs[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidSurnameError)
		So(castErr.Description, ShouldEqual, models.InvalidSurnameErrorDescription)
		castErr = errs[1].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidStatusNotesError)
		So(castErr.Description, ShouldEqual, models.TooLongStatusNotesDescription)
	})

	Convey("returns an InvalidForename, InvalidSurname and InvalidStatusNotes errors if no forename or surname and invalid notes are submitted", t, func() {
		user := models.UserParams{
			Forename:    "",
			Lastname:    "",
			StatusNotes: invalidStatusNotes,
		}

		errs := user.ValidateUpdate(ctx)

		So(len(errs), ShouldEqual, 3)
		castErr := errs[0].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidForenameError)
		So(castErr.Description, ShouldEqual, models.InvalidForenameErrorDescription)
		castErr = errs[1].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidSurnameError)
		So(castErr.Description, ShouldEqual, models.InvalidSurnameErrorDescription)
		castErr = errs[2].(*models.Error)
		So(castErr.Code, ShouldEqual, models.InvalidStatusNotesError)
		So(castErr.Description, ShouldEqual, models.TooLongStatusNotesDescription)
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

		response := user.BuildCreateUserRequest(userID, userPoolID)

		So(reflect.TypeOf(*response), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.AdminCreateUserInput{}))
		So(*response.Username, ShouldEqual, userID)
		So(*response.UserPoolId, ShouldEqual, userPoolID)
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

		response := user.BuildUpdateUserRequest(userPoolID)

		So(reflect.TypeOf(*response), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.AdminUpdateUserAttributesInput{}))
		So(*response.Username, ShouldEqual, user.ID)
		So(*response.UserPoolId, ShouldEqual, userPoolID)
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

		response, err := createdUser.BuildSuccessfulJSONResponse(ctx)

		So(err, ShouldBeNil)
		So(reflect.TypeOf(response), ShouldEqual, reflect.TypeOf([]byte{}))
		var userJSON map[string]interface{}
		err = json.Unmarshal(response, &userJSON)
		So(err, ShouldBeNil)
		So(userJSON["id"], ShouldEqual, name)
		So(userJSON["status"], ShouldEqual, status)
	})
}

func TestUserParams_BuildAdminGetUserRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito AdminGetUserInput request body", t, func() {
		user := models.UserParams{
			ID: userID,
		}

		request := user.BuildAdminGetUserRequest(userPoolID)

		So(reflect.TypeOf(*request), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.AdminGetUserInput{}))
		So(*request.Username, ShouldEqual, userID)
		So(*request.UserPoolId, ShouldEqual, userPoolID)
	})
}

func TestUserParams_BuildEnableUserRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito AdminEnableUserInput request body", t, func() {
		user := models.UserParams{
			ID: userID,
		}

		request := user.BuildEnableUserRequest(userPoolID)

		So(reflect.TypeOf(*request), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.AdminEnableUserInput{}))
		So(*request.Username, ShouldEqual, userID)
		So(*request.UserPoolId, ShouldEqual, userPoolID)
	})
}

func TestUserParams_BuildDisableUserRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito AdminDisableUserInput request body", t, func() {
		user := models.UserParams{
			ID: userID,
		}

		request := user.BuildDisableUserRequest(userPoolID)

		So(reflect.TypeOf(*request), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.AdminDisableUserInput{}))
		So(*request.Username, ShouldEqual, userID)
		So(*request.UserPoolId, ShouldEqual, userPoolID)
	})
}

func TestUserParams_MapCognitoDetails(t *testing.T) {
	Convey("maps the returned user details to the UserParam attributes", t, func() {
		var forename, surname, email, status, id = "Bob", "Smith", "email@ons.gov.uk", "CONFIRMED", "user-1"
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
			Enabled:    aws.Bool(true),
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
		var forename, surname, email, status, id = "Bob", "Smith", "email@ons.gov.uk", "CONFIRMED", "user-1"
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
			Enabled:    aws.Bool(true),
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

		clientAuthFlow := "authflow"

		response := signIn.BuildCognitoRequest(clientID, clientSecret, clientAuthFlow)

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

		response, err := signIn.BuildSuccessfulJSONResponse(ctx, &result, 1)

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

		response, err := signIn.BuildSuccessfulJSONResponse(ctx, &result, 1)

		So(err, ShouldBeNil)
		So(reflect.TypeOf(response), ShouldEqual, reflect.TypeOf([]byte{}))
		var body map[string]interface{}
		_ = json.Unmarshal(response, &body)
		So(body["expirationTime"], ShouldNotBeNil)
		So(body["refreshTokenExpirationTime"], ShouldNotBeNil)
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

		response := passwordChangeParams.BuildAuthChallengeResponseRequest(clientSecret, clientID, api.NewPasswordChallenge)

		So(*response.ChallengeResponses["USERNAME"], ShouldEqual, passwordChangeParams.Email)
		So(*response.ChallengeResponses["NEW_PASSWORD"], ShouldEqual, passwordChangeParams.NewPassword)
		So(*response.ChallengeResponses["SECRET_HASH"], ShouldNotBeEmpty)
		So(*response.ChallengeName, ShouldEqual, api.NewPasswordChallenge)
		So(*response.Session, ShouldEqual, passwordChangeParams.Session)
		So(*response.ClientId, ShouldResemble, clientID)
	})
}

func TestChangePassword_BuildAuthChallengeSuccessfulJsonResponse(t *testing.T) {
	ctx := context.Background()

	Convey("returns an InternalServerError if the Cognito response does not meet expected format", t, func() {
		passwordChangeParams := models.ChangePassword{}
		result := cognitoidentityprovider.RespondToAuthChallengeOutput{}

		response, err := passwordChangeParams.BuildAuthChallengeSuccessfulJSONResponse(ctx, &result, 1)

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

		response, err := passwordChangeParams.BuildAuthChallengeSuccessfulJSONResponse(ctx, &result, 1)

		So(err, ShouldBeNil)
		So(reflect.TypeOf(response), ShouldEqual, reflect.TypeOf([]byte{}))
		var body map[string]interface{}
		_ = json.Unmarshal(response, &body)
		So(body["expirationTime"], ShouldNotBeNil)
		So(body["refreshTokenExpirationTime"], ShouldNotBeNil)
	})
}

func TestChangePassword_ValidateForgottenPasswordRequest(t *testing.T) {
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
				[]string{models.InvalidUserIDError},
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
				[]string{models.InvalidUserIDError, models.InvalidTokenError},
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
				[]string{models.InvalidPasswordError, models.InvalidUserIDError},
			},
			{
				// missing VerificationToken, email and password
				"",
				"",
				"",
				[]string{models.InvalidPasswordError, models.InvalidUserIDError, models.InvalidTokenError},
			},
		}
		for _, tt := range missingParamsTests {
			passwordChangeParams := models.ChangePassword{
				ChangeType:        models.NewPasswordRequiredType,
				VerificationToken: tt.VerificationToken,
				Email:             tt.Email,
				NewPassword:       tt.Password,
			}

			validationErrs := passwordChangeParams.ValidateForgottenPasswordRequest(ctx)

			So(len(validationErrs), ShouldEqual, len(tt.ExpectedErrors))
			for i, expectedErrCode := range tt.ExpectedErrors {
				castErr := validationErrs[i].(*models.Error)
				So(castErr.Code, ShouldEqual, expectedErrCode)
			}
		}
	})
}

func TestForgottenPassword_BuildConfirmForgotPasswordRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito BuildConfirmForgotPasswordRequest request body", t, func() {
		passwordChangeParams := models.ChangePassword{
			VerificationToken: "verification_token",
			Email:             "email@gmail.com",
			NewPassword:       "Password2",
		}

		clientID := "awsclientid"
		clientSecret := "awsSecret"

		response := passwordChangeParams.BuildConfirmForgotPasswordRequest(clientSecret, clientID)

		So(*response.ClientId, ShouldResemble, clientID)
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

		response := passwordResetParams.BuildCognitoRequest(clientSecret, clientID)

		So(*response.Username, ShouldEqual, passwordResetParams.Email)
		So(*response.SecretHash, ShouldNotBeEmpty)
		So(*response.ClientId, ShouldResemble, clientID)
	})
}

func TestUserParams_BuildListUserGroupsRequest(t *testing.T) {
	Convey("builds a correctly populated Cognito AdminListUserGroupsInput request body with empty nextToken", t, func() {
		nextToken := ""
		user := models.UserParams{
			ID: userID,
		}

		request := user.BuildListUserGroupsRequest(userPoolID, nextToken)

		So(reflect.TypeOf(*request), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.AdminListGroupsForUserInput{}))
		So(*request.Username, ShouldEqual, userID)
		So(*request.UserPoolId, ShouldEqual, userPoolID)
	})

	Convey("builds a correctly populated Cognito AdminDisableUserInput request body with nextToken", t, func() {
		nextToken := "abc1234"
		user := models.UserParams{
			ID: userID,
		}

		request := user.BuildListUserGroupsRequest(userPoolID, nextToken)

		So(reflect.TypeOf(*request), ShouldEqual, reflect.TypeOf(cognitoidentityprovider.AdminListGroupsForUserInput{}))
		So(*request.Username, ShouldEqual, userID)
		So(*request.UserPoolId, ShouldEqual, userPoolID)
	})
}

func TestListUserGroups_BuildListUserGroupsSuccessfulJsonResponse(t *testing.T) {
	Convey("add the returned groups for given user", t, func() {
		ctx := context.Background()
		input := models.ListUserGroups{}

		timestamp := time.Now()
		result := &cognitoidentityprovider.AdminListGroupsForUserOutput{
			Groups: []*cognitoidentityprovider.GroupType{
				{
					CreationDate:     &timestamp,
					Description:      aws.String("A test group1"),
					GroupName:        aws.String("test-group1"),
					LastModifiedDate: &timestamp,
					Precedence:       aws.Int64(4),
					RoleArn:          aws.String(""),
					UserPoolId:       aws.String(""),
				},
				{
					CreationDate:     &timestamp,
					Description:      aws.String("A test group1"),
					GroupName:        aws.String("test-group1"),
					LastModifiedDate: &timestamp,
					Precedence:       aws.Int64(4),
					RoleArn:          aws.String(""),
					UserPoolId:       aws.String(""),
				},
			},
		}

		response, err := input.BuildListUserGroupsSuccessfulJSONResponse(ctx, result)
		So(err, ShouldBeNil)
		So(reflect.TypeOf(response), ShouldEqual, reflect.TypeOf([]byte{}))

		var userGroupsJSON models.ListUserGroups
		err = json.Unmarshal(response, &userGroupsJSON)
		So(err, ShouldBeNil)
		So(len(userGroupsJSON.Groups), ShouldEqual, len(result.Groups))
		So(userGroupsJSON.Count, ShouldEqual, len(result.Groups))
		So(userGroupsJSON.NextToken, ShouldBeNil)

		So(*userGroupsJSON.Groups[0].ID, ShouldEqual, *result.Groups[0].GroupName)
		So(*userGroupsJSON.Groups[1].ID, ShouldEqual, *result.Groups[1].GroupName)
		So(*userGroupsJSON.Groups[0].Name, ShouldEqual, *result.Groups[0].Description)
		So(*userGroupsJSON.Groups[1].Name, ShouldEqual, *result.Groups[1].Description)
	})

	Convey("Check empty response from cognito i.e valid user with no groups", t, func() {
		ctx := context.Background()
		input := models.ListUserGroups{}

		result := &cognitoidentityprovider.AdminListGroupsForUserOutput{}

		response, err := input.BuildListUserGroupsSuccessfulJSONResponse(ctx, result)
		So(err, ShouldBeNil)

		var userGroupsJSON models.ListUserGroups
		err = json.Unmarshal(response, &userGroupsJSON)
		So(err, ShouldBeNil)
		So(len(userGroupsJSON.Groups), ShouldEqual, len(result.Groups))
		So(userGroupsJSON.Count, ShouldEqual, 0)
		So(userGroupsJSON.NextToken, ShouldBeNil)
	})

	Convey("force nil return for cognitoidentityprovider.AdminListGroupsForUserOutput", t, func() {
		var result *cognitoidentityprovider.AdminListGroupsForUserOutput
		ctx := context.Background()
		input := models.ListUserGroups{}

		result = nil

		response, err := input.BuildListUserGroupsSuccessfulJSONResponse(ctx, result)
		castErr := err.(*models.Error)
		So(castErr.Code, ShouldEqual, models.InternalError)
		So(response, ShouldBeNil)
	})
}
