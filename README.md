Squirrelbot
===========

Squirrelbot is a Telegram bot that saves links that you send it for viewing
later. Currently, it only saves Youtube videos and formats them for easy viewing
in Kodi. On the roadmap are handling videos from other websites and saving
arbitrary links in a personal RSS feed.

Building
--------

Squirrelbot is written in [Go](https://golang.org). You'll need the [Go
toolchain](https://golang.org/doc/install) to build the bot.

### Dependencies

Dependencies are managed by [Go dep](https://github.com/golang/dep), so install
dep and then run `dep ensure`.

[youtube-dl](https://rg3.github.io/youtube-dl/) is a runtime dependency. You
should be able to install it through your system package manager.

### Build

To build squirrelbot, simply run `make` in the root directory. To install
squirrelbot, run `make install` as root.

Some build options are available as `make` variables. You can change the system
config file location, systemd unit file location, and others. Look at the
Makefile for details.

Running
-------

### Telegram Token

To run this bot, you first need to get an API token from Telegram. The
directions for doing that are [here](https://core.telegram.org/bots).

Once you have your API token, run Squirrelbot with the required command-line
arguments:

```sh
squirrelbot --address=http://myserver.example.com --port=80 \
  --token=<your Telegram token>
```

### Ports

I use a reverse proxy to forward traffic from port 443 to Squirrelbot's default
port. You can also set up Squirrelbot directly on port 80 or 443. Just make sure
to use the appropriate port for your url scheme:

*	Port 80 for http://
*	Port 443 for https://

### Download Directory

You can optionally specify a directory to download the videos to with the
`--dir` argument:

```sh
squirrelbot --address=http://myserver.example.com --port=80 \
  --token=<your Telegram token> --dir="Youtube Videos"
```

### MOTD

If you want Squirrelbot to respond to "/start" messages (which is a commonly
used phrase for starting conversations with Telegram bots), you should set the
motd (message of the day) flag:

```sh
squirrelbot --address=http://myserver.example.com --port=80 \
  --token=<your Telegram token> --motd="Hello! Try sending me a link."
```

### Configuration with a YAML file

All command line options can also be set in a YAML config file (command line
options override config file options). By default, Squirrelbot looks for the
config file at `/etc/squirrelbot/config.yaml`

Here is an example config file:

```yaml
address: https://myserver.example.com
port: 80
token: <your Telegram token>
directory: Youtube Videos
motd: Hello! Try sending me a link.
```

Tips
----

### Transfering Video Files

Squirrelbot downloads videos to a local directory that can be specified with the
`--dir` option. The videos are formatted to be easy to view in Kodi. If you run
this bot on a different server than your Kodi/media server, you will want to
transfer your video files to your Kodi or media server box. Here are a couple
ways to do that.

#### Sync with rsync

*	On the receiving end, set up an rsync daemon that allows write-only access
	to the right directory.
*	On the sending end, write a simple script that uses `inotifywait` to copy
	files via rsync.
*	Optionally, after the file is successfully send, the script should delete
	the file from the server.

#### Sync with Syncthing

You can use Syncthing to send video files to their final destination.
See https://docs.syncthing.net/intro/getting-started.html for instructions.

License
-------

Copyright © 2017-2018 Jordan Christiansen.

Squirrelbot is distributed under the terms of the GNU GPL version 3 or later.
This is free software: you are free to change and redistribute it. There is
NO WARRANTY, to the extent permitted by law.

The full text of the GPL version 3 is in the file LICENSE. For more
information, visit http://gnu.org/licenses/gpl.html.
