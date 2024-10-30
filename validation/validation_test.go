package validation

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

	Convey("A capitalised email conforms to the expected format and is validated", t, func() {
		email := "EMAIL.EMAIL@DOMAIN.HOST"

		emailResponse := IsEmailValid(email)
		So(emailResponse, ShouldBeTrue)
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

	Convey("The email conforms to the expected format and is validated", t, func() {
		email := "\"email.email\"@ons.gov.uk"
		emailResponse := IsEmailValid(email)
		So(emailResponse, ShouldBeTrue)
	})
}

func TestEmailDomainValidationCorrectlyValidatesAllowedDomains(t *testing.T) {
	allowedDomains := []string{"@ons.gov.uk", "@ext.ons.gov.uk"}

	Convey("An empty email does not conform to the expected format and is not validated", t, func() {
		email := ""
		emailResponse := IsAllowedEmailDomain(email, allowedDomains)
		So(emailResponse, ShouldBeFalse)
	})

	Convey("A badly formatted email is not validated", t, func() {
		email := "string@string@string.string"
		emailResponse := IsAllowedEmailDomain(email, allowedDomains)
		So(emailResponse, ShouldBeFalse)
	})

	Convey("The email conforms to the expected format but is not an allowed domain so it is not validated", t, func() {
		email := "email.email@gmail.com"
		emailResponse := IsAllowedEmailDomain(email, allowedDomains)
		So(emailResponse, ShouldBeFalse)
	})

	Convey("The allowed domain is only part of the domain so it is not validated", t, func() {
		email := "email.email@ons.gov.uk.test.com"
		emailResponse := IsAllowedEmailDomain(email, allowedDomains)
		So(emailResponse, ShouldBeFalse)
	})

	Convey("The email conforms to the expected format and is an allowed domain then it is validated", t, func() {
		email := "email.email@ons.gov.uk"
		emailResponse := IsAllowedEmailDomain(email, allowedDomains)
		So(emailResponse, ShouldBeTrue)

		email = "email.email@ext.ons.gov.uk"
		emailResponse = IsAllowedEmailDomain(email, allowedDomains)
		So(emailResponse, ShouldBeTrue)
	})

	Convey("A new domain is added to the allowed list, emails using it are now validated", t, func() {
		email := "email.email@gmail.com"
		emailResponse := IsAllowedEmailDomain(email, allowedDomains)
		So(emailResponse, ShouldBeFalse)

		allowedDomains = append(allowedDomains, "@gmail.com")

		emailResponse = IsAllowedEmailDomain(email, allowedDomains)
		So(emailResponse, ShouldBeTrue)
	})
}
