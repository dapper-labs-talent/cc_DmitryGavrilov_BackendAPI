package main

import (
	"github.com/dapper-labs/identity-server/api"
	"github.com/dapper-labs/identity-server/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func main() {
	config, err := config.LoadConfigWithPath("./config/config.ini")
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "cannot read configuration file"))
	}

	api, err := api.NewAPI(config)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "cannot create an api to server identity requests"))
	}

	err = api.ListenAndServe()
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "cannot start an api server"))
	}
}
