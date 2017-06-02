package bot

import (
	"neolog.xyz/squirrelbot/telegram"

	"errors"
	"fmt"
)

func handleYoutube(message map[string]interface{}, url, token string) error {
	recipient, err := getID(message)
	if err != nil {
		return err
	}

	fmt.Println(recipient) //debug

	return telegram.SendMessage(
		recipient,
		"I saved your Youtube video.",
		token,
	)
}

func handleLink(message map[string]interface{}, url, token string) error {
	recipient, err := getID(message)

	if err != nil {
		return err
	}

	return telegram.SendMessage(
		recipient,
		"This link does not look like a video. Stashing ordinary links is not yet implemented.",
		token,
	)
}

func getID(message map[string]interface{}) (string, error) {
	var id int

	if _, ok := message["from"]; ok {
		id = int(message["from"].(map[string]interface{})["id"].(float64))
	} else if _, ok := message["chat"]; ok {
		id = int(message["chat"].(map[string]interface{})["id"].(float64))
	} else {
		return "", errors.New("Error: message has no sender")
	}

	return fmt.Sprintf("%d", id), nil
}
