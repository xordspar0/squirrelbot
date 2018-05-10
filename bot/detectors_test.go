package bot

import (
	"testing"
)

func TestIsYoutubeSourceValid(t *testing.T) {
	url := "https://www.youtube.com/v=12345"
	if !isYoutubeSource(url) {
		t.Fail()
	}
}

func TestIsYoutubeSourceInvalid(t *testing.T) {
	url := "https://example.com"
	if isYoutubeSource(url) {
		t.Fail()
	}
}
