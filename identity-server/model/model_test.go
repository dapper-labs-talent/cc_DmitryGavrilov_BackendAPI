package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser_CorrectUser_Ok(t *testing.T) {

	user, err := NewUser("alex", "wolfgang", "alex.wolfgang@gmail.com", "wunderland")
	assert.Nil(t, err)

	assert.Equal(t, "alex", user.Firstname)
	assert.Equal(t, "wolfgang", user.Lastname)
	assert.Equal(t, "alex.wolfgang@gmail.com", user.Email)
	assert.NotEqual(t, "wunderland", user.PasswordHash)
}

func TestChallengePassword_Password_NotMatch(t *testing.T) {
	user, err := NewUser("alex", "wolfgang", "alex.wolfgang@gmail.com", "wunderland")
	assert.Nil(t, err)

	assert.False(t, user.Challenge("wrongpassword"))
}

func TestChallengePassword_Password_Match(t *testing.T) {
	user, err := NewUser("alex", "wolfgang", "alex.wolfgang@gmail.com", "wunderland")
	assert.Nil(t, err)

	assert.True(t, user.Challenge("wunderland"))
}
