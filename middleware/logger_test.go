package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/mefellows/muxy/muxy"
)

func TestLogger_Setup(t *testing.T) {
	l := LoggerMiddleware{}
	l.Setup()

	if l.HexOutput == true {
		t.Fatal("Want false got true")
	}
	if l.format != "%s" {
		t.Fatal("Want '%s' got'", l.HexOutput, "'")
	}

	l = LoggerMiddleware{HexOutput: true}
	l.Setup()

	if l.format == "%s" {
		t.Fatal("Want '%x' got'", l.HexOutput, "'")
	}
}

func TestLogger_Teardown(t *testing.T) {
	l := LoggerMiddleware{}
	l.Teardown()
}

func TestLogger_HandleEventPreDispatchWithHTTP(t *testing.T) {
	l := LoggerMiddleware{}

	ctx := &muxy.Context{
		Request: &http.Request{
			URL: &url.URL{
				Path:   "/foo/bar",
				Scheme: "https",
			},
			Host:   "foo.com",
			Method: "GET",
		},
	}

	l.HandleEvent(muxy.EventPreDispatch, ctx)
}

func TestLogger_HandleEventPostDispatchWithHTTP(t *testing.T) {
	l := LoggerMiddleware{}

	cl := ioutil.NopCloser(bytes.NewReader([]byte("my new response body")))

	ctx := &muxy.Context{
		Request: &http.Request{
			URL: &url.URL{
				Path:   "/foo/bar",
				Scheme: "https",
			},
			Host:   "foo.com",
			Method: "GET",
		},
		Response: &http.Response{
			Body:       cl,
			Status:     "301",
			StatusCode: 301,
			Header: map[string][]string{
				"MyOriginalHeader": []string{"MyOriginalHeader"},
			},
		},
	}

	l.HandleEvent(muxy.EventPostDispatch, ctx)
}

func TestLogger_HandleEventPreDispatchWithTCP(t *testing.T) {
	l := LoggerMiddleware{}

	ctx := &muxy.Context{
		Bytes: []byte("this is a message"),
	}

	l.HandleEvent(muxy.EventPreDispatch, ctx)
}
func TestLogger_HandleEventPostDispatchWithTCP(t *testing.T) {
	l := LoggerMiddleware{}

	ctx := &muxy.Context{
		Bytes: []byte("this is a message"),
	}

	l.HandleEvent(muxy.EventPostDispatch, ctx)
}
