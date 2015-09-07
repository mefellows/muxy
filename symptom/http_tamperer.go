package symptom

import (
	"fmt"
	"github.com/mefellows/muxy/log"
	"github.com/mefellows/muxy/muxy"
	"github.com/mefellows/plugo/plugo"
	"io"
	"net/http"
	"strings"
)

// 50x, 40x etc.

type HttpTampererSymptom struct {
	Status  int `required:"true" default:"500"`
	Headers map[string]string
	Body    string
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
	case muxy.EVENT_POST_DISPATCH:
		log.Debug("Spoofing status code, responding with %s", log.Colorize(log.YELLOW, fmt.Sprintf("%d %s", m.Status, http.StatusText(m.Status))))

		// Body
		if m.Body != "" {
			var cl io.ReadCloser
			cl = &responseBody{body: []byte(m.Body)}
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
			log.Debug("Injecting HTTP Response Body with %s", log.Colorize(log.BLUE, m.Body))
			*ctx.Response = *r
		}

		// Set Headers
		for k, v := range m.Headers {
			key := strings.ToTitle(strings.Replace(k, "_", "-", -1))
			log.Debug("Spoofing Header %s => %s", log.Colorize(log.LIGHTMAGENTA, key), v)
			ctx.Response.Header.Add(key, v)
			ctx.ResponseWriter.Header().Add(key, v)
		}

		// This Writes all headers, setting status code - so call this last
		ctx.Response.StatusCode = m.Status
		ctx.Response.Status = http.StatusText(m.Status)
		ctx.ResponseWriter.WriteHeader(m.Status)
	}
}
