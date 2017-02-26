package log

import (
	"fmt"
	"log"
	"os"

	"github.com/mgutz/ansi"
)

// Level to set
// nolint
type Level int

// Colour type to print in log messages
type Colour int

// nolint
const (
	TRACE Level = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

// nolint
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

// MuxyLogger is the logging facility for Muxy
type MuxyLogger struct {
	log.Logger
	Level Level
}

func init() {
	if fl := log.Flags(); fl&log.Ltime != 0 {
		log.SetFlags(fl | log.Lmicroseconds)
	}
}

// NewLogger creates a new logger
func NewLogger() *MuxyLogger {
	return &MuxyLogger{Level: INFO}
}

var std = NewLogger()

// Trace Logging
func (m *MuxyLogger) Trace(format string, v ...interface{}) {
	m.Log(TRACE, format, v...)
}

// Debug logging
func (m *MuxyLogger) Debug(format string, v ...interface{}) {
	m.Log(DEBUG, format, v...)
}

// Info loggig
func (m *MuxyLogger) Info(format string, v ...interface{}) {
	m.Log(INFO, format, v...)
}

// Warn logging
func (m *MuxyLogger) Warn(format string, v ...interface{}) {
	m.Log(WARN, format, v...)
}

// Error logging
func (m *MuxyLogger) Error(format string, v ...interface{}) {
	m.Log(ERROR, format, v...)
}

// Fatal logging
func (m *MuxyLogger) Fatal(v ...interface{}) {
	s := fmt.Sprint(v...)
	m.Log(FATAL, s)
	os.Exit(1)
}

// Fatalf formatted fatal logging
func (m *MuxyLogger) Fatalf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	m.Log(FATAL, s)
	os.Exit(1)
}

// Log is general log facility
func (m *MuxyLogger) Log(l Level, format string, v ...interface{}) {
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

// SetLevel sets the log level
func (m *MuxyLogger) SetLevel(l Level) {
	m.Level = l
}

// Colorize returns a coloured log string
func Colorize(colour Colour, format string) string {
	return fmt.Sprintf("%s%s%s", coloursMap[colour], format, ansi.Reset)
}

// Trace logging
func Trace(format string, v ...interface{}) {
	std.Log(TRACE, format, v...)
}

// Debug logging
func Debug(format string, v ...interface{}) {
	std.Log(DEBUG, format, v...)
}

// Info logging
func Info(format string, v ...interface{}) {
	std.Log(INFO, format, v...)
}

// Warn logging
func Warn(format string, v ...interface{}) {
	std.Log(WARN, format, v...)
}

// Error logging
func Error(format string, v ...interface{}) {
	std.Log(ERROR, format, v...)
}

// Fatal logging
func Fatal(v ...interface{}) {
	std.Fatal(v...)
}

// Fatalf is formatted fatal logging
func Fatalf(format string, v ...interface{}) {
	std.Fatalf(format, v...)
}

// Log general log method
func Log(l Level, format string, v ...interface{}) {
	std.Log(l, format, v...)
}

// SetLevel sets log level
func SetLevel(l Level) {
	std.SetLevel(l)
}
