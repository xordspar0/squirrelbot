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

// The Server contains all of the settings related to a running instance of the
// bot. You can set the fields individually or load them from a YAML file.
type Server struct {
	Address   string `yaml:"address"`
	Port      string `yaml:"port"`
	Token     string `yaml:"token"`
	Directory string `yaml:"directory"`
	Motd      string `yaml:"motd"`
	Endpoint  string
}

// LoadConfigFromFile fills in settings for the server from a YAML file,
// specified by a filename.
func (s *Server) LoadConfigFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, s)
	if err != nil {
		return err
	}

	return nil
}

// Start starts the bot as an HTTP server. It will listen on the port configured
// in s.Port.
func (s *Server) Start() error {
	log.WithFields(log.Fields{
		"url": s.Address + s.Endpoint,
	}).Info("Setting up endpoint")
	http.HandleFunc(s.Endpoint, s.botListener)
	err := telegram.SetWebhook(s.Address+s.Endpoint, s.Token)
	if err != nil {
		return err
	}

	return http.ListenAndServe(":"+s.Port, nil)
}

func (s *Server) botListener(w http.ResponseWriter, r *http.Request) {
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
			if s.Motd != "" {
				err = s.sendMotd(message.Chat.ID)
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
				go s.handleYoutube(url, s.Directory, message.Chat.ID)
			} else {
				go s.handleUnknown(message.Chat.ID)
			}
		}
	}
}

func (s *Server) sendMotd(recipient int) error {
	return telegram.SendMessage(
		recipient,
		s.Motd,
		s.Token,
	)
}
