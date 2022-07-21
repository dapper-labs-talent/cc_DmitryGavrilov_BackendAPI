package storage

import (
	"errors"
	"fmt"

	"github.com/dapper-labs/identity-server/config"
	"github.com/dapper-labs/identity-server/model"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/sirupsen/logrus"
)

var (
	ErrorDBCannotConnect  = errors.New("could not connect to remote database")
	ErrorDBConnection     = errors.New("the db connection is in invalid state")
	ErrorInsertUser       = errors.New("could not save user")
	ErrorSelectUsers      = errors.New("could not select all the users")
	ErrorConfigInvalid    = errors.New("configuration file is nil")
	ErrorUserAlreadyExist = errors.New("user with the same email address already exist")
)

type postgresMigrator struct {
	conn *pg.DB
}

type postgresUserRepository struct {
	conn *pg.DB
}

func NewPostgresUserRepository(config *config.Config) (*postgresUserRepository, error) {
	conn, err := getPostgresConnection(config)
	if err != nil {
		return nil, err
	}

	return &postgresUserRepository{conn: conn}, nil
}

func (pur *postgresUserRepository) Insert(user *model.User) error {
	if pur.conn == nil {
		return ErrorDBConnection
	}

	_, err := pur.conn.Model(user).Insert()
	if err != nil {
		logrus.Error(err)

		pgErr, ok := err.(pg.Error)
		if ok && pgErr.IntegrityViolation() {
			return ErrorUserAlreadyExist
		}
		return err
	}

	return nil
}
func (pur *postgresUserRepository) GetUsers() (*[]model.User, error) {
	if pur.conn == nil {
		return nil, ErrorDBConnection
	}

	var users []model.User
	err := pur.conn.Model(&users).Select()
	if err != nil {
		logrus.Error(err)
		return nil, ErrorSelectUsers
	}

	return &users, nil
}

func (pur *postgresUserRepository) GetUserWithEmail(email string) (*model.User, error) {
	user := &model.User{}

	err := pur.conn.Model(user).Where("email = ?", email).Select()
	if err != nil {
		logrus.Error(err)
		return nil, errors.New("an error ocurred while selecting by email")
	}
	return user, nil
}

func (pur *postgresUserRepository) UpdateUserWithEmail(updateUser *model.UpdateUser, email string) error {
	user := &model.User{}

	if updateUser.Firstname != "" {
		user.Firstname = updateUser.Firstname
	}

	if updateUser.Lastname != "" {
		user.Lastname = updateUser.Lastname
	}
	_, err := pur.conn.Model(user).Where("email = ?", email).UpdateNotNull()
	if err != nil {
		logrus.Error(err)
		return errors.New("an error ocurred while updating an user")
	}
	return nil
}

func (pm *postgresMigrator) CreateSchema() error {

	logrus.Info("Starting db schema creation")

	models := []interface{}{
		(*model.User)(nil),
	}
	defer pm.conn.Close()
	for _, model := range models {
		err := pm.conn.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}

	logrus.Info("Db schema created")
	return nil
}

func createPosgressMigrator(config *config.Config) (Migrator, error) {
	conn, err := getPostgresConnection(config)
	if err != nil {
		return nil, err
	}

	return &postgresMigrator{conn: conn}, nil
}

func getPostgresConnection(config *config.Config) (*pg.DB, error) {
	if config == nil {
		return nil, ErrorConfigInvalid
	}

	addr := fmt.Sprintf("%v:%d", config.Database.Host, config.Database.Port)
	logrus.Infof("Getting db connection at %s", addr)
	opts := &pg.Options{
		Database: config.DbName,
		Addr:     addr,
		User:     config.Database.User,
		Password: config.Database.Password,
	}

	db := pg.Connect(opts)
	if db == nil {
		return nil, ErrorDBCannotConnect
	}
	logrus.Info("Connection successfully established")

	return db, nil
}
