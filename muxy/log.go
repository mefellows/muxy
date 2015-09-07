package muxy

import (
	"fmt"
	"github.com/mgutz/ansi"
	"log"
)

type LogLevel int
type Colour string

const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
)

var (
	/*
		RESET        Colour = ansi.Reset
		BLACK        Colour = Colour(ansi.Black)
		BLUE         Colour = Colour(ansi.Blue)
		RED          Colour = Colour(ansi.Red)
		GREEN        Colour = Colour(ansi.Green)
		YELLOW       Colour = Colour(ansi.Yellow)
		MAGENTA      Colour = Colour(ansi.Magenta)
		CYAN         Colour = Colour(ansi.Cyan)
		WHITE        Colour = Colour(ansi.White)
		LIGHTBLACK   Colour = Colour(ansi.LightBlack)
		LIGHTRED     Colour = Colour(ansi.LightRed)
		LIGHTGREEN   Colour = Colour(ansi.LightGreen)
		LIGHTYELLOW  Colour = Colour(ansi.LightYellow)
		LIGHTBLUE    Colour = Colour(ansi.LightBlue)
		LIGHTMAGENTA Colour = Colour(ansi.LightMagenta)
		LIGHTCYAN    Colour = Colour(ansi.LightCyan)
		LIGHTWHITE   Colour = Colour(ansi.LightWhite)
	*/
	RESET        Colour = ansi.Reset
	BLACK        Colour = Colour(ansi.ColorCode("black"))
	RED          Colour = Colour(ansi.ColorCode("red"))
	GREEN        Colour = Colour(ansi.ColorCode("green"))
	YELLOW       Colour = Colour(ansi.ColorCode("yellow"))
	BLUE         Colour = Colour(ansi.ColorCode("blue"))
	MAGENTA      Colour = Colour(ansi.ColorCode("magenta"))
	CYAN         Colour = Colour(ansi.ColorCode("cyan"))
	WHITE        Colour = Colour(ansi.ColorCode("white"))
	LIGHTBLACK   Colour = Colour(ansi.ColorCode("black+h"))
	LIGHTRED     Colour = Colour(ansi.ColorCode("red+h"))
	LIGHTGREEN   Colour = Colour(ansi.ColorCode("green+h"))
	LIGHTYELLOW  Colour = Colour(ansi.ColorCode("yellow+h"))
	LIGHTBLUE    Colour = Colour(ansi.ColorCode("blue+h"))
	LIGHTMAGENTA Colour = Colour(ansi.ColorCode("magenta+h"))
	LIGHTCYAN    Colour = Colour(ansi.ColorCode("cyan+h"))
	LIGHTWHITE   Colour = Colour(ansi.ColorCode("white+h"))
)

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

func (m *MuxyLogger) Log(l LogLevel, format string, v ...interface{}) {
	if l >= m.Level {
		var level string
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
		}
		log.Printf("["+level+"]\t\t"+format+"\n", v...)
	}
}

func (m *MuxyLogger) Colorize(colour Colour, format string) string {
	return fmt.Sprintf("%s%s%s", string(colour), format, RESET)
}
