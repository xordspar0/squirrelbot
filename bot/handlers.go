package bot

import (
	"github.com/xordspar0/squirrelbot/article"
	"github.com/xordspar0/squirrelbot/telegram"
	"github.com/xordspar0/squirrelbot/video"

	"fmt"
	"log"
)

// handleYoutube takes Youtube url strings, downloads them, and sends a message
// back to the user.
func handleYoutube(message telegram.Message, url, directory, token string) error {
	recipient, err := telegram.GetSenderID(message)
	if err != nil {
		return err
	}

	// Get the video metadata.
	v := video.NewVideo(url)

	var videoTitle string
	if v.Title != "" {
		videoTitle = v.Title
	} else {
		videoTitle = "that video"
	}

	err = v.WriteVideo(directory)
	if err != nil {
		var message string
		if v.Title != "" {
			message = fmt.Sprintf("I couldn't save your video, \"%s\".", v.Title)
		} else {
			message = "I couldn't save that video."
		}

		telegram.SendMessage(recipient, message, token)

		return err
	}

	err = v.WriteThumb(directory)
	if err != nil {
		telegram.SendMessage(
			recipient,
			fmt.Sprintf("I couldn't save a thumbnail for \"%s\".", videoTitle),
			token,
		)
		log.Println(err.Error())
	}

	err = v.WriteNfo(directory)
	if err != nil {
		telegram.SendMessage(
			recipient,
			fmt.Sprintf("I couldn't save the metadata for \"%s\".", videoTitle),
			token,
		)
		log.Println(err.Error())
	}

	// Finally, report a successful save to the user.
	err = telegram.SendMessage(
		recipient,
		fmt.Sprintf("I saved your Youtube video, \"%s\".", videoTitle),
		token,
	)
	if err != nil {
		return err
	}

	return nil
}

func handleLink(message telegram.Message, url, token, pocketKey, pocketUserToken string) error {
	recipient, err := telegram.GetSenderID(message)
	if err != nil {
		return err
	}

	a := article.NewArticle(url)

	err = a.Save(pocketKey, pocketUserToken)
	if err != nil {
		var message string
		if a.Title != "" {
			message = fmt.Sprintf("I couldn't save your article, \"%s\".", a.Title)
		} else {
			message = "I couldn't save that article."
		}

		telegram.SendMessage(recipient, message, token)

		return err
	}

	// Finally, report a successful save to the user.
	err = telegram.SendMessage(
		recipient,
		fmt.Sprintf("I saved your article, \"%s\".", a.Title),
		token,
	)
	if err != nil {
		return err
	}

	return nil
}
