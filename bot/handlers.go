package bot

import (
	"github.com/xordspar0/squirrelbot/article"
	"github.com/xordspar0/squirrelbot/telegram"
	"github.com/xordspar0/squirrelbot/video"

	log "github.com/sirupsen/logrus"

	"fmt"
)

// handleYoutube takes Youtube url strings, downloads them, and sends a message
// back to the user.
func handleYoutube(url, directory string, recipient int, token string) error {
	// Get the video metadata.
	v, err := video.NewVideo(url)

	var videoTitle string
	if v.Title != "" {
		videoTitle = v.Title
	} else {
		videoTitle = "that video"
	}

	if err != nil {
		goto videoFail
	}

	err = v.WriteVideo(directory)
	if err != nil {
		goto videoFail
	}

	err = v.WriteThumb(directory)
	if err != nil {
		telegram.SendMessage(
			recipient,
			fmt.Sprintf("I couldn't save a thumbnail for \"%s\".", videoTitle),
			token,
		)
		log.Error(err.Error())
	}

	err = v.WriteNfo(directory)
	if err != nil {
		telegram.SendMessage(
			recipient,
			fmt.Sprintf("I couldn't save the metadata for \"%s\".", videoTitle),
			token,
		)
		log.Error(err.Error())
	}

	// Finally, report a successful save to the user.
	err = telegram.SendMessage(
		recipient,
		fmt.Sprintf("I saved your Youtube video, \"%s\".", videoTitle),
		token,
	)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil

videoFail:
	var message string
	if v.Title != "" {
		message = fmt.Sprintf("I couldn't save your video, \"%s\".", v.Title)
	} else {
		message = "I couldn't save that video."
	}

	telegram.SendMessage(recipient, message, token)
	log.Error(err.Error())
	return err
}

func handleLink(message telegram.Message, url string, recipient int, token, pocketKey string) error {
	a := article.NewArticle(url)

	err := a.Save(pocketKey)
	if err != nil {
		var message string
		if a.Title != "" {
			message = fmt.Sprintf("I couldn't save your article, \"%s\".", a.Title)
		} else {
			message = "I couldn't save that article."
		}

		telegram.SendMessage(recipient, message, token)
		log.Error(err.Error())
		return err
	}

	// Finally, report a successful save to the user.
	err = telegram.SendMessage(
		recipient,
		fmt.Sprintf("I saved your article, \"%s\".", a.Title),
		token,
	)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
