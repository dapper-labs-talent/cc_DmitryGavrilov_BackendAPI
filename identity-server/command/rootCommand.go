package command

import (
	"github.com/dapper-labs/identity-server/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCommand = cobra.Command{
	Use: "identity-server",
	Run: func(cmd *cobra.Command, args []string) {
		executeCommand(startServer)
	},
}

// configuration filename
var cfile = ""

func RootCommand() *cobra.Command {
	rootCommand.AddCommand(&serverCommand, &migrateCommand)
	rootCommand.PersistentFlags().StringVar(&cfile, "config", "", "the configuratin filename to use")
	return &rootCommand
}

func executeCommand(fh func(*config.Config)) {
	config, err := config.LoadConfigWithPath(cfile)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "cannot read configuration file"))
	}
	logrus.SetLevel(logrus.Level(config.LogLevel))

	fh(config)
}
