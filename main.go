package main

import (
	"neolog.xyz/squirrelbot/bot"

	"errors"
	"fmt"
	"log"
	"os"
)

var botname = "squirrelbot"

func main() {
	c := &bot.ServerConfig{}
	if err := run(c); err != nil {
		log.Fatalln(err.Error())
	}
}

func run(c *bot.ServerConfig) error {
	c.Name = os.Getenv("SERVER_NAME")
	c.Endpoint = fmt.Sprintf("/%s/", botname)

	if c.Name == "" {
		return errors.New("Server domain name is not set")
	}

	return bot.Exec(c)
}
