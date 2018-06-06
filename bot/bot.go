package bot

import (
	"github.com/xordspar0/squirrelbot/telegram"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"mvdan.cc/xurls"

	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type BotServer struct {
	Address         string `yaml:"address"`
	Port            string `yaml:"port"`
	Token           string `yaml:"token"`
	Directory       string `yaml:"directory"`
	Motd            string `yaml:"motd"`
	Endpoint        string
	PocketKey       string `yaml:"pocket_key"`
	PocketUserToken string `yaml:"pocket_user_token"`
}

func (b *BotServer) LoadConfigFromFile(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, b)
	if err != nil {
		return err
	}

	return nil
}

func (b *BotServer) Start() error {
	log.WithFields(log.Fields{
		"url": b.Address + b.Endpoint,
	}).Info("Setting up endpoint")
	http.HandleFunc(b.Endpoint, b.botListener)
	err := telegram.SetWebhook(b.Address+b.Endpoint, b.Token)
	if err != nil {
		return err
	}

	return http.ListenAndServe(":"+b.Port, nil)
}

func (b *BotServer) botListener(w http.ResponseWriter, r *http.Request) {
	rawBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err.Error())
	}

	log.WithFields(log.Fields{
		"request body": string(rawBody),
	}).Debug("Received request")

	message, err := telegram.GetMessage(rawBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err.Error())
	}

	log.WithFields(log.Fields{
		"message ID": message.ID,
		"body":       message.Text,
		"from":       message.From.ID,
	}).Debug("Parsed message")

	var username string
	if message.From.Username != "" {
		username = message.From.Username
	} else if message.From.FirstName != "" || message.From.LastName != "" {
		username = strings.TrimSpace(fmt.Sprintf("%s %s",
			message.From.FirstName,
			message.From.LastName,
		))
	}

	if message.Text == "" {
		log.Error("Message has no body")
	} else if message.Chat.ID == 0 {
		log.Error("Message has no sender")
	} else {
		if message.Text == "/start" {
			if b.Motd != "" {
				err = b.SendMotd(message.Chat.ID)
				if err != nil {
					log.Error(err.Error())
				}
			}
		} else if url := xurls.Strict.FindString(message.Text); url != "" {
			infoLogger := log.WithFields(log.Fields{
				"url":        url,
				"user":       username,
				"message ID": message.ID,
			})

			if isYoutubeSource(url) {
				infoLogger.Info("Stashing video")
				handleYoutube(url, b.Directory, message.Chat.ID, b.Token)
			} else {
				infoLogger.Info("Stashing link")
				handleLink(message, url, message.Chat.ID, b.Token, b.PocketKey, b.PocketUserToken)
			}
		}
	}
}

func (b *BotServer) SendMotd(recipient int) error {
	return telegram.SendMessage(
		recipient,
		b.Motd,
		b.Token,
	)
}
