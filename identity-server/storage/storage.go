package storage

import (
	"github.com/dapper-labs/identity-server/config"
	"github.com/dapper-labs/identity-server/model"
)

type UserRepository interface {
	Insert(*model.User) error
	GetUsers() (*[]model.User, error)
	GetUserWithEmail(email string) (*model.User, error)
	UpdateUser() error
}

func NewUserRepository(config *config.Config) (UserRepository, error) {
	if config.Driver == "memory" {
		return createInMemoryRepo()
	} else {
		return createDBRepo(config)
	}
}

func createInMemoryRepo() (UserRepository, error) {
	return &inMemoryUserRepository{users: make(map[string]model.User)}, nil
}

func createDBRepo(config *config.Config) (UserRepository, error) {
	return nil, nil
}
