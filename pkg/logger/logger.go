package logger

import (
	"github.com/GFLdev/gorrent/pkg/utils"
	"sync"
	"time"
)

type Color string

const (
	Reset   Color = "\033[0m"
	Red     Color = "\033[31;1m"
	Green   Color = "\033[32;1m"
	Blue    Color = "\033[34;1m"
	Magenta Color = "\033[35;1m"
	Gray    Color = "\033[37;1m"
)

const logLevelLength = 5

type LogLevel string

const (
	Debug LogLevel = "DEBUG"
	Info  LogLevel = "INFO"
	Warn  LogLevel = "WARN"
	Error LogLevel = "ERROR"
	Fatal LogLevel = "FATAL"
)

type Logger struct {
	Config
}

type Config struct {
	Development bool
	SaveLogs    bool
	LogDir      string
	LogFileName string
	PadLen      int
	mux         *sync.Mutex
}

func NewLogger(config Config) *Logger {
	if config.PadLen <= 0 {
		config.PadLen = 3
	}
	config.mux = &sync.Mutex{}

	return &Logger{config}
}

func (l *Logger) getTimestamp() string {
	ts := time.Now().Format("2006-01-02 15:04:05 -0700")
	return string(Gray) + ts + string(Reset)
}

func (l *Logger) getType(t LogLevel) string {
	pad := logLevelLength + l.PadLen
	text := utils.RPad(string(t), pad, " ")
	text = utils.LPad(text, pad+l.PadLen, " ")

	color := Blue
	switch t {
	case Debug:
		color = Green
	case Info:
		color = Blue
	case Warn:
		color = Magenta
	case Error:
	case Fatal:
		color = Red
	}

	return string(color) + text + string(Reset)
}

func (l *Logger) log(msg string) {
	println(msg)
}

func (l *Logger) Debug(msg string) {
	if !l.Development {
		return
	}

	msg = l.getTimestamp() + l.getType(Debug) + msg
	l.log(msg)
}

func (l *Logger) Info(msg string) {
	msg = l.getTimestamp() + l.getType(Info) + msg
	l.log(msg)
}

func (l *Logger) Warn(msg string) {
	msg = l.getTimestamp() + l.getType(Warn) + msg
	l.log(msg)
}

func (l *Logger) Error(msg string) {
	msg = l.getTimestamp() + l.getType(Error) + msg
	l.log(msg)
}

func (l *Logger) Fatal(msg string) {
	msg = l.getTimestamp() + l.getType(Fatal) + msg
	l.log(msg)
	panic(nil)
}
