package main

import (
	"github.com/xordspar0/squirrelbot/bot"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
)

var botname = "squirrelbot"
var version = "devel"
var systemConfigFile = ""
var app *cli.App

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	app := cli.NewApp()
	app.Name = botname
	app.Usage = "A Telegram bot that stashes away links that you send it so " +
		"that you can view them later."
	app.Action = run
	app.Version = version
	app.HideHelp = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "address",
			Usage:  "The address of the server where this bot can be reached (required)",
			EnvVar: "SQUIRRELBOT_ADDRESS",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "The authentication token for the Telegram API (required)",
			EnvVar: "TELEGRAM_TOKEN",
		},
		cli.IntFlag{
			Name:   "port",
			Value:  1327,
			Usage:  "The port to run the server on",
			EnvVar: "SQUIRRELBOT_PORT",
		},
		cli.StringFlag{
			Name:   "dir, d",
			Usage:  "The directory to store downloaded files",
			EnvVar: "SQUIRRELBOT_DIR",
		},
		cli.StringFlag{
			Name:   "config",
			Usage:  "The location of the server config file to use",
			EnvVar: "SQUIRRELBOT_CONFIG_FILE",
		},
		cli.StringFlag{
			Name:   "motd, m",
			Usage:  "A message of the day to send to new users",
			EnvVar: "SQUIRRELBOT_MOTD",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err.Error())
	}
}

func run(c *cli.Context) error {
	// Generate a random secret for the webhook endpoint. If the endpoint is a
	// secret between Telegram and the bot, we can be sure that requests to this
	// port are from Telegram.
	var max big.Int
	max.Exp(big.NewInt(2), big.NewInt(128), nil)
	randomSecret, err := rand.Int(rand.Reader, &max)
	if err != nil {
		return err
	}

	squirrelbotServer := &bot.BotServer{}

	// Load config from a file first.
	if fileName := c.String("config"); fileName != "" {
		if err := squirrelbotServer.LoadConfigFromFile(fileName); err != nil {
			return errors.New("Could not load config file: " + err.Error())
		}
	} else if systemConfigFile != "" {
		if err := squirrelbotServer.LoadConfigFromFile(systemConfigFile); err != nil {
			log.Printf("Could not load system config file %s", err.Error())
		}
	}

	// Add in settings from the command line.
	if address := c.String("address"); address != "" {
		squirrelbotServer.Address = address
	}
	if port := c.String("port"); port != "" {
		squirrelbotServer.Port = port
	}
	if token := c.String("token"); token != "" {
		squirrelbotServer.Token = token
	}
	if directory := c.String("dir"); directory != "" {
		squirrelbotServer.Directory = directory
	}
	if motd := c.String("motd"); motd != "" {
		squirrelbotServer.Motd = motd
	}

	squirrelbotServer.Endpoint = fmt.Sprintf("/%s_%x/", botname, randomSecret)

	if squirrelbotServer.Address == "" {
		return errors.New("Server address is not set")
	}
	if squirrelbotServer.Token == "" {
		return errors.New("Telegram API token is not set")
	}

	return squirrelbotServer.Start()
}
