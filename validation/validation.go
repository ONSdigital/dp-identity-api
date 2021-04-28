package validation

import (
	"regexp"
)

//Description of regex
//Before @: the initial squared brackets matches one the characters the email can start with, there is no full stop in these characters. The following statement in the round brackets checks if two full stops are consecutive. The next statement after the or (|) but before the @ checks for the characters that match these symbols.
//@ matches an @ symbol
//After @: the domain name is matched on the characters in the square brackets and the \. checks if there are two consecutive full stops.
var emailRegex = regexp.MustCompile(`^(?:[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[a-z0-9-!#$%&'*+/=?^_` + "`" + `{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)$`)

//IsEmailValid is a function to validate email addresses. Here an email is valid if it follows the standard rules for valid email addresses. The unit tests contain examples of what are considered valid or invalid email addresses. Here we do not validate on the domain name.
func IsEmailValid(e string) bool {

	minimumEmailLength := 3
	maximumEmailLength := 254

	if len(e) < minimumEmailLength && len(e) > maximumEmailLength {
		return false
	}

	return emailRegex.MatchString(e)
}
