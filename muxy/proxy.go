package muxy

// Proxy is the interface for a Proxy plugin
type Proxy interface {
	Setup([]Middleware)
	Proxy()
	Teardown()
}
