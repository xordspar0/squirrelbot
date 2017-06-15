package bot

import (
	"github.com/xordspar0/squirrelbot/telegram"
	"github.com/xordspar0/squirrelbot/youtubedl"

	"fmt"
)

// handleYoutube takes Youtube url strings, downloads them, and sends a message
// back to the user.
func handleYoutube(message telegram.Message, url, token string) error {
	recipient, err := telegram.GetSenderID(message)
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

func handleLink(message telegram.Message, url, token string) error {
	recipient, err := telegram.GetSenderID(message)
	if err != nil {
		return err
	}

	return telegram.SendMessage(
		recipient,
		"This link does not look like a video. Stashing ordinary links is not yet implemented.",
		token,
	)
}
