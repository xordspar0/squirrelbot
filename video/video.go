package video

import (
	"encoding/xml"
)

type Video struct {
	XMLName     xml.Name `xml:"movie"`
	Title       string   `xml:"title"`
	Description string   `xml:"plot"`
}

func NewVideo(title, description string) *Video {
	newVideo := &Video{
		Title:       title,
		Description: description,
	}

	if newVideo.Title == "" {
		newVideo.Title = "that video"
	}

	return newVideo
}
