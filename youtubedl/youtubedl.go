package youtubedl

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

// Make sure that youtube-dl is installed.
func init() {
	_, err := exec.LookPath("youtube-dl")
	if err != nil {
		panic("youtube-dl is not installed.")
	}
}

// Download uses youtube-dl to download a video.
func Download(url, directory, prefix string) error {
	cmd := exec.Command(
		"youtube-dl",
		"--output",
		fmt.Sprintf("%s/%s%s.%s", directory, prefix, "%(title)s", "%(ext)s"),
		url,
	)

	return downloadFile(cmd)
}

// Download uses youtube-dl to download the thumbnail of a video.
func DownloadThumbnail(url, directory, prefix string) error {
	// Skip the download of the video itself, but download the thumbnail. Put a
	// "-thumb" suffix on the file name because that is the format that Kodi
	// recognizes.
	cmd := exec.Command(
		"youtube-dl",
		"--skip-download",
		"--write-thumbnail",
		"--output",
		fmt.Sprintf("%s/%s%s-thumb.%s", directory, prefix, "%(title)s", "%(ext)s"),
		url,
	)

	return downloadFile(cmd)
}

// GetTitle uses youtube-dl to get the name of a video.
func GetTitle(url string) (videoTitle string) {
	cmd := exec.Command(
		"youtube-dl",
		"--get-title",
		url,
	)

	return getCmdOutput(cmd)
}

// GetDescription uses youtube-dl to get the description of a video.
func GetDescription(url string) string {
	cmd := exec.Command(
		"youtube-dl",
		"--get-description",
		url,
	)

	return getCmdOutput(cmd)
}

func downloadFile(cmd *exec.Cmd) error {
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

// getCmdResponse runs a command and returns its output. If there are any errors
// in the command, then just return an empty string.
func getCmdOutput(cmd *exec.Cmd) string {
	stdout, err := cmd.StdoutPipe()
	var out []byte
	if err != nil {
		return ""
	}

	err = cmd.Start()
	if err != nil {
		return ""
	}

	out, err = ioutil.ReadAll(stdout)
	if err != nil {
		return ""
	}

	err = cmd.Wait()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}
