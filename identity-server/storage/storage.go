package storage

import (
	"github.com/dapper-labs/identity-server/model"
)

type UserRepository interface {
	Insert(*model.User) error
	GetUsers() (*[]model.User, error)
	UpdateUser() error
}
