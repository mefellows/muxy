package muxy

// ProxyEvent is the event set to a proxy
type ProxyEvent int

const (
	// EventPreDispatch is the event sent prior to dispatching a request
	EventPreDispatch ProxyEvent = iota

	// EventPostDispatch is the event sent directly after dispatching a request
	EventPostDispatch
)

// Middleware is a plugin that intercepts requests and injects chaos
// Middleware's are executed in stacked order before or after a Middleware,
// and are perfect for jobs like instrumentation. They are given a read/write copy
// of the runtime context and are executed synchronously.
//
// Middleware's are registered via a plugin factory e.g.
//
//    func init() {
//      muxy.PluginFactories.Register(NewMiddlewarex, "sypmtomname")
//    }
//
// Where NewMiddlewarex is a func that returns a interface{} (Middleware) when called
type Middleware interface {
	// Setup is called when the plugin is registered
	Setup()

	// HandleEvent takes a ProxyEvent and acts on the information provided
	HandleEvent(event ProxyEvent, ctx *Context)

	// Teardown is used to cleanup any resources on plugin destray
	Teardown()
}
