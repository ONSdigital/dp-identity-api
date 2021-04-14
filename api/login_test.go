package api

import(
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)
func TestAValidEmailHasBeenProvided(t *testing.T) {

	Convey("Both an email and password have been provided", t, func() {
		//Given the body we need a function that will check that the body contains both an email and a password and these are not null

		email := "invalidEmail"
		password := "password"

		body := make(map[string]string)
		body["email"] = email
		body["password"] = password

		emailResponse, passwordResponse, _ := funcName(body)

		So(emailResponse, ShouldBeTrue)
		So(passwordResponse, ShouldBeTrue)

	})
}