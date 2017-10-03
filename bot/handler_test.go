package bot

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNfo(t *testing.T) {
	messageJson := []byte(`{"from": {"id": 1234}}`)
	var message map[string]interface{}
	err := json.Unmarshal(messageJson, &message)
	if err != nil {
		t.Error(err.Error())
	}

	err = handleYoutube(message, "http://example.com", "/tmp", "")
	if err != nil {
		if !strings.HasPrefix(err.Error(), "Failed to send message: 404 Not Found") {
			t.Error(err.Error())
		}
	}
}
