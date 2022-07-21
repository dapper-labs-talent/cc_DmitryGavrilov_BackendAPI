package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int64     `json:"-" pg:",pk"`
	Email        string    `sql:"unique:email"`
	PasswordHash string    `pg:"password_hash" json:"-"`
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	CreatedAt    time.Time `sql:"default:now()" json:"-"`
	UpdatedAt    time.Time `sql:"default:now()" json:"-"`
}

type UpdateUser struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func (u *User) Challenge(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func NewUser(firstname string, lastname string, email string, password string) (*User, error) {
	ph, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Firstname:    firstname,
		Lastname:     lastname,
		Email:        email,
		PasswordHash: ph,
	}
	return user, nil
}
func hashPassword(password string) (string, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(pw), nil
}
