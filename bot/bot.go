package bot

import (
	"neolog.xyz/squirrelbot/config"
	"neolog.xyz/squirrelbot/telegram"

	"github.com/mvdan/xurls"

	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func Start(c *config.ServerConfig) error {
	log.Println("Setting up endpoint at " + c.Endpoint)
	http.HandleFunc(c.Endpoint, botListener)
	err := telegram.SetWebhook(c, c.Name+c.Endpoint)
	if err != nil {
		return err
	}

	return http.ListenAndServe(":"+c.Port, nil)
}

func botListener(w http.ResponseWriter, r *http.Request) {
	jsonBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var body map[string]interface{}
	err = json.Unmarshal(jsonBody, &body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	if message, ok := body["message"].(map[string]interface{}); ok {
		bodyText := message["text"].(string)
		url := xurls.Strict.FindString(bodyText)
		if url != "" {
			if len(url) > 23 && url[:23] == "https://www.youtube.com" {
				log.Println("Found this Youtube video: " + url)
				handleYoutube(url)
			} else {
				log.Println("Found this link: " + url)
				handleLink(url)
			}
		}
	} else {
		log.Printf("Update %d has no message", message["update_id"])
	}
}
