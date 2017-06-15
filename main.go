package main

import (
	"github.com/xordspar0/squirrelbot/bot"

	"github.com/urfave/cli"

	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
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
			EnvVar: "SQUIRRELBOT_PORT",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "The authentication token for the Telegram API",
			EnvVar: "TELEGRAM_TOKEN",
		},
		cli.StringFlag{
			Name:   "dir, d",
			Usage:  "The directory to store downloaded files",
			EnvVar: "SQUIRRELBOT_DIR",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err.Error())
	}
}

func run(c *cli.Context) error {
	// Change to the specified directory so that downloads will go there.
	os.Chdir(c.String("dir"))

	// Generate a random secret for the webhook endpoint. If the endpoint is a
	// secret between Telegram and the bot, we can be sure that requests to this
	// port are from Telegram.
	var max big.Int
	max.Exp(big.NewInt(2), big.NewInt(128), nil)
	randomSecret, err := rand.Int(rand.Reader, &max)
	if err != nil {
		return err
	}

	squirrelbotServer := &bot.BotServer{
		Name:     c.String("server-name"),
		Endpoint: fmt.Sprintf("/%s_%x/", botname, randomSecret),
		Port:     c.String("port"),
		Token:    c.String("token"),
	}

	if squirrelbotServer.Name == "" {
		return errors.New("Server domain name is not set")
	}

	return squirrelbotServer.Start()
}
