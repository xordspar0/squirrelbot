package bot

import "strings"

func isYoutubeSource(url string) bool {
	validPrefixes := []string{
		"http://www.youtube.com",
		"https://www.youtube.com",
		"http://m.youtube.com",
		"https://m.youtube.com",
		"http://youtu.be",
		"https://youtu.be",
		"http://vimeo.com",
		"https://vimeo.com",
		"http://player.vimeo.com",
		"https://player.vimeo.com",
	}

	for _, prefix := range validPrefixes {
		if strings.HasPrefix(url, prefix) {
			return true
		}
	}

	return false
}
