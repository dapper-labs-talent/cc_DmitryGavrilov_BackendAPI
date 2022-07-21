package storage

import (
	"errors"
	"fmt"

	"github.com/dapper-labs/identity-server/config"
	"github.com/dapper-labs/identity-server/model"
)

type UserRepository interface {
	Insert(*model.User) error
	GetUsers() (*[]model.User, error)
	GetUserWithEmail(email string) (*model.User, error)
	UpdateUserWithEmail(*model.UpdateUser, string) error
}

type Migrator interface {
	CreateSchema() error
}

func NewUserRepository(config *config.Config) (UserRepository, error) {
	if config.Driver == "memory" {
		return &inMemoryUserRepository{users: make(map[string]model.User)}, nil
	} else if config.Driver == "postgres" {
		return NewPostgresUserRepository(config)
	} else {
		return nil, errors.New(fmt.Sprintf("driver %s is not supported", config.Driver))
	}
}

func NewMigrator(config *config.Config) (Migrator, error) {
	return createPosgressMigrator(config)
}
