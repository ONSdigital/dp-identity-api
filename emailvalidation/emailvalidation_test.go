package emailvalidation

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEmailValidationConformsToExpectedFormat(t *testing.T) {

	Convey("The email conforms to the expected format and is validated", t, func() {

		email := "email.email@ons.gov.uk"

		emailResponse := IsEmailValid(email)
		So(emailResponse, ShouldBeTrue)
	})

	Convey("The email conforms to the expected format and is validated", t, func() {

		email := "email.email@domain.host"

		emailResponse := IsEmailValid(email)
		So(emailResponse, ShouldBeTrue)
	})

	Convey("The single domain email does not conform to the expected format and is validated", t, func() {

		email := "emailemail@domainhost"

		emailResponse := IsEmailValid(email)
		So(emailResponse, ShouldBeFalse)
	})

	Convey("The empty email does not conform to the expected format and is validated", t, func() {

		email := ""

		emailResponse := IsEmailValid(email)
		So(emailResponse, ShouldBeFalse)
	})

	Convey("The small email does not conform to the expected format and is validated", t, func() {

		email := "aaa"

		emailResponse := IsEmailValid(email)
		So(emailResponse, ShouldBeFalse)
	})

	Convey("The large email does not conform to the expected format and is validated", t, func() {
		str1 := "a"
		email := strings.Repeat(str1, 260)

		emailResponse := IsEmailValid(email)
		So(emailResponse, ShouldBeFalse)
	})

	Convey("The email has more than one @ does not conform to the expected format and is validated", t, func() {
		email := "string@string@string.string"

		emailResponse := IsEmailValid(email)
		So(emailResponse, ShouldBeFalse)
	})

	Convey("The email has .. in prefix and does not conform to the expected format and is validated", t, func() {
		email := "string..string@string.string"

		emailResponse := IsEmailValid(email)
		So(emailResponse, ShouldBeFalse)
	})

	Convey("The email has . at the start of prefix and does not conform to the expected format and is validated", t, func() {
		email := ".string.string@string.string"

		emailResponse := IsEmailValid(email)
		So(emailResponse, ShouldBeFalse)
	})

}
