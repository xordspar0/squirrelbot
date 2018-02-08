package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type update struct {
	ID      int `json:"update_id"`
	Message Message
}

type Message struct {
	ID   int `json:"message_id"`
	Date int64
	Chat struct {
		ID int
	}
	From struct {
		ID        int
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string
	}
	Text string
}

// SetWebhook establishes a connection with the Telegram server and tells
// Telegram where to send updates and messages.
func SetWebhook(address, token string) error {
	// Form request JSON.
	reqMap := make(map[string]string)
	reqMap["url"] = address
	reqJson, err := json.Marshal(reqMap)
	if err != nil {
		return errors.New("Failed to connect to Telegram API: " + err.Error())
	}

	// Make the request.
	resp, err := http.Post(
		fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook", token),
		"application/json",
		bytes.NewReader(reqJson),
	)

	if err != nil {
		return errors.New("Failed to connect to Telegram API: " + err.Error())
	}
	if resp.StatusCode != 200 {
		message, _ := ioutil.ReadAll(resp.Body)
		return errors.New("Failed to connect to Telegram API: " + resp.Status + " " + string(message))
	}

	return nil
}

func SendMessage(recipient int, messageBody, token string) error {
	// Form request JSON.
	reqMap := make(map[string]interface{})
	reqMap["chat_id"] = recipient
	reqMap["text"] = messageBody
	reqJson, err := json.Marshal(reqMap)
	if err != nil {
		return errors.New("Failed to send message: " + err.Error())
	}

	// Send the message.
	resp, err := http.Post(
		fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token),
		"application/json",
		bytes.NewReader(reqJson),
	)

	if err != nil {
		return errors.New("Failed to send message: " + err.Error())
	}
	if resp.StatusCode != 200 {
		message, _ := ioutil.ReadAll(resp.Body)
		return errors.New("Failed to send message: " + resp.Status + " " + string(message))
	}

	return nil
}

// GetMessage accepts json as a byte array and returns a message object if the
// json contains one.
func GetMessage(raw []byte) (Message, error) {
	var update update
	err := json.Unmarshal(raw, &update)
	if err != nil {
		return Message{}, errors.New("Error parsing message: " + err.Error())
	}

	return update.Message, nil
}
