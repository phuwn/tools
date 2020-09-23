package log

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/k0kubun/pp"

	"github.com/phuwn/tools/util"
)

var logger Logger

// Logger - uses of logger in app
type Logger interface {
	Error(err error, args ...interface{})
	Info(args ...interface{})
	Status(msg string)
	Fatal(err error, args ...interface{})
	Color(color int, msg string)
}

func init() {
	logger = &DefaultLogger{os.Stdout, nil}
}

// Error - print out error
func Error(err error, args ...interface{}) {
	logger.Error(err, args...)
}

// Errorf - print out error using a formatted message
func Errorf(err error, msg string, args ...interface{}) {
	logger.Error(err, fmt.Sprintf(msg, args...))
}

// Info - print out server's info for debugging
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Infof - print out server's info for debugging using a formatted message
func Infof(msg string, args ...interface{}) {
	logger.Info(fmt.Sprintf(msg, args...))
}

// Status - print out server's status
func Status(msg string, args ...interface{}) {
	logger.Status(fmt.Sprintf(msg, args...))
}

// Fatal - print out server's error followed by a call to os.Exit(1)
func Fatal(err error, args ...interface{}) {
	logger.Fatal(err, args...)
}

// Fatalf - print out server's error followed by a call to os.Exit(1) using a formatted message
func Fatalf(msg string, args ...interface{}) {
	logger.Fatal(fmt.Errorf(msg, args...))
}

// Color - print colorful message
func Color(color int, msg string, args ...interface{}) {
	logger.Color(color, fmt.Sprintf(msg, args...))
}

// PP - pretty printer implement
func PP(args ...interface{}) {
	pp.Println(args...)
}

// DefaultLogger - simple default logger if user dont define their own
type DefaultLogger struct {
	writer io.Writer
	Locale *time.Location
}

func (l DefaultLogger) Error(err error, args ...interface{}) {
	l.time()
	l.cW(ErrorLog, "[ERROR] ")
	if len(args) > 0 {
		fmt.Fprintln(l.writer, args...)
	}
	fmt.Fprintf(l.writer, "%+v\n", err)
}

func (l DefaultLogger) Info(args ...interface{}) {
	l.time()

	l.cW(InfoLog, "[INFO] ")
	fmt.Fprint(l.writer, args...)
	l.cW(&formatter{3, 6}, fmt.Sprintf(" %s\n", baseTrace()))
}

func (l DefaultLogger) Status(msg string) {
	l.time()
	l.cW(StatusLog, "[STATUS] ")
	fmt.Fprintf(l.writer, "%s\n", msg)
}

func (l DefaultLogger) Color(c int, msg string) {
	l.cW(&formatter{0, color(c)}, msg)
}

func (l DefaultLogger) Fatal(err error, args ...interface{}) {
	l.Error(err, args...)
	os.Exit(1)
}

func (l DefaultLogger) time() {
	now := time.Now()
	if l.Locale != nil {
		now = now.In(l.Locale)
	}
	l.cW(&formatter{2, 7}, now.Format("2006/01/02 15:04:05MST "))
}

func (l DefaultLogger) cW(f *formatter, s string, args ...interface{}) {
	if !isTTY {
		fmt.Fprint(l.writer, s)
		return
	}
	fmt.Fprintf(l.writer, fmt.Sprintf("\033[%d;3%dm", f.style, f.color)+s+"\033[0m", args...)
}

func baseTrace() string {
	frame := util.Caller(4)
	return fmt.Sprintf("%v", *frame)
}
