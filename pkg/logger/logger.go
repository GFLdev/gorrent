package logger

import (
	"github.com/GFLdev/gorrent/pkg/utils"
	"sync"
	"time"
)

type Color string

const (
	Reset   Color = "\033[0m"
	Red     Color = "\033[31m"
	Green   Color = "\033[32m"
	Blue    Color = "\033[34m"
	Magenta Color = "\033[35m"
)

const logLevelLength = 5

type LogLevel string

const (
	Debug LogLevel = "DEBUG"
	Info  LogLevel = "INFO"
	Warn  LogLevel = "WARN"
	Error LogLevel = "ERROR"
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

func NewLogger(config *Config) *Logger {
	if config == nil {
		config = &Config{
			Development: false,
			SaveLogs:    false,
			LogDir:      "",
			LogFileName: "",
			PadLen:      3,
		}
	}
	config.mux = &sync.Mutex{}

	return &Logger{*config}
}

func (l *Logger) getTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05 -0700")
}

func (l *Logger) getType(t LogLevel) string {
	text := utils.RPad(string(t), logLevelLength+l.PadLen, " ")
	text = utils.LPad(text, l.PadLen, " ")

	color := Blue
	switch t {
	case Debug:
		color = Green
	case Info:
		color = Blue
	case Warn:
		color = Magenta
	case Error:
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
	msg = l.getTimestamp() + l.getType(Error) + msg
	l.log(msg)
	panic(nil)
}
