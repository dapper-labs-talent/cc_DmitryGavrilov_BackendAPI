package main

import (
	"fmt"

	"github.com/dapper-labs/identity-server/config"
)

func main() {
	c, err := config.LoadConfigWithPath("")

	if c != nil {

	}
	if err != nil {
		fmt.Println(err.Error())

	}
}
