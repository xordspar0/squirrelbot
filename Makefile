binname=squirrelbot
version=devel
prefix=/usr/local
systemd_unit_path=/etc/systemd/system
system_config_file=/etc/squirrelbot/config.yaml

.PHONY: build clean fmt snap install uninstall

build:
	go build -ldflags "-X main.version=$(version) -X main.systemConfigFile=$(system_config_file)" -o "$(binname)" ./cmd/squirrelbot

squirrelbot.1: doc/squirrelbot.txt
	a2x -f manpage doc/squirrelbot.txt

fmt:
	gofmt -s -l -w $(shell find . -name '*.go' -not -path '*vendor*')

snap: build squirrelbot.1
	packages/build_snap.sh $(version)

install: squirrelbot.1
	install -Dm 755 "$(binname)" "$(prefix)/bin/$(binname)"
	install -Dm 644 system/squirrelbot.service "$(systemd_unit_path)/squirrelbot.service"
	install -Dm 644 doc/squirrelbot.1 "$(prefix)/share/man/man1/squirrelbot.1"

uninstall:
	-rm -f "$(prefix)/bin/$(binname)"
	-rm -f "$(systemd_unit_path)/squirrelbot.service"
	-rm -f "$(prefix)/share/man/man1/squirrelbot.1"

clean:
	-rm -f squirrelbot
	-rm -f squirrelbot.1
