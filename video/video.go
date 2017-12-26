package video

import (
	"github.com/xordspar0/squirrelbot/youtubedl"

	"encoding/xml"
	"fmt"
	"io/ioutil"
	"time"
)

type Video struct {
	XMLName      xml.Name `xml:"movie"`
	Url          string   `xml:"-"`
	Title        string   `xml:"title"`
	Description  string   `xml:"plot"`
	downloadDate time.Time
	fileName     string
}

func NewVideo(url string) *Video {
	currentTime := time.Now().Local()
	v := &Video{
		Url:          url,
		Title:        youtubedl.GetTitle(url),
		Description:  youtubedl.GetDescription(url),
		downloadDate: currentTime,
		fileName:     fmt.Sprintf("%s %s", currentTime.Format(time.RFC3339), youtubedl.GetTitleSafe(url)),
	}

	if v.Title == "" {
		v.Title = "that video"
	}

	return v
}

func (v *Video) WriteVideo(directory string) error {
	return youtubedl.DownloadTo(
		v.Url,
		fmt.Sprintf("%s/%s.%s", directory, v.fileName, "%(ext)s"),
	)
}

// WriteNfo saves a thumbnail file for the video.
func (v *Video) WriteThumb(directory string) error {
	return youtubedl.DownloadThumbnailTo(
		v.Url,
		fmt.Sprintf("%s/%s-thumb.%s", directory, v.fileName, "%(ext)s"),
	)
}

// WriteNfo makes a .nfo file for the video, which includes the video's Title
// and Description.
func (v *Video) WriteNfo(directory string) error {
	nfoXml, err := xml.MarshalIndent(v, "", "    ")
	if err != nil {
		return err
	}

	nfoXml = []byte(fmt.Sprintf("%s\n%s\n%s",
		`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`,
		fmt.Sprintf("<!-- Created on %s by SquirrelBot -->", v.downloadDate.Format(time.RFC3339)),
		nfoXml,
	))
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s.nfo", directory, v.fileName), nfoXml, 0644)
	if err != nil {
		return err
	}

	return nil
}
