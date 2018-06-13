package bot

import (
	"testing"
)

func TestIsYoutubeSourceValid(t *testing.T) {
	url := "https://www.youtube.com/watch?v=CKIye4RZ-5k"
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
