package telegram

import (
	"neolog.xyz/squirrelbot/config"

	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// SetWebhook establishes a connection with the Telegram server and tells
// Telegram where to send updates and messages.
func SetWebhook(c *config.ServerConfig, address string) error {
	// Form request JSON.
	reqMap := make(map[string]string)
	reqMap["url"] = address
	reqJson, err := json.Marshal(reqMap)
	if err != nil {
		return errors.New("Failed to connect to Telegram API: " + err.Error())
	}

	// Make the request.
	resp, err := http.Post(
		fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook", c.Token),
		"application/json",
		bytes.NewReader(reqJson),
	)

	if err != nil {
		return errors.New("Failed to connect to Telegram API: " + err.Error())
	}
	if resp.StatusCode == 200 {
		log.Println("Successfully connected to Telegram")
	} else {
		message, _ := ioutil.ReadAll(resp.Body)
		return errors.New("Failed to connect to Telegram API: " + resp.Status + " " + string(message))
	}

	return nil
}

func SendMessage() {
}
