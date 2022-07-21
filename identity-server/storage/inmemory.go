package storage

import (
	"errors"

	"github.com/dapper-labs/identity-server/model"
)

var (
	ErrorUserInvalid            = errors.New("user invalid")
	ErrorUserNotFound           = errors.New("requesting user was not found")
	ErrorEmailExist             = errors.New("email was already used")
	ErrorRepositoryInvalidState = errors.New("user repository has an invalid state, please use a NewUserRepository function to create it")
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
	res := make([]model.User, 0)
	for _, user := range r.users {
		res = append(res, user)
	}

	return &res, nil
}

func (r *inMemoryUserRepository) GetUserWithEmail(email string) (*model.User, error) {
	if r.users == nil {
		return nil, ErrorRepositoryInvalidState
	}

	user, _ := r.users[email]
	return &user, nil
}

func (r *inMemoryUserRepository) UpdateUserWithEmail(updateUser *model.UpdateUser, email string) error {
	if r.users == nil {
		return ErrorRepositoryInvalidState
	}

	user, ok := r.users[email]
	if !ok {
		return ErrorUserNotFound
	}

	if updateUser.Firstname != "" {
		user.Firstname = updateUser.Firstname
	}

	if updateUser.Lastname != "" {
		user.Lastname = updateUser.Lastname
	}

	r.users[email] = user
	return nil
}
