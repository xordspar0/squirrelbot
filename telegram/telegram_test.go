package telegram

import "testing"

func TestGetMessage(t *testing.T) {
	update := []byte(`{
		"update_id": 1,
		"message": {
			"message_id": 2,
			"text": "abc",
			"from": {
				"id": 3,
				"username": "name"
			}
		}
	}`)

	message, err := GetMessage(update)
	if err != nil {
		t.Error(err.Error())
	}

	errNoMatch := "Field %s does not match: %#v should be %#v"
	if message.ID != 2 {
		t.Errorf(errNoMatch, "ID", message.ID, 2)
	}
	if message.Text != "abc" {
		t.Errorf(errNoMatch, "Text", message.Text, "abc")
	}
	if message.From.ID != 3 {
		t.Errorf(errNoMatch, "From.ID", message.From.ID, 3)
	}
	if message.From.Username != "name" {
		t.Errorf(errNoMatch, "From.Username", message.From.Username, "name")
	}
}
