// Package symptom contains the default muxy.Symptom implementations:
// HTTP Delay, HTTP Tamperer, Network Shape and the TCP Tamperer.
package symptom

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"strings"
	"time"

	"github.com/mefellows/muxy/log"
	"github.com/mefellows/muxy/muxy"
	"github.com/mefellows/plugo/plugo"
)

// RequestConfig contains details of the HTTP request to tamper with prior to
// sending on to the target system
type RequestConfig struct {
	Method  string
	Headers map[string]string
	Cookies []http.Cookie
	Body    string
	Path    string
	Host    string
}

// ResponseConfig contains details of the HTTP response to tamper with prior to
// sending on to the initiating system
type ResponseConfig struct {
	Headers map[string]string
	Cookies []http.Cookie
	Body    string
	Status  int
}

// HTTPTampererSymptom is a plugin to mess with request/responses between
// a consumer and provider system
type HTTPTampererSymptom struct {
	Request       RequestConfig
	Response      ResponseConfig
	MatchingRules []MatchingRule `required:"false" mapstructure:"matching_rules"`
}

func init() {
	plugo.PluginFactories.Register(func() (interface{}, error) {
		return &HTTPTampererSymptom{}, nil
	}, "http_tamperer")
}

var defaultMatchingRule = MatchingRule{
	Path:   "/",
	Host:   ".*",
	Method: ".*",
}

// Setup sets up the plugin
func (m *HTTPTampererSymptom) Setup() {
	log.Debug("HTTP Tamperer Setup()")

	// Add default (catch all) matching rule
	// Only applicable if none supplied
	if len(m.MatchingRules) == 0 {
		m.MatchingRules = []MatchingRule{
			defaultMatchingRule,
		}
	}
}

// Teardown shuts down the plugin
func (m *HTTPTampererSymptom) Teardown() {
	log.Debug("HTTP Tamperer Teardown()")
}

// Crude implementation of an io.ReadCloser
type responseBody struct {
	body   []byte
	closed bool
}

func (r *responseBody) Close() error {
	r.closed = true
	return nil
}

func (r *responseBody) Read(p []byte) (int, error) {
	if r.closed {
		return 0, io.EOF
	}
	copy(p, r.body)
	r.closed = true
	return len(r.body), nil
}

// HandleEvent is a hook to allow the plugin to intervene with a request/response
// event
func (m *HTTPTampererSymptom) HandleEvent(e muxy.ProxyEvent, ctx *muxy.Context) {
	fmt.Println("MatchingRules:", m.MatchingRules)
	if MatchSymptoms(m.MatchingRules, *ctx) {
		log.Trace("HTTP Tamperer Symptom Hit")
		switch e {
		case muxy.EventPreDispatch:
			m.MuckRequest(ctx)
		case muxy.EventPostDispatch:
			m.MuckResponse(ctx)
		}
	} else {
		log.Trace("HTTP Tamperer Symptom Miss")
	}
}

