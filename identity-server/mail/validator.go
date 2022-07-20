package mail

import (
	"errors"
	"regexp"
)

var (
	BadFormatError = errors.New("email has an invalid format")
	//https://www.w3.org/TR/html5/forms.html#valid-e-mail-address
	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func Validate(email string) error {
	if !emailRegex.MatchString(email) {
		return BadFormatError
	}

	return nil
}
