package video

import (
	"github.com/xordspar0/squirrelbot/youtubedl"

	"encoding/xml"
	"fmt"
	"io/ioutil"
	"time"
)

// Video contains the metadata of a video and all the information required to
// download it.
type Video struct {
	XMLName      xml.Name `xml:"movie"`
	URL          string   `xml:"-"`
	Title        string   `xml:"title"`
	Description  string   `xml:"plot"`
	downloadDate time.Time
	filename     string
}

// NewVideo takes a URL for a video and discovers the metadata for it using
// youtube-dl.
func NewVideo(url string) (*Video, error) {
	currentTime := time.Now().Local()
	title, err := youtubedl.GetTitle(url)
	if err != nil {
		return &Video{}, err
	}

	description, err := youtubedl.GetDescription(url)
	if err != nil {
		return &Video{}, err
	}

	filename, err := youtubedl.GetTitleSafe(url)
	if err != nil {
		return &Video{}, err
	}

	v := &Video{
		URL:          url,
		Title:        title,
		Description:  description,
		downloadDate: currentTime,
		filename:     fmt.Sprintf("%s %s", currentTime.Format(time.RFC3339), filename),
	}

	return v, nil
}

// WriteVideo downloads a video and writes it to the specified directory. The
// name of the file is determined by the filename field.
func (v *Video) WriteVideo(directory string) error {
	return youtubedl.DownloadTo(
		v.URL,
		fmt.Sprintf("%s/%s.%s", directory, v.filename, "%(ext)s"),
	)
}

// WriteThumb saves a thumbnail file for the video.
func (v *Video) WriteThumb(directory string) error {
	return youtubedl.DownloadThumbnailTo(
		v.URL,
		fmt.Sprintf("%s/%s-thumb.%s", directory, v.filename, "%(ext)s"),
	)
}

// WriteNfo makes a .nfo file for the video, which includes the video's Title
// and Description.
func (v *Video) WriteNfo(directory string) error {
	nfoXML, err := xml.MarshalIndent(v, "", "    ")
	if err != nil {
		return err
	}

	nfoXML = []byte(fmt.Sprintf("%s\n%s\n%s",
		`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`,
		fmt.Sprintf("<!-- Created on %s by SquirrelBot -->", v.downloadDate.Format(time.RFC3339)),
		nfoXML,
	))
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s.nfo", directory, v.filename), nfoXML, 0644)
	if err != nil {
		return err
	}

	return nil
}
