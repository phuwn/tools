package log

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/k0kubun/pp"

	"github.com/phuwn/tools/util"
)

var logger Logger

// Logger - uses of logger in app
type Logger interface {
	Error(err error, msg ...string)
	Info(info, msg string, args ...interface{})
	Status(msg string)
	Fatal(err error)
	Color(color int, msg string)
}

func init() {
	logger = &DefaultLogger{}
}

func Error(err error) {
	logger.Error(err)
}

func Errorf(err error, msg string, args ...interface{}) {
	logger.Error(err, fmt.Sprintf(msg, args...))
}

func Info(info string, args ...interface{}) {
	logger.Info(fmt.Sprintf(info, args...), "")
}

func Infof(info, msg string, args ...interface{}) {
	logger.Info(info, msg, args...)
}

// Status - log server status
func Status(msg string, args ...interface{}) {
	logger.Status(fmt.Sprintf(msg, args...))
}

func Fatal(err error) {
	logger.Fatal(err)
}

func Fatalf(msg string, args ...interface{}) {
	logger.Fatal(fmt.Errorf(msg, args...))
}

func Color(color int, msg string, args ...interface{}) {
	logger.Color(color, fmt.Sprintf(msg, args...))
}

func PP(args ...interface{}) {
	pp.Println(args...)
}

type DefaultLogger struct {
	msg    string
	Locale *time.Location
}

func (l DefaultLogger) Error(err error, msg ...string) {
	l.time()
	cW(ErrorLog, "[ERROR] ")
	fmt.Printf("%s%+v\n", strings.Join(msg, "\n"), err)
}

func (l DefaultLogger) Info(info, msg string, args ...interface{}) {
	l.time()
	s := fmt.Sprintf("%s %s\n", info, baseTrace())
	cW(InfoLog, "[INFO] ")
	if msg != "" {
		s += msg + "\n"
	}
	fmt.Printf(s, args...)
}

func (l DefaultLogger) Status(msg string) {
	l.time()
	cW(StatusLog, "[STATUS] ")
	fmt.Printf("%s\n", msg)
}

func (l DefaultLogger) Color(c int, msg string) {
	cW(&formatter{0, color(c)}, msg)
}

func (l DefaultLogger) Fatal(err error) {
	logger.Error(err)
	os.Exit(1)
}

func (l DefaultLogger) time() {
	now := time.Now()
	if l.Locale != nil {
		now = now.In(l.Locale)
	}
	cW(&formatter{2, 7}, now.Format("2006/01/02 15:04:05MST "))
}

func baseTrace() string {
	frame := util.Caller()
	return fmt.Sprintf("\033[3;36mat %v\033[0m", *frame)
}
