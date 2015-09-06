package muxy

type Proxy interface {
	Setup([]Middleware)
	Proxy()
	Teardown()
}
