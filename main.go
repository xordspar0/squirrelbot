package main

import (
	"neolog.xyz/squirrelbot/bot"
	"neolog.xyz/squirrelbot/config"

	"errors"
	"fmt"
	"log"
	"os"
)

var botname = "squirrelbot"

func main() {
	var c *config.ServerConfig
	if err := run(c); err != nil {
		log.Fatalln(err.Error())
	}
}

func run(c *config.ServerConfig) error {
	c.Name = os.Getenv("SERVER_NAME")
	c.Endpoint = fmt.Sprintf("/%s/", botname)

	if c.Name == "" {
		return errors.New("Server domain name is not set")
	}

	return bot.Exec(c)
}
