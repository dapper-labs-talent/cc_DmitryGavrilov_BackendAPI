package main

import (
	"github.com/dapper-labs/identity-server/command"
	"github.com/sirupsen/logrus"
)

func main() {
	err := command.RootCommand().Execute()
	if err != nil {
		logrus.Fatal(err)
	}
}
