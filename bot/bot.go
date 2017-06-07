package bot

import (
	"neolog.xyz/squirrelbot/telegram"

	"github.com/mvdan/xurls"

	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Bot struct {
	Name     string
	Endpoint string
	Port     string
	Token    string
}

func (b *Bot) Start() error {
	log.Println("Setting up endpoint at " + b.Name + b.Endpoint)
	http.HandleFunc(b.Endpoint, b.botListener)
	err := telegram.SetWebhook(b.Name+b.Endpoint, b.Token)
	if err != nil {
		return err
	}

	return http.ListenAndServe(":"+b.Port, nil)
}

func (b *Bot) botListener(w http.ResponseWriter, r *http.Request) {
	rawBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var jsonBody map[string]interface{}
	err = json.Unmarshal(rawBody, &jsonBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	if message, ok := jsonBody["message"].(telegram.Message); ok {
		messageText := message["text"].(string)
		url := xurls.Strict.FindString(messageText)
		if url != "" {
			if len(url) > 23 && url[:23] == "https://www.youtube.com" {
				err := handleYoutube(message, url, b.Token)
				if err != nil {
					log.Printf(err.Error())
				}
			} else {
				err := handleLink(message, url, b.Token)
				if err != nil {
					log.Printf(err.Error())
				}
			}
		}
	} else {
		log.Printf("Update %d has no message", message["update_id"])
	}
}
