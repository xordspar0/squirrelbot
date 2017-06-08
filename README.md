SquirrelBot
===========

SquirrelBot is a Telegram bot that saves links that you send it for viewing
later. Currently, it only saves Youtube videos and formats them for easy viewing
in Kodi. On the roadmap are handling videos from other websites and saving
arbitrary links in a  personal RSS feed.

Video File Sync via rsync
-------------------------

*	On the receiving end, set up an rsync daemon that allows write-only access
	to the right directory.
*	On the sending end, write a simple script that uses `inotifywait` to copy
	files via rsync.
