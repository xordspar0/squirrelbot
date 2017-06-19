binname=squirrelbot
prefix=/usr/local
systemd_unit_path=/etc/systemd/system

.PHONY: build clean fmt install uninstall

build:
	go build -o "$(binname)" ./cmd/squirrelbot

fmt:
	gofmt -s -l -w $(shell find . -name '*.go' -not -path '*vendor*')

install:
	install -m 755 "$(binname)" "$(prefix)/bin/"
	install -m 644 system/squirrelbot.service "$(systemd_unit_path)/"

uninstall:
	rm -f "$(prefix)/bin/$(binname)"
	rm -f "$(systemd_unit_path)/squirrelbot.service"

clean:
	rm squirrelbot
