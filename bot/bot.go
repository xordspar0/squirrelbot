package bot

import (
	"github.com/xordspar0/squirrelbot/telegram"

	"gopkg.in/yaml.v2"
	"mvdan.cc/xurls"

	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type BotServer struct {
	Address   string `yaml:"address"`
	Port      string `yaml:"port"`
	Token     string `yaml:"token"`
	Directory string `yaml:"directory"`
	Motd      string `yaml:"motd"`
	Endpoint  string
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
	log.Println("Setting up endpoint at " + b.Address + b.Endpoint)
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
	}

	message, err := telegram.GetMessage(rawBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
	}

	var recipient int
	if message.From.ID != 0 {
		recipient = message.From.ID
	} else if message.Chat.ID != 0 {
		recipient = message.Chat.ID
	}

	if message.Text != "" && recipient != 0 {
		if message.Text == "/start" {
			err = b.SendMotd(recipient)
			if err != nil {
				log.Printf(err.Error())
			}

		} else if url := xurls.Strict.FindString(message.Text); url != "" {
			if strings.HasPrefix(url, "http://www.youtube.com") ||
				strings.HasPrefix(url, "https://www.youtube.com") ||
				strings.HasPrefix(url, "http://m.youtube.com") ||
				strings.HasPrefix(url, "https://m.youtube.com") ||
				strings.HasPrefix(url, "http://youtu.be") ||
				strings.HasPrefix(url, "https://youtu.be") ||
				strings.HasPrefix(url, "http://vimeo.com") ||
				strings.HasPrefix(url, "https://vimeo.com") ||
				strings.HasPrefix(url, "http://player.vimeo.com") ||
				strings.HasPrefix(url, "https://player.vimeo.com") {
				err = handleYoutube(url, b.Directory, recipient, b.Token)
				if err != nil {
					log.Printf(err.Error())
				}
			} else {
				err = handleLink(url, recipient, b.Token)
				if err != nil {
					log.Printf(err.Error())
				}
			}
		}
	} else {
		log.Printf("Update %d has no message")
	}
}

func (b *BotServer) SendMotd(recipient int) error {
	return telegram.SendMessage(
		recipient,
		b.Motd,
		b.Token,
	)
}
