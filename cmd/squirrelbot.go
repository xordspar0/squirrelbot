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
	"strings"
)

var botname = "squirrelbot"
var version = "devel"
var systemConfigFile = ""

var Cmd = &cobra.Command{
	Use:               botname,
	Short:             "a Telegram bot that stashes videos that you send it",
	Version:           version,
	DisableAutoGenTag: true,
	Long: strings.Title(botname) +
		` is a Telegram bot that stashes videos. You can send it links to
Youtube videos and it will save them for you to view later, complete with
thumbnails, metadata, subtitles, etc.

Configuration options may be specified as command line options, as environment
variables, or in a configuration file. The default configuration file can be
specified at compile time, but it is /etc/squirrelbot/config.yaml by default. A
different file may be specified with the --config command line option (see
OPTIONS).`,
	Run: func(Cmd *cobra.Command, args []string) {
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
		Cmd.Usage()
		os.Exit(1)
		return
	},
}

var address string

func init() {
	Cmd.PersistentFlags().StringVar(
		&address,
		"address",
		"",
		"The address of the server where this bot can be reached (required) ($SQUIRRELBOT_ADDRESS)",
	)
	Cmd.PersistentFlags().String(
		"telegram-token",
		"",
		"The authentication token for the Telegram API. You can find directions for obtaining your token at https://core.telegram.org/bots. (required) ($SQUIRRELBOT_TELEGRAM_TOKEN)",
	)
	Cmd.PersistentFlags().IntP(
		"port",
		"p",
		80,
		"The port to run the server on ($SQUIRRELBOT_PORT)",
	)
	Cmd.PersistentFlags().StringP(
		"dir",
		"d",
		"",
		"The directory to store downloaded files ($SQUIRRELBOT_DIR)",
	)
	Cmd.PersistentFlags().String(
		"config",
		"",
		"A configuration file to use ($SQUIRRELBOT_CONFIG)",
	)
	Cmd.PersistentFlags().StringP(
		"motd",
		"m",
		"",
		"A message of the day to send to new users ($SQUIRRELBOT_MOTD)",
	)
	Cmd.PersistentFlags().String(
		"debug",
		"",
		"Enable debug messages ($SQUIRRELBOT_DEBUG)",
	)

	// Get configuration from flags, environment variables, and configuration
	// files.
	viper.BindPFlags(Cmd.PersistentFlags())

	viper.SetEnvPrefix(botname)
	viper.AutomaticEnv()

	if userConfigFile := viper.GetString("config"); userConfigFile != "" {
		viper.SetConfigFile(userConfigFile)
	} else if systemConfigFile != "" {
		viper.SetConfigFile(systemConfigFile)
	}
}

func Execute() {
	Cmd.Execute()
}
