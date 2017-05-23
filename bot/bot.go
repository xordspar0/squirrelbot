package bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func Start(c *ServerConfig) error {
	log.Println("Setting up endpoint at " + c.Endpoint)
	http.HandleFunc(c.Endpoint, botListener)
	err := setWebhook(c, c.Name+c.Endpoint)
	if err != nil {
		return err
	}

	// Switch to this after I get the cert figured out.
	// Telegram requires https connections.
	//func ListenAndServeTLS(addr, certFile, keyFile string, handler Handler) error
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
		log.Println(message["text"].(string))
	} else {
		log.Printf("Update %d has no message", message["update_id"])
	}
}

// setWebhook establishes a connection with the Telegram server and tells
// Telegram where to send updates and messages.
func setWebhook(c *ServerConfig, address string) error {
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
