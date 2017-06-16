binname=squirrelbot
prefix=/usr/local
systemd_unit_path=/etc/systemd/system

.PHONY: build clean fmt install

build:
	go build -o "${binname}" .

fmt:
	gofmt -s -l -w $(shell find . -name '*.go' -not -path '*vendor*')

install:
	install -m 755 "${binname}" "${prefix}/bin/"
	install -m 644 system/squirrelbot.service "${systemd_unit_path}/"

clean:
	rm squirrelbot
