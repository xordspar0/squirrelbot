package main

import (
	"github.com/xordspar0/squirrelbot/cmd"

	"github.com/spf13/cobra/doc"

	"bytes"
	"fmt"
	"strings"
)

var tips = `
.SH TIPS
.SS Configuration Files
.PP
Config files are YAML files.
Specify options with the following syntax:
.IP
.nf
\f[C]
address:\ example.com
token:\ 123456:ABC\-DEF1234ghIkl\-zyx57W2v1u123ew11
\f[]
.fi
.SS Transferring Video Files
.PP
Squirrelbot downloads videos to a local directory that can be specified
with the \[en]dir option.
The videos are formatted to be easy to view in Kodi.
If you run this bot on a different server than your Kodi/media center,
you will want to transfer your video files to your Kodi or media center
box.
Here are a couple ways to do that.
.PP
Sync with rsync:
.IP \[bu] 2
On the receiving end, set up an rsync daemon that allows write\-only
access to the right directory.
.IP \[bu] 2
On the sending end, write a simple script that uses inotifywait to copy
files via rsync.
.IP \[bu] 2
Optionally, after the file is successfully send, the script should
delete the file from the server.
.PP
Sync with Syncthing:
.PP
You can use Syncthing to send video files to their final destination.
See https://docs.syncthing.net/intro/getting\-started.html for
instructions.
`

var footer = `
.SH "COPYRIGHT"
.sp
Copyright \(co 2017 Jordan Christiansen\&. License GPLv3+: GNU GPL version 3 or later http://gnu\&.org/licenses/gpl\&.html\&. This is free software: you are free to change and redistribute it\&. There is NO WARRANTY, to the extent permitted by law\&.
.SH "SEE ALSO"
.sp
youtube\-dl(1)
`

func main() {
	header := &doc.GenManHeader{
		Section: "1",
		Source:  fmt.Sprintf("%s %s", cmd.Cmd.Use, cmd.Cmd.Version),
		Manual:  fmt.Sprintf("%s Manual", strings.Title(cmd.Cmd.Use)),
	}

	buf := new(bytes.Buffer)
	err := doc.GenMan(cmd.Cmd, header, buf)
	if err != nil {
		panic(err)
	}

	buf.WriteString(tips)
	buf.WriteString(footer)

	fmt.Print(buf.String())
}
