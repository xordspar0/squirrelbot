package bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func Start(c *ServerConfig) error {
	listenAddr := c.Name + c.Endpoint
	setWebhook(c, listenAddr)

	return http.ListenAndServe(":"+c.Port, nil)
}

func setWebhook(c *ServerConfig, address string) error {
	// Form request JSON.
	reqMap := make(map[string]string)
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

	if resp.StatusCode != 200 {
		return errors.New("Failed to connect to Telegram API: " + resp.Status)
	}

	return nil
}
