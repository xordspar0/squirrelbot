package bot

import (
	"neolog.xyz/squirrelbot/telegram"

	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"time"
)

// handleYoutube takes Youtube url strings, downloads them, and sends a message
// back to the user.
func handleYoutube(message map[string]interface{}, url, token string) error {
	recipient, err := getID(message)
	if err != nil {
		return err
	}

	// Make sure that youtube-dl is installed.
	_, err = exec.LookPath("youtube-dl")
	if err != nil {
		_ = telegram.SendMessage(
			recipient,
			"Error: The server does not have youtube-dl installed on it.",
			token,
		)
		return err
	}

	// youtube-dl downloads the video for us.
	timestamp := time.Now().Local().Format(time.RFC3339)
	cmd := exec.Command(
		"youtube-dl",
		"--write-thumbnail",
		"--output",
		fmt.Sprintf("%s %s.%s", timestamp, "%(title)s", "%(ext)s"),
		url,
	)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	// Catch any error messgages while the process is running.
	errMessages, err := ioutil.ReadAll(stderr)
	if err != nil {
		return err
	}
	err = cmd.Wait()

	// If there is an error, log the standard error.
	if err != nil {
		return errors.New("Failed to download video:\n" + string(errMessages))
	}

	// Finally, send a message back to the user.
	err = telegram.SendMessage(
		recipient,
		"I saved your Youtube video.",
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