// MuckRequest adds chaos to the request
func (m *HTTPTampererSymptom) MuckRequest(ctx *muxy.Context) {
	if ctx.Request != nil {

		// Path
		if m.Request.Path != "" {
			log.Debug("HTTP Tamperer Spoofing HTTP URI from [%s] to [%s]", ctx.Request.URL.Path, log.Colorize(log.BLUE, m.Request.Path))
			ctx.Request.URL.Path = m.Request.Path
		}

		// Host
		if m.Request.Host != "" {
			log.Debug("HTTP Tamperer Spoofing HTTP Host from [%s] to [%s]", ctx.Request.URL.Host, log.Colorize(log.BLUE, m.Request.Host))
			ctx.Request.Host = m.Request.Host
			ctx.Request.URL.Host = m.Request.Host
		}

		// Body
		if m.Request.Body != "" {
			newreq, err := http.NewRequest(ctx.Request.Method, ctx.Request.URL.String(), bytes.NewBuffer([]byte(m.Request.Body)))
			if err != nil {
				log.Error(err.Error())
			}
			*ctx.Request = *newreq
			log.Debug("HTTP Tamperer Spoofing HTTP Request Body With [%s]", log.Colorize(log.BLUE, m.Request.Body))
		}

		// Set Cookies
		for _, c := range m.Request.Cookies {
			c.Expires = stringToDate(c.RawExpires)
			log.Debug("HTTP Tamperer Spoofing Request Cookie [%s]=> %s]", log.Colorize(log.LIGHTMAGENTA, c.Name), c.String())
			ctx.Request.Header.Add("Cookie", c.String())
		}

		// Set Headers
		for k, v := range m.Request.Headers {
			key := strings.ToTitle(strings.Replace(k, "_", "-", -1))
			log.Debug("HTTP Tamperer Spoofing Request Header [%s => %s]", log.Colorize(log.LIGHTMAGENTA, key), v)
			ctx.Request.Header.Set(key, v)
		}

		// This Writes all headers, setting status code - so call this last
		if m.Request.Method != "" {
			log.Debug("HTTP Tamperer Spoofing Request Method from [%s] to [%s]", ctx.Request.Method, log.Colorize(log.LIGHTMAGENTA, m.Request.Method))
			ctx.Request.Method = m.Request.Method
		}
	}
}

// MuckResponse adds chaos to the response
func (m *HTTPTampererSymptom) MuckResponse(ctx *muxy.Context) {

	// Body
	if m.Response.Body != "" {
		var cl io.ReadCloser
		cl = &responseBody{body: []byte(m.Response.Body)}
		r := &http.Response{
			Request:          ctx.Request,
			Header:           ctx.Response.Header,
			Close:            ctx.Response.Close,
			ContentLength:    ctx.Response.ContentLength,
			Trailer:          ctx.Response.Trailer,
			TLS:              ctx.Response.TLS,
			TransferEncoding: ctx.Response.TransferEncoding,
			Status:           ctx.Response.Status,
			StatusCode:       ctx.Response.StatusCode,
			Proto:            ctx.Response.Proto,
			ProtoMajor:       ctx.Response.ProtoMajor,
			ProtoMinor:       ctx.Response.ProtoMinor,
			Body:             cl,
		}
		log.Debug("HTTP Tamperer Injecting HTTP Response Body with [%s]", log.Colorize(log.BLUE, m.Response.Body))
		*ctx.Response = *r
	}

	// Set Cookies
	for _, c := range m.Response.Cookies {
		c.Expires = stringToDate(c.RawExpires)
		log.Debug("HTTP Tamperer Spoofing Response Cookie [%s => %s]", log.Colorize(log.LIGHTMAGENTA, c.Name), c.String())
		ctx.Response.Header.Add("Set-Cookie", c.String())
	}

	// Set Headers
	for k, v := range m.Response.Headers {
		key := strings.ToTitle(strings.Replace(k, "_", "-", -1))
		log.Debug("HTTP Tamperer Spoofing Response Header [%s => %s]", log.Colorize(log.LIGHTMAGENTA, key), v)
		ctx.Response.Header.Add(key, v)
	}

	// This Writes all headers, setting status code - so call this last
	if m.Response.Status != 0 {
		log.Debug("HTTP Tamperer Spoofing Response Code From [%d] to [%s]", ctx.Response.StatusCode, log.Colorize(log.LIGHTMAGENTA, fmt.Sprintf("%d", m.Response.Status)))
		ctx.Response.StatusCode = m.Response.Status
		ctx.Response.Status = http.StatusText(m.Response.Status)
	}
}

func stringToDate(val string) time.Time {

	exptime, err := time.Parse(time.RFC1123, val)
	if err != nil {
		exptime, err = time.Parse("Mon, 02-Jan-2006 15:04:05 MST", val)
		if err != nil {
			return time.Time{}
		}
	}
	return exptime.UTC()
}
