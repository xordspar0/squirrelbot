package bot

import (
	"neolog.xyz/squirrelbot/telegram"
	"neolog.xyz/squirrelbot/youtubedl"

	"errors"
	"fmt"
)

// handleYoutube takes Youtube url strings, downloads them, and sends a message
// back to the user.
func handleYoutube(message map[string]interface{}, url, token string) error {
	recipient, err := getID(message)
	if err != nil {
		return err
	}

	videoTitle := youtubedl.GetTitle(url)
	if videoTitle == "" {
		videoTitle = "that video"
	} else {
		videoTitle = "\"" + videoTitle + "\""
	}

	err = youtubedl.Download(url)

	// If there was an error, log the standard error and send a report to the
	// user.
	if err != nil {
		telegram.SendMessage(
			recipient,
			fmt.Sprintf("I couldn't save %s.", videoTitle),
			token,
		)
		return err
	}

	// Finally, send a message back to the user.
	err = telegram.SendMessage(
		recipient,
		fmt.Sprintf("I saved your Youtube video, %s.", videoTitle),
		token,
	)
	if err != nil {
		return err
	}

	return nil
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
