package bot

import (
	"github.com/xordspar0/squirrelbot/telegram"
	"github.com/xordspar0/squirrelbot/video"

	log "github.com/sirupsen/logrus"

	"fmt"
)

// handleYoutube takes Youtube url strings, downloads them, and sends a message
// back to the user.
func (s *Server) handleYoutube(url, directory string, recipient int) error {
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
			s.Token,
		)
		log.Error(err.Error())
	}

	err = v.WriteNfo(directory)
	if err != nil {
		telegram.SendMessage(
			recipient,
			fmt.Sprintf("I couldn't save the metadata for \"%s\".", videoTitle),
			s.Token,
		)
		log.Error(err.Error())
	}

	// Finally, report a successful save to the user.
	err = telegram.SendMessage(
		recipient,
		fmt.Sprintf("I saved your Youtube video, \"%s\".", videoTitle),
		s.Token,
	)
	if err != nil {
		log.Error(err.Error())
	}
	return err

videoFail:
	var message string
	if v.Title != "" {
		message = fmt.Sprintf("I couldn't save your video, \"%s\".", v.Title)
	} else {
		message = "I couldn't save that video."
	}

	telegram.SendMessage(recipient, message, s.Token)
	log.Error(err.Error())
	return err
}

func (s *Server) handleUnknown(recipient int) error {
	err := telegram.SendMessage(
		recipient,
		"That doesn't look like a video that I can save. Contact the developer"+
			"if you would like me to be able to save this type of video.",
		s.Token,
	)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}
