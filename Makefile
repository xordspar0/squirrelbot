binname=squirrelbot
version=devel
prefix=/usr/local
systemd_unit_path=/etc/systemd/system
system_config_file=/etc/squirrelbot/config.yaml

.PHONY: build clean docker docs fmt install snap test uninstall

# Building Commands

build:
	env CGO_ENABLED=0 go build -ldflags "-X github.com/xordspar0/squirrelbot/cmd.version=$(version) -X github.com/xordspar0/squirrelbot/cmd.systemConfigFile=$(system_config_file)" -o "build/$(binname)"

docs: build
	go run -ldflags "-X github.com/xordspar0/squirrelbot/cmd.version=$(version) -X github.com/xordspar0/squirrelbot/cmd.systemConfigFile=$(system_config_file)" tools/gendocs.go > "build/$(binname).1"

buildall: build docs

docker:
	docker build . --tag $(binname)

snap: buildall
	packages/build_snap.sh $(version)

# Installing Commands

install: buildall
	install -d -m 755 "$(prefix)/bin/"
	install -m 755 "build/$(binname)" "$(prefix)/bin/$(binname)"
	install -d -m 755 "$(systemd_unit_path)"
	install -m 644 system/squirrelbot.service "$(systemd_unit_path)/squirrelbot.service"
	install -d -m 755 "$(prefix)/share/man/man1/"
	install -m 644 build/squirrelbot.1 "$(prefix)/share/man/man1/squirrelbot.1"

uninstall:
	-rm -f "$(prefix)/bin/$(binname)"
	-rm -f "$(systemd_unit_path)/squirrelbot.service"
	-rm -f "$(prefix)/share/man/man1/squirrelbot.1"

# Maintenance Commands

fmt:
	go fmt ./...

test:
	go test ./...

clean:
	-rm -rf build/
	-rm -f squirrelbot.1
