package log

import (
	"fmt"
	"github.com/mgutz/ansi"
	"log"
	"os"
)

type LogLevel int
type Colour int

const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

const (
	BLACK Colour = iota
	BLUE
	RED
	GREEN
	GREY
	YELLOW
	MAGENTA
	CYAN
	WHITE
	LIGHTBLACK
	LIGHTRED
	LIGHTGREEN
	LIGHTYELLOW
	LIGHTBLUE
	LIGHTMAGENTA
	LIGHTCYAN
	LIGHTWHITE
)

var coloursMap = map[Colour]string{
	BLACK:        ansi.ColorCode("black"),
	RED:          ansi.ColorCode("red"),
	GREEN:        ansi.ColorCode("green"),
	GREY:         string([]byte{'\033', '[', '3', '2', ';', '1', 'm'}),
	YELLOW:       ansi.ColorCode("yellow"),
	BLUE:         ansi.ColorCode("blue"),
	MAGENTA:      ansi.ColorCode("magenta"),
	CYAN:         ansi.ColorCode("cyan"),
	WHITE:        ansi.ColorCode("white"),
	LIGHTBLACK:   ansi.ColorCode("black+h"),
	LIGHTRED:     ansi.ColorCode("red+h"),
	LIGHTGREEN:   ansi.ColorCode("green+h"),
	LIGHTYELLOW:  ansi.ColorCode("yellow+h"),
	LIGHTBLUE:    ansi.ColorCode("blue+h"),
	LIGHTMAGENTA: ansi.ColorCode("magenta+h"),
	LIGHTCYAN:    ansi.ColorCode("cyan+h"),
	LIGHTWHITE:   ansi.ColorCode("white+h"),
}

// Logging facility
type MuxyLogger struct {
	log.Logger
	Level LogLevel
}

func init() {
	if fl := log.Flags(); fl&log.Ltime != 0 {
		log.SetFlags(fl | log.Lmicroseconds)
	}
}

func NewLogger() *MuxyLogger {
	return &MuxyLogger{Level: INFO}
}

var std = NewLogger()

func (m *MuxyLogger) Trace(format string, v ...interface{}) {
	m.Log(TRACE, format, v...)
}

func (m *MuxyLogger) Debug(format string, v ...interface{}) {
	m.Log(DEBUG, format, v...)
}

func (m *MuxyLogger) Info(format string, v ...interface{}) {
	m.Log(INFO, format, v...)
}

func (m *MuxyLogger) Warn(format string, v ...interface{}) {
	m.Log(WARN, format, v...)
}

func (m *MuxyLogger) Error(format string, v ...interface{}) {
	m.Log(ERROR, format, v...)
}

func (m *MuxyLogger) Fatal(v ...interface{}) {
	s := fmt.Sprint(v...)
	m.Log(FATAL, s)
	os.Exit(1)
}

func (m *MuxyLogger) Fatalf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	m.Log(FATAL, s)
	os.Exit(1)
}

func (m *MuxyLogger) Log(l LogLevel, format string, v ...interface{}) {
	if l >= m.Level {
		var level string
		var colorFormat = ""
		switch l {
		case TRACE:
			level = "TRACE"
		case DEBUG:
			level = "DEBUG"
		case INFO:
			level = "INFO"
		case WARN:
			level = "WARN"
		case ERROR:
			level = "ERROR"
			colorFormat = coloursMap[LIGHTRED]
		case FATAL:
			level = "FATAL"
			colorFormat = coloursMap[LIGHTRED]
		}
		log.Printf("["+level+"]\t\t"+colorFormat+format+ansi.Reset+"\n", v...)
	}
}

func (m *MuxyLogger) SetLevel(l LogLevel) {
	m.Level = l
}

func Colorize(colour Colour, format string) string {
	return fmt.Sprintf("%s%s%s", coloursMap[colour], format, ansi.Reset)
}

func Trace(format string, v ...interface{}) {
	std.Log(TRACE, format, v...)
}

func Debug(format string, v ...interface{}) {
	std.Log(DEBUG, format, v...)
}

func Info(format string, v ...interface{}) {
	std.Log(INFO, format, v...)
}

func Warn(format string, v ...interface{}) {
	std.Log(WARN, format, v...)
}

func Error(format string, v ...interface{}) {
	std.Log(ERROR, format, v...)
}

func Fatal(v ...interface{}) {
	std.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	std.Fatalf(format, v...)
}

func Log(l LogLevel, format string, v ...interface{}) {
	std.Log(l, format, v...)
}

func SetLevel(l LogLevel) {
	std.SetLevel(l)
}
