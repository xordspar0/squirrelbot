package bot

import (
	"github.com/xordspar0/squirrelbot/telegram"

	"gopkg.in/yaml.v2"
	"mvdan.cc/xurls"

	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type BotServer struct {
	Address         string `yaml:"address"`
	Port            string `yaml:"port"`
	Token           string `yaml:"token"`
	Directory       string `yaml:"directory"`
	Endpoint        string
	pocketKey       string `yaml:"pocket_key"`
	pocketUserToken string `yaml:"pocket_user_token"`
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

	var jsonBody map[string]interface{}
	err = json.Unmarshal(rawBody, &jsonBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	// TODO: Make this more safe when there are missing fields.
	if message, ok := jsonBody["message"].(map[string]interface{}); ok {
		messageText := message["text"].(string)
		if url := xurls.Strict.FindString(messageText); url != "" {
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
				err := handleYoutube(message, url, b.Directory, b.Token)
				if err != nil {
					log.Printf(err.Error())
				}
			} else {
				err := handleLink(message, url, b.Token, b.pocketKey, b.pocketUserToken)
				if err != nil {
					log.Printf(err.Error())
				}
			}
		}
	} else {
		log.Printf("Update %d has no message", jsonBody["update_id"])
	}
}
