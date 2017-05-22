package main

import (
	"neolog.xyz/squirrelbot/bot"

	"github.com/urfave/cli"

	"errors"
	"fmt"
	"log"
	"os"
)

var botname = "squirrelbot"
var version = "devel"
var app *cli.App

func main() {
	app := cli.NewApp()
	app.Name = botname
	app.Usage = "A Telegram bot that stashes away links that you send it so " +
		"that you can view them later."
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "server-name",
			Usage:  "The domain name of the server where this bot can be reached",
			EnvVar: "SERVER_NAME",
		},
		cli.IntFlag{
			Name:   "port",
			Value:  1327,
			Usage:  "The port to run the server on",
			EnvVar: "BOT_PORT",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "The authentication token for the Telegram API",
			EnvVar: "TELEGRAM_TOKEN",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err.Error())
	}
}

func run(c *cli.Context) error {
	config := &bot.ServerConfig{
		Name:     c.String("server-name"),
		Endpoint: fmt.Sprintf("/%s_%s/", botname, c.String("token")),
		Port:     c.String("port"),
		Token:    c.String("token"),
	}

	if config.Name == "" {
		return errors.New("Server domain name is not set")
	}

	return bot.Start(config)
}
