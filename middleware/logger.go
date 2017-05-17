// Package middleware contains default middleware implementations.
package middleware

import (
	"fmt"

	"github.com/mefellows/muxy/log"
	"github.com/mefellows/muxy/muxy"
	"github.com/mefellows/plugo/plugo"
)

// LoggerMiddleware is middleware that will log all request/responses at
// a configured level
type LoggerMiddleware struct {
	HexOutput bool `mapstructure:"hex_output"`
	format    string
}

const bytesTab = "\n\t\t\t\t\t\t"

func init() {
	plugo.PluginFactories.Register(func() (interface{}, error) {
		return &LoggerMiddleware{}, nil
	}, "logger")
}

// Setup sets up the middleware
func (l *LoggerMiddleware) Setup() {
	if l.HexOutput {
		l.format = "%x"
	} else {
		l.format = "%s"
	}
}

// Teardown shuts down the middleware
func (LoggerMiddleware) Teardown() {
}

// HandleEvent takes a ProxyEvent and acts on the information provided
func (l *LoggerMiddleware) HandleEvent(e muxy.ProxyEvent, ctx *muxy.Context) {
	switch e {
	case muxy.EventPreDispatch:
		if ctx.Request == nil {
			if len(ctx.Bytes) > 0 {
				log.Info("Handle TCP event " + log.Colorize(log.GREY, "PRE_DISPATCH") + fmt.Sprintf(" Received %d%s", len(ctx.Bytes), " bytes"))
				data := fmt.Sprintf(l.format, ctx.Bytes)
				log.Debug("Handle TCP event " + log.Colorize(log.GREY, "PRE_DISPATCH") + " Received request: " + bytesTab + log.Colorize(log.BLUE, data))
			}
		} else {
			log.Info("Handle HTTP event " + log.Colorize(log.GREY, "PRE_DISPATCH") + " Proxying request " +
				log.Colorize(log.LIGHTMAGENTA, ctx.Request.Method) +
				log.Colorize(log.BLUE, " \""+ctx.Request.URL.String()+"\""))
		}
	case muxy.EventPostDispatch:
		if ctx.Request == nil {
			if len(ctx.Bytes) > 0 {
				log.Info("Handle TCP event " + log.Colorize(log.GREY, "POST_DISPATCH") + fmt.Sprintf(" Sent %d%s", len(ctx.Bytes), " bytes"))
				data := fmt.Sprintf(l.format, ctx.Bytes)
				log.Debug("Handle TCP event " + log.Colorize(log.GREY, "POST_DISPATCH") + " Sent response: " + bytesTab + log.Colorize(log.BLUE, data))
			}
		} else {
			log.Info("Handle HTTP event " + log.Colorize(log.GREY, "POST_DISPATCH") + " Returning " +
				log.Colorize(log.GREY, fmt.Sprintf("%s", ctx.Response.Status)))
		}
	}
}
