package bot

import (
	"neolog.xyz/squirrelbot/telegram"

	"github.com/mvdan/xurls"

	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type BotServer struct {
	Name     string
	Endpoint string
	Port     string
	Token    string
}

func (b *BotServer) Start() error {
	log.Println("Setting up endpoint at " + b.Name + b.Endpoint)
	http.HandleFunc(b.Endpoint, b.botListener)
	err := telegram.SetWebhook(b.Name+b.Endpoint, b.Token)
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

	if message, ok := jsonBody["message"].(map[string]interface{}); ok {
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
		log.Printf("Update %d has no message", jsonBody["update_id"])
	}
}
