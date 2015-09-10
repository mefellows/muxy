package symptom

import (
	"bytes"
	"github.com/mefellows/muxy/log"
	"github.com/mefellows/muxy/muxy"
	"github.com/mefellows/plugo/plugo"
	"io"
	"net/http"
	//"net/http/httptest"
	"strings"
	"time"
)

type RequestConfig struct {
	Method  string
	Headers map[string]string
	Cookies []http.Cookie
	Body    string
}

type ResponseConfig struct {
	Headers map[string]string
	Cookies []http.Cookie
	Body    string
	Status  int
}

type HttpTampererSymptom struct {
	Request  RequestConfig
	Response ResponseConfig
}

func init() {
	plugo.PluginFactories.Register(func() (interface{}, error) {
		return &HttpTampererSymptom{}, nil
	}, "http_tamperer")
}

func (m HttpTampererSymptom) Setup() {
	log.Debug("HTTP Error Setup()")
}

func (m HttpTampererSymptom) Teardown() {
	log.Debug("HTTP Error Teardown()")
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

func (m *HttpTampererSymptom) HandleEvent(e muxy.ProxyEvent, ctx *muxy.Context) {
	switch e {
	case muxy.EVENT_PRE_DISPATCH:
		m.MuckRequest(ctx)
	case muxy.EVENT_POST_DISPATCH:
		m.MuckResponse(ctx)
	}
}

func (m *HttpTampererSymptom) MuckRequest(ctx *muxy.Context) {

	// Body
	if m.Request.Body != "" {
		newreq, err := http.NewRequest(ctx.Request.Method, ctx.Request.URL.String(), bytes.NewBuffer([]byte(m.Request.Body)))
		if err != nil {
			log.Error(err.Error())
		}
		*ctx.Request = *newreq
		log.Debug("Spoofing HTTP Request Body with %s", log.Colorize(log.BLUE, m.Request.Body))
	}

	// Set Cookies
	for _, c := range m.Request.Cookies {
		c.Expires = stringToDate(c.RawExpires)
		log.Debug("Spoofing Request Cookie %s => %s", log.Colorize(log.LIGHTMAGENTA, c.Name), c.String())
		ctx.Request.Header.Add("Cookie", c.String())
	}

	// Set Headers
	for k, v := range m.Request.Headers {
		key := strings.ToTitle(strings.Replace(k, "_", "-", -1))
		log.Debug("Spoofing Request Header %s => %s", log.Colorize(log.LIGHTMAGENTA, key), v)
		ctx.Request.Header.Set(key, v)
	}

	// This Writes all headers, setting status code - so call this last
	if m.Request.Method != "" {
		ctx.Request.Method = m.Request.Method
	}
}
func (m *HttpTampererSymptom) MuckResponse(ctx *muxy.Context) {

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
		log.Debug("Injecting HTTP Response Body with %s", log.Colorize(log.BLUE, m.Response.Body))
		*ctx.Response = *r
	}

	// Set Cookies
	for _, c := range m.Response.Cookies {
		c.Expires = stringToDate(c.RawExpires)
		log.Debug("Spoofing Response Cookie %s => %s", log.Colorize(log.LIGHTMAGENTA, c.Name), c.String())
		ctx.Response.Header.Add("Set-Cookie", c.String())
	}

	// Set Headers
	for k, v := range m.Response.Headers {
		key := strings.ToTitle(strings.Replace(k, "_", "-", -1))
		log.Debug("Spoofing Response Header %s => %s", log.Colorize(log.LIGHTMAGENTA, key), v)
		ctx.Response.Header.Add(key, v)
	}

	// This Writes all headers, setting status code - so call this last
	if m.Response.Status != 0 {
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
