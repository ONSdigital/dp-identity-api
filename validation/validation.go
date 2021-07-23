package validation

import (
	"regexp"
	"strings"
)

// This regex checks entered email address and if match found, email address deemed valid, else invalid.
//
// Regex delimits entered email into its user and domain name sections using `@` character.
//
// Note:
// —————
// `(?:…)` - denotes a non-capturing group
// `[…]`   - denotes a character set
// `*`     - denotes zero or more matches
// `+`     - denotes one or more matches
// `?`     - denotes a non-greedy match
//
// User name section:
// ——————————————————
// Grouped(
//   first character set ensures that email starts with either a valid character (one or more)
//   Grouped(second character set ensures characters are present after an escaped `.` (one or more)) (zero or more)
//   OR
//   Grouped(third character set matching permitted ascii hex characters OR a fourth character set matching a different set of permitted ascii hex characters) (zero or more)
//)
//
// Domain name section:
// ————————————————————
// Grouped(
//   first character set ensures that email domain section starts with an alphanumeric character
//   Grouped(second character set ensures alphanumeric/hyphen (`-`) characters followed by an alphanumeric character with a non-greedy match to first escaped `.` encountered) (one or more)
//   Third character set ensures that character immediately after `.` is an alphanumeric character
//   Grouped(fourth character set ensures alphanumeric/hyphen (`-`) characters) (zero or more) followed by fifth character set, ensuring domain name section ends in an alphanumeric character)
//)
//
// Example:
// ————————
// my.name@myself.com - match found - valid
// .my.name@myself.com - match not found - invalid
// my.name@myself.com. - match not found - invalid
var emailRegex = regexp.MustCompile(`^(?:[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[a-z0-9-!#$%&'*+/=?^_` + "`" + `{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)$`)
var onsEmailRegex = regexp.MustCompile(`^(?:[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[a-z0-9-!#$%&'*+/=?^_` + "`" + `{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(:?ext\.)?ons.gov.uk$`)
var (
	minimumEmailLength, maximumEmailLength int = 3, 254
)

//IsEmailValid is a function to validate email addresses. Here an email is valid if it follows the standard rules for valid email addresses. The unit tests contain examples of what are considered valid or invalid email addresses. Here we do not validate on the domain name.
func IsEmailValid(e string) bool {
	if !emailLengthValid(len(e)) {
		return false
	}

	return emailRegex.MatchString(strings.ToLower(e))
}

// IsAllowedEmailDomain - validates email address is a valid email format and the domain is in the allowed list in config
func IsAllowedEmailDomain(email string, allowedDomains []string) bool {
	if isValidStructure := IsEmailValid(email); !isValidStructure {
		return isValidStructure
	}
	for _, domain := range allowedDomains {
		if strings.HasSuffix(email, domain) {
			return true
		}
	}
	return false
}

func emailLengthValid(l int) bool {
	if l < minimumEmailLength && l > maximumEmailLength {
		return false
	}
	return true
}

func IsPasswordValid(p string) bool {
	if len(p) == 0 {
		return false
	}
	return true
}
