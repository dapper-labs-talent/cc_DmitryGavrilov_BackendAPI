package command

import (
	"fmt"

	"github.com/dapper-labs/identity-server/api"
	"github.com/dapper-labs/identity-server/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serverCommand = cobra.Command{
	Use:  "server",
	Long: "Start Identity Server",
	Run: func(cmd *cobra.Command, args []string) {
		executeCommand(startServer)
	},
}

func startServer(config *config.Config) {
	logrus.Info("Starting server")
	api, err := api.NewAPI(config)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "cannot create an api to server identity requests"))
	}
	address := fmt.Sprintf("0.0.0.0:%d", config.ListenPort)

	err = api.ListenAndServe(address)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "cannot start an api server"))
	}
}
