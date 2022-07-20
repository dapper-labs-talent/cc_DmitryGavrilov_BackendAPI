package storage

import (
	"errors"

	"github.com/dapper-labs/identity-server/model"
)

var (
	ErrorUserInvalid = errors.New("user invalid")
	ErrorEmailExist  = errors.New("email was already used")
)

type inMemoryUserRepository struct {
	users map[string]model.User
}

func (r *inMemoryUserRepository) Insert(user *model.User) error {
	if user == nil {
		return ErrorUserInvalid
	}

	_, ok := r.users[user.Email]
	if ok {
		return ErrorEmailExist
	}

	r.users[user.Email] = *user
	return nil
}

func (r *inMemoryUserRepository) GetUsers() (*[]model.User, error) {
	return nil, nil
}

func (r *inMemoryUserRepository) UpdateUser() error {
	return nil
}
