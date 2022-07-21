package command

import (
	"github.com/dapper-labs/identity-server/config"
	"github.com/dapper-labs/identity-server/storage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var migrateCommand = cobra.Command{
	Use: "migrate",
	Run: func(cmd *cobra.Command, args []string) {
		executeCommand(migrate)
	},
}

func MigrateCommand() *cobra.Command {
	return &migrateCommand
}

func migrate(config *config.Config) {
	logrus.Info("Starting database migration")

	migrator, err := storage.NewMigrator(config)
	if err != nil {
		logrus.Fatalf("could not create a new database migrator using provided configuration: %v", err)
	}

	err = migrator.CreateSchema()
	if err != nil {
		logrus.Fatalf("could not finish database migration %v", err)
	}

	logrus.Info("Database migration completed")
}
