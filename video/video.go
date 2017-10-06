package video

import (
	"github.com/xordspar0/squirrelbot/youtubedl"

	"encoding/xml"
	"fmt"
	"io/ioutil"
	"time"
)

type Video struct {
	Url          string
	XMLName      xml.Name `xml:"movie"`
	Title        string   `xml:"title"`
	Description  string   `xml:"plot"`
	DownloadDate string
	FileName     string
}

func NewVideo(url string) *Video {
	currentTime := time.Now().Local().Format(time.RFC3339)
	v := &Video{
		Url:          url,
		Title:        youtubedl.GetTitle(url),
		Description:  youtubedl.GetDescription(url),
		DownloadDate: currentTime,
		FileName:     fmt.Sprintf("%s %s", currentTime, youtubedl.GetTitleSafe(url)),
	}

	if v.Title == "" {
		v.Title = "that video"
	}

	return v
}

func (v *Video) WriteVideo(directory string) error {
	return youtubedl.DownloadTo(
		v.Url,
		fmt.Sprintf("%s/%s.%s", directory, v.FileName, "%(ext)s"),
	)
}

// WriteNfo saves a thumbnail file for the video.
func (v *Video) WriteThumb(directory string) error {
	return youtubedl.DownloadThumbnailTo(
		v.Url,
		fmt.Sprintf("%s/%s-thumb.%s", directory, v.FileName, "%(ext)s"),
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
		fmt.Sprintf("<!-- Created on %s by SquirrelBot -->", v.DownloadDate),
		nfoXml,
	))
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s.nfo", directory, v.FileName), nfoXml, 0644)
	if err != nil {
		return err
	}

	return nil
}
