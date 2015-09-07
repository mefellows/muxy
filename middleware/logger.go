package middleware

import (
	"github.com/mefellows/muxy/log"
	"github.com/mefellows/muxy/muxy"
	"github.com/mefellows/plugo/plugo"
	"time"
)

type LoggerMiddleware struct {
}

const DEFAULT_DELAY = 2 * time.Second

func init() {
	plugo.PluginFactories.Register(func() (interface{}, error) {
		return &LoggerMiddleware{}, nil
	}, "logger")
}

func (LoggerMiddleware) Setup() {
}

func (LoggerMiddleware) Teardown() {
}

func (LoggerMiddleware) HandleEvent(e muxy.ProxyEvent, ctx *muxy.Context) {
	switch e {
	case muxy.EVENT_PRE_DISPATCH:
		if ctx.Request == nil {
			if len(ctx.Bytes) > 0 {
				log.Info("Handle TCP event " + log.Colorize(log.GREY, "PRE_DISPATCH") + " Proxying request: " + string(ctx.Bytes))
			}
		} else {
			log.Info("Handle HTTP event " + log.Colorize(log.GREY, "PRE_DISPATCH") + " Proxying request " +
				log.Colorize(log.LIGHTMAGENTA, ctx.Request.Method) +
				log.Colorize(log.BLUE, " \""+ctx.Request.URL.String()+"\""))
		}
	case muxy.EVENT_POST_DISPATCH:
		if ctx.Request == nil {
			if len(ctx.Bytes) > 0 {
				log.Info("Handle TCP event " + log.Colorize(log.GREY, "POST_DISPATCH") + " Returning " + string(ctx.Bytes))
			}
		} else {
			log.Info("Handle HTTP event " + log.Colorize(log.GREY, "PRE_DISPATCH") + " Returning " +
				log.Colorize(log.GREY, ctx.Response.Status))
		}
	}
}
