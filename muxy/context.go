package muxy

import (
	"net/http"
)

// Context is the request context given to Middlewares and Symptoms.
type Context struct {
	Request        *http.Request
	Response       *http.Response
	ResponseWriter http.ResponseWriter
	Bytes          []byte
}
