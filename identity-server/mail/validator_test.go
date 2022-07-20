package mail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailValidator_CorrectEmail_Success(t *testing.T) {
	email := "hello@gmail.com"
	err := Validate(email)
	assert.Nil(t, err, "error expected to be nil")
}

func TestEmailValidator_IncorrectEmail_ErrorExpected(t *testing.T) {
	email := "hellogmail.com"
	err := Validate(email)
	assert.NotNil(t, err, "an error expected to be thrown")
	assert.Equal(t, BadFormatError, err)
}

func TestEmailValidator_EmptyEmail_ErrorExpected(t *testing.T) {
	email := ""
	err := Validate(email)
	assert.NotNil(t, err, "an error expected to be thrown")
	assert.Equal(t, BadFormatError, err)
}
