package bot

import (
	"github.com/xordspar0/squirrelbot/telegram"
	"github.com/xordspar0/squirrelbot/video"
	"github.com/xordspar0/squirrelbot/youtubedl"

	"encoding/xml"
	"fmt"
	"io/ioutil"
	"time"
)

// handleYoutube takes Youtube url strings, downloads them, and sends a message
// back to the user.
func handleYoutube(message telegram.Message, url, dir, token string) error {
	recipient, err := telegram.GetSenderID(message)
	if err != nil {
		return err
	}

	// Get the video metadata.
	newVideo := video.NewVideo(
		youtubedl.GetTitle(url),
		youtubedl.GetDescription(url),
	)
	timestamp := time.Now().Local().Format(time.RFC3339)
	newVideo.Comment = fmt.Sprintf("Created on %s by SquirrelBot", timestamp)

	// youtube-dl downloads the video and its thumbnail.
	err = youtubedl.Download(url, dir, timestamp+" ")
	if err != nil {
		telegram.SendMessage(
			recipient,
			fmt.Sprintf("I couldn't save \"%s\".", newVideo.Title),
			token,
		)
		return err
	}

	err = youtubedl.DownloadThumbnail(url, dir, timestamp+" ")
	if err != nil {
		telegram.SendMessage(
			recipient,
			fmt.Sprintf("I couldn't save a thumbnail for \"%s\".", newVideo.Title),
			token,
		)
		return err
	}

	// Make a .nfo file for the video.
	nfoXml, err := xml.MarshalIndent(newVideo, "", "    ")
	nfoXml = []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n" + string(nfoXml))
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s%s.nfo", dir, timestamp+" ", newVideo.Title), nfoXml, 0644)
	if err != nil {
		telegram.SendMessage(
			recipient,
			fmt.Sprintf("I couldn't save the metadata for \"%s\".", newVideo.Title),
			token,
		)
		return err
	}

	// Finally, send a message back to the user.
	err = telegram.SendMessage(
		recipient,
		fmt.Sprintf("I saved your Youtube video, \"%s\".", newVideo.Title),
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
