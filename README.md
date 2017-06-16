SquirrelBot
===========

SquirrelBot is a Telegram bot that saves links that you send it for viewing
later. Currently, it only saves Youtube videos and formats them for easy viewing
in Kodi. On the roadmap are handling videos from other websites and saving
arbitrary links in a personal RSS feed.


Transfering Video Files
-----------------------

SquirrelBot downloads videos to a local directory that can be specified with the
`--dir` option. The videos are formatted to be easy to view in Kodi. If you run
this bot on a different server than your Kodi/media server, you will want to
transfer your video files to your Kodi or media server box. Here are a couple
ways to do that.

### Sync with rsync

*	On the receiving end, set up an rsync daemon that allows write-only access
	to the right directory.
*	On the sending end, write a simple script that uses `inotifywait` to copy
	files via rsync.
*	Optionally, after the file is successfully send, the script should delete
	the file from the server.

### Sync with Syncthing

You can use Syncthing to send video files to their final destination.
See https://docs.syncthing.net/intro/getting-started.html for instructions.
