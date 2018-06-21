package cmd

import (
	"github.com/xordspar0/squirrelbot/bot"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"crypto/rand"
	"fmt"
	"math/big"
	"os"
)

var botname = "squirrelbot"
var version = "devel"
var systemConfigFile = ""

var cmd = &cobra.Command{
	Use:   botname,
	Short: botname + " is a Telegram bot that stashes videos that you send it",
	Long: botname + " is a Telegram bot that stashes videos. You can send it " +
		"links to Youtube videos and it will save them for you to view " +
		"later, complete with thumbnails, metadata, subtitles, etc.",
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		var missingParameters int

		if viper.GetBool("debug") {
			log.SetLevel(log.DebugLevel)
		}

		squirrelbotServer := &bot.Server{
			Address:   viper.GetString("address"),
			Port:      viper.GetString("port"),
			Token:     viper.GetString("telegram-token"),
			Directory: viper.GetString("dir"),
			Motd:      viper.GetString("motd"),
		}

		// Generate a random secret for the webhook endpoint. If the endpoint is a
		// secret between Telegram and the bot, we can be sure that requests to this
		// port are from Telegram.
		var max big.Int
		max.Exp(big.NewInt(2), big.NewInt(128), nil)
		randomSecret, err := rand.Int(rand.Reader, &max)
		if err != nil {
			goto exitWithError
		}
		squirrelbotServer.Endpoint = fmt.Sprintf("/%s_%x/", botname, randomSecret)

		if squirrelbotServer.Address == "" {
			log.Error("Server address is not set")
			missingParameters++
		}
		if squirrelbotServer.Token == "" {
			log.Error("Telegram API token is not set")
			missingParameters++
		}
		if missingParameters > 0 {
			log.Error("Missing required parameters")
			goto exitWithUsage
		}

		err = squirrelbotServer.Start()
		if err != nil {
			goto exitWithError
		}

		return

	exitWithError:
		log.Error(err)
		os.Exit(1)
		return

	exitWithUsage:
		cmd.Usage()
		os.Exit(1)
		return
	},
}

var address string

func init() {
	cmd.PersistentFlags().StringVar(
		&address,
		"address",
		"",
		"The address of the server where this bot can be reached (required)",
	)
	cmd.PersistentFlags().String(
		"telegram-token",
		"",
		"The authentication token for the Telegram API (required)",
	)
	cmd.PersistentFlags().IntP(
		"port",
		"p",
		80,
		"The port to run the server on",
	)
	cmd.PersistentFlags().StringP(
		"dir",
		"d",
		"",
		"The directory to store downloaded files",
	)
	cmd.PersistentFlags().String(
		"config",
		"",
		"The location of the server config file to use",
	)
	cmd.PersistentFlags().StringP(
		"motd",
		"m",
		"",
		"A message of the day to send to new users",
	)
	cmd.PersistentFlags().String(
		"debug",
		"",
		"Enable debug messages",
	)

	// Get configuration from flags, environment variables, and configuration
	// files.
	viper.BindPFlags(cmd.PersistentFlags())

	viper.SetEnvPrefix(botname)
	viper.AutomaticEnv()

	if userConfigFile := viper.GetString("config"); userConfigFile != "" {
		viper.SetConfigFile(userConfigFile)
	} else if systemConfigFile != "" {
		viper.SetConfigFile(systemConfigFile)
	}
}

func Execute() {
	cmd.Execute()
}
