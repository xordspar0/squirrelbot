binname=squirrelbot
version=devel
prefix=/usr/local
systemd_unit_path=/etc/systemd/system
system_config_file=/etc/squirrelbot/config.yaml

.PHONY: build docker clean fmt install snap test uninstall

# Building Commands

build:
	go build -ldflags "-X main.version=$(version) -X main.systemConfigFile=$(system_config_file)" -o "bin/$(binname)" ./cmd/squirrelbot

squirrelbot.1: doc/squirrelbot.txt
	a2x -f manpage doc/squirrelbot.txt

buildall: build squirrelbot.1

docker:
	env CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=$(version) -X main.systemConfigFile=$(system_config_file)" -o "bin/$(binname).docker" ./cmd/squirrelbot
	docker build . -f packages/Dockerfile -t squirrelbot:$(version)

snap: buildall
	packages/build_snap.sh $(version)

# Installing Commands

install: buildall
	install -Dm 755 "bin/$(binname)" "$(prefix)/bin/$(binname)"
	install -Dm 644 system/squirrelbot.service "$(systemd_unit_path)/squirrelbot.service"
	install -Dm 644 doc/squirrelbot.1 "$(prefix)/share/man/man1/squirrelbot.1"

uninstall:
	-rm -f "$(prefix)/bin/$(binname)"
	-rm -f "$(systemd_unit_path)/squirrelbot.service"
	-rm -f "$(prefix)/share/man/man1/squirrelbot.1"

# Maintenance Commands

fmt:
	gofmt -s -l -w $(shell find . -name '*.go' -not -path '*vendor*')

test:
	go test ./...

clean:
	-rm -rf bin/
	-rm -f squirrelbot.1
