package youtubedl

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"
)

// Make sure that youtube-dl is installed.
func init() {
	_, err := exec.LookPath("youtube-dl")
	if err != nil {
		panic("youtube-dl is not installed.")
	}
}

// Download uses youtube-dl to download a video and its thumbnail.
func Download(url string) error {
	timestamp := time.Now().Local().Format(time.RFC3339)
	cmd := exec.Command(
		"youtube-dl",
		"--write-thumbnail",
		"--output",
		fmt.Sprintf("%s %s.%s", timestamp, "%(title)s", "%(ext)s"),
		url,
	)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return errors.New("youtubedl: Failed to download video: " + err.Error())
	}

	err = cmd.Start()
	if err != nil {
		return errors.New("youtubedl: Failed to download video: " + err.Error())
	}

	// Catch any error messgages while the process is running.
	errMessages, err := ioutil.ReadAll(stderr)
	if err != nil {
		return errors.New("youtubedl: " + err.Error())
	}
	err = cmd.Wait()
	if err != nil {
		return errors.New("External program youtube-dl: \n" + string(errMessages))
	}

	return nil
}

// GetTitle uses youtube-dl to get the name of a video.
func GetTitle(url string) (videoTitle string) {
	cmd := exec.Command(
		"youtube-dl",
		"--get-title",
		url,
	)

	stdout, err := cmd.StdoutPipe()
	var out []byte
	if err != nil {
		return
	}

	err = cmd.Start()
	if err != nil {
		return
	}

	out, err = ioutil.ReadAll(stdout)
	if err != nil {
		return
	}

	err = cmd.Wait()
	if err != nil {
		return
	}

	videoTitle = strings.TrimSpace(string(out))

	// If there are any errors getting the video name, then just leave it blank.
	return
}
