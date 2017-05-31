package bot

import (
	"neolog.xyz/squirrelbot/telegram"
)

func handleYoutube(url string) {
	telegram.SendMessage("I saved your Youtube video.")
}

func handleLink(url string) {
	telegram.SendMessage("This link does not look like a video. Stashing " +
		"ordinary links is not yet implemented.")
}
