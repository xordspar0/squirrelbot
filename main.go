package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type serverConfig struct {
	name string
	endpoint string
}

var botname = "squirrelbot"
var config *serverConfig

func main() {
	if err := run(config); err != nil {
		log.Fatalln(err.Error())
	}
}

func run(config *serverConfig) error {
	config.name = os.Getenv("SERVER_NAME")
	config.endpoint = fmt.Sprintf("/%s/", botname)

	if config.name == "" {
		return errors.New("Server domain name is not set")
	}

	return bot.Exec(config)
}
