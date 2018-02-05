package bot

import (
	"github.com/xordspar0/squirrelbot/telegram"
	"github.com/xordspar0/squirrelbot/video"

	log "github.com/sirupsen/logrus"

	"fmt"
)

// handleYoutube takes Youtube url strings, downloads them, and sends a message
// back to the user.
func handleYoutube(url, directory string, recipient int, token string) error {
	// Get the video metadata.
	v := video.NewVideo(url)

	err := v.WriteVideo(directory)
	if err != nil {
		telegram.SendMessage(
			recipient,
			fmt.Sprintf("I couldn't save \"%s\".", v.Title),
			token,
		)
		log.Error(err.Error())
		return err
	}

	err = v.WriteThumb(directory)
	if err != nil {
		telegram.SendMessage(
			recipient,
			fmt.Sprintf("I couldn't save a thumbnail for \"%s\".", v.Title),
			token,
		)
		log.Error(err.Error())
	}

	err = v.WriteNfo(directory)
	if err != nil {
		telegram.SendMessage(
			recipient,
			fmt.Sprintf("I couldn't save the metadata for \"%s\".", v.Title),
			token,
		)
		log.Error(err.Error())
	}

	// Finally, send a message back to the user.
	err = telegram.SendMessage(
		recipient,
		fmt.Sprintf("I saved your Youtube video, \"%s\".", v.Title),
		token,
	)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func handleLink(url string, recipient int, token string) error {
	err := telegram.SendMessage(
		recipient,
		"This link does not look like a video. Stashing ordinary links is not yet implemented.",
		token,
	)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
