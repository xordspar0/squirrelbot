package telegram

import "testing"

func TestGetMessage(t *testing.T) {
	update := []byte(`{
		"update_id": 1,
		"message": {
			"message_id": 2,
			"text": "abc",
			"from": {
				"id": 3
			},
			"chat": {
				"id": 4
			}
		}
	}`)

	message, err := GetMessage(update)
	if err != nil {
		t.Error(err.Error())
	}

	noMatch := "Field %s does not match: %#v should be %#v"
	if message.ID != 2 {
		t.Errorf(noMatch, "ID", message.ID, 2)
	}
	if message.Text != "abc" {
		t.Errorf(noMatch, "Text", message.Text, "abc")
	}
	if message.From.ID != 3 {
		t.Errorf(noMatch, "From.ID", message.From.ID, 3)
	}
	if message.Chat.ID != 4 {
		t.Errorf(noMatch, "Chat.ID", message.Chat.ID, 4)
	}
}
