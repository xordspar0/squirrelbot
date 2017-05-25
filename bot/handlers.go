package bot

func handleYoutube(url string) {
	telegram.SendMessage("I saved your Youtube video.")
}

func handleLink(url string) {
	telelgram.SendMessage("Stashing ordinary links is not yet implemented.")
}
