package bot

import (
	"strings"
	"testing"
)

func TestYoutube(t *testing.T) {
	// TODO: Mock out network and filesystem operations
	err := handleYoutube("http://example.com", "/tmp", 1234, "")
	if err != nil {
		if !strings.HasPrefix(err.Error(), "Failed to send message: 404 Not Found") {
			t.Error(err.Error())
		}
	}
}
