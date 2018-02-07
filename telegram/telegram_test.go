package telegram

import "testing"

func TestGetMessage(t *testing.T) {
	// Setup
	update := []byte(`{
		"update_id": 1,
		"message": {
			"message_id": 2,
			"text": "abc",
			"chat": {
				"id": 3,
				"username": "chat"
			},
			"from": {
				"id": 4,
				"username": "name"
			}
		}
	}`)

	message, err := GetMessage(update)
	if err != nil {
		t.Error(err.Error())
	}

	// Tests
	errNoMatch := "Field %s does not match: %#v should be %#v"
	if message.ID != 2 {
		t.Errorf(errNoMatch, "ID", message.ID, 2)
	}
	if message.Text != "abc" {
		t.Errorf(errNoMatch, "Text", message.Text, "abc")
	}
	if message.Chat.ID != 3 {
		t.Errorf(errNoMatch, "Chat.ID", message.Chat.ID, 3)
	}
	if message.Chat.Username != "chat" {
		t.Errorf(errNoMatch, "Chat.Username", message.Chat.Username, "chat")
	}
	if message.From.ID != 4 {
		t.Errorf(errNoMatch, "From.ID", message.From.ID, 4)
	}
	if message.From.Username != "name" {
		t.Errorf(errNoMatch, "From.Username", message.From.Username, "name")
	}
}
