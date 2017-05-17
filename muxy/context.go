package muxy

import (
	"net/http"
)

// Context is the request context given to Middlewares and Symptoms.
type Context struct {
	// Request contains a reference to the HTTP Request
	// if it's an HTTP proxied event.
	// It may be mutated by prior and future middlewares/plugins
	Request *http.Request

	// Response contains a reference to the HTTP Response
	// if it's an HTTP proxied event.
	// It may be mutated by prior and future middlewares/plugins
	Response *http.Response

	// ResponseWriter for the current HTTP session if it exists.
	ResponseWriter http.ResponseWriter

	// Bytes contains the current message for TCP sessions.
	Bytes []byte
}
