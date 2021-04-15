package api

import(
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)
func TestEmailAndPasswordHaveBeenProvided(t *testing.T) {

	Convey("Both an email and password have been provided and the body is validated", t, func() {
		//Given the body we need a function that will check that the body contains both an email and a password and these are not null

		email := "email"
		password := "password"

		body := make(map[string]string)
		body["email"] = email
		body["password"] = password

		emailResponse, passwordResponse := bodyValidation(body)

		So(emailResponse, ShouldBeTrue)
		So(passwordResponse, ShouldBeTrue)
	})

	Convey("There isn't an email field in the body and the body isn't validated", t, func() {

		password := "password"
		body := make(map[string]string)
		body["password"] = password

		emailResponse, passwordResponse := bodyValidation(body)
		So(emailResponse, ShouldBeFalse)
		So(passwordResponse, ShouldBeTrue)
	})

	Convey("There isn't an email value in the body and the body isn't validated", t, func() {
		
		password := "password"
		email := ""
		body := make(map[string]string)
		body["email"] = email
		body["password"] = password

		emailResponse, passwordResponse := bodyValidation(body)
		So(emailResponse, ShouldBeFalse)
		So(passwordResponse, ShouldBeTrue)
	})
}

func TestEmailConformsToExpectedFormat(t *testing.T) {

	Convey("The email conforms to the expected format and is validated", t, func() {

		email := "email.email@ons.gov.uk"
		//email = "email.email@ext.ons.gov.uk"

		body := make(map[string]string)
		body["email"] = email

		emailResponse := emailValidation(body)
		So(emailResponse, ShouldBeTrue)
	})

	Convey("There isn't an email field in the body and it isn't validated", t, func() {

		password := "password"
		body := make(map[string]string)
		body["password"] = password

		emailResponse := emailValidation(body)
		So(emailResponse, ShouldBeFalse)
	})

	Convey("The email doesn't conform to the expected format and it isn't validated", t, func() {

		email := "email@ons2.gov.uk"

		body := make(map[string]string)
		body["email"] = email

		emailResponse := emailValidation(body)
		So(emailResponse, ShouldBeFalse)
	})
}