package log

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

var isTTY bool

func init() {
	isTTY = terminal.IsTerminal(int(os.Stdout.Fd()))
}

type color int8

const (
	black color = iota
	red
	green
	yellow
	blue
	magenta
	cyan
	white
)

type style int8

const (
	normal style = iota
	bold
	faded
	italic
	underscore
)

type formatter struct {
	style style
	color color
}

var (
	// Status Log
	InfoLog    = &formatter{1, 2}
	WarningLog = &formatter{1, 3}
	ErrorLog   = &formatter{1, 1}
	StatusLog  = &formatter{1, 6}
	DebugLog   = &formatter{1, 4}

	// Request Log
	ReqUrlLog    = &formatter{3, 6}
	RespBytesLog = &formatter{3, 4}

	reset = []byte{'\033', '[', '0', 'm'}
)

func cW(f *formatter, s string, args ...interface{}) {
	if !isTTY {
		fmt.Print(s)
		return
	}
	fmt.Printf(fmt.Sprintf("\033[%d;3%dm", f.style, f.color)+s+"\033[0m", args...)
}

func methodForm(method string) *formatter {
	var c color
	switch method {
	case "POST":
		c = green
	case "PUT":
		c = yellow
	case "DELETE":
		c = red
	default:
		c = magenta
	}
	return &formatter{1, c}
}

func statusForm(status int) *formatter {
	var c color
	switch {
	case status < 200:
		c = blue
	case status < 300:
		c = green
	case status < 400:
		c = cyan
	case status < 500:
		c = yellow
	default:
		c = red
	}
	return &formatter{1, c}
}
