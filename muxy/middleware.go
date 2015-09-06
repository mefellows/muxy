package muxy

// Middleware's are executed in stacked order before or after a Middleware,
// and are perfect for jobs like instrumentation. They are given a read/write copy
// of the runtime context and are executed synchronously.

type ProxyEvent int

const (
	EVENT_PRE_DISPATCH ProxyEvent = iota
	EVENT_POST_DISPATCH
	EVENT_PRE_RESPONSE
)

type Middleware interface {
	Setup()
	HandleEvent(event ProxyEvent, ctx *Context)
	Teardown()
}

// Register with plugin factory:
//
//func init() {
//  muxy.PluginFactories.Register(NewMiddlewarex, "sypmtomname")
//}
//
// Where NewMiddlewarex is a func that returns a interface{} (Middleware) when called
