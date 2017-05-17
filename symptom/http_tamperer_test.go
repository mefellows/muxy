package symptom

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"
	"time"

	"io/ioutil"

	"github.com/mefellows/muxy/muxy"
)

func TestHTTPTampererSymptom_CreateCookieHeader(t *testing.T) {
	cookies := []http.Cookie{
		http.Cookie{
			Domain:   "foo.com",
			Expires:  time.Now(),
			HttpOnly: true,
			MaxAge:   1,
			Name:     "Foo",
			Path:     "/",
			Secure:   false,
			Value:    "blah",
		},
	}
	for _, c := range cookies {
		t.Logf("cookie: %s", c.String())

	}
}

func TestHTTPTampererSymptom_MuckRequest(t *testing.T) {
	tamperer := HTTPTampererSymptom{
		Request: RequestConfig{
			Method: "POST",
			Body:   "my new body",
			Path:   "/my/new/path",
			Host:   "mynewhost.com",
			Headers: map[string]string{
				"MyNewHeader": "MyNewHeader",
			},
			Cookies: []http.Cookie{
				http.Cookie{
					Name:  "MyNewCookie",
					Value: "MyNewValue",
				},
			},
		},
	}
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

	tamperer.MuckRequest(ctx)

	if ctx.Request.Method != "POST" {
		t.Fatal("Expected request to be tampered with a POST request, got", ctx.Request.Method)
	}
	if ctx.Request.URL.Path != "/my/new/path" {
		t.Fatal("Expected request to be tampered to path /my/new/path, got", ctx.Request.URL.Path)
	}
	if ctx.Request.Host != "mynewhost.com" {
		t.Fatal("Expected request to be tampered to host mynewhost.com, got", ctx.Request.Host)
	}
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	if string(body) != "my new body" {
		t.Fatal("Expected request body to be tampered to 'my new body', got", string(body))
	}
	if ctx.Request.Header.Get("MyNewHeader") != "MyNewHeader" {
		t.Fatal("Expected request header to be tampered to contain 'MyNewHeader', got", ctx.Request.Header.Get("MyNewHeader"))
	}
	if _, err := ctx.Request.Cookie("MyNewCookie"); err != nil {
		t.Fatal("Expected request cookie 'MyNewCookie' to be added, got", err)
	}
}

func TestHTTPTampererSymptom_MuckResponse(t *testing.T) {
	cl := ioutil.NopCloser(bytes.NewReader([]byte("my new response body")))
	tamperer := HTTPTampererSymptom{
		Response: ResponseConfig{
			Body: "my new body",
			Headers: map[string]string{
				"MyNewHeader": "MyNewHeader",
			},
			Cookies: []http.Cookie{
				http.Cookie{
					Name:  "MyNewCookie",
					Value: "MyNewValue",
				},
			},
			Status: 200,
		},
	}
	ctx := &muxy.Context{
		Response: &http.Response{
			Body:       cl,
			Status:     "301",
			StatusCode: 301,
			Header: map[string][]string{
				"MyOriginalHeader": []string{"MyOriginalHeader"},
			},
		},
	}

	tamperer.MuckResponse(ctx)

	body, _ := ioutil.ReadAll(ctx.Response.Body)
	if string(body) != "my new body" {
		t.Fatal("Expected response body to be tampered to 'my new body', got", string(body))
	}
	if ctx.Response.Header.Get("MyNewHeader") != "MyNewHeader" {
		t.Fatal("Expected response header to be tampered to contain 'MyNewHeader', got", ctx.Response.Header.Get("MyNewHeader"))
	}
	if ctx.Response.StatusCode != 200 {
		t.Fatal("Expected response status code to be tampered to be 200, got", ctx.Response.StatusCode)
	}
	if ctx.Response.Cookies()[0].Name != "MyNewCookie" {
		t.Fatal("Expected response cookies to contain 'MyNewCookie' but did not")
	}
}

func TestHTTPTampererSymptom_HandleEventPostDispatch(t *testing.T) {
	cl := ioutil.NopCloser(bytes.NewReader([]byte("my response body")))
	tamperer := HTTPTampererSymptom{
		Response: ResponseConfig{
			Body: "my new body",
			Headers: map[string]string{
				"MyNewHeader": "MyNewHeader",
			},
			Cookies: []http.Cookie{
				http.Cookie{
					Name:  "MyNewCookie",
					Value: "MyNewValue",
				},
			},
			Status: 200,
		},
	}
	ctx := &muxy.Context{
		Response: &http.Response{
			Body:       cl,
			Status:     "301",
			StatusCode: 301,
			Header: map[string][]string{
				"MyOriginalHeader": []string{"MyOriginalHeader"},
			},
		},
	}

	tamperer.Setup()
	tamperer.HandleEvent(muxy.EventPostDispatch, ctx)

	body, _ := ioutil.ReadAll(ctx.Response.Body)
	if string(body) != "my new body" {
		t.Fatal("Expected response body to be tampered to 'my new body', got", string(body))
	}
	if ctx.Response.Header.Get("MyNewHeader") != "MyNewHeader" {
		t.Fatal("Expected response header to be tampered to contain 'MyNewHeader', got", ctx.Response.Header.Get("MyNewHeader"))
	}
	if ctx.Response.StatusCode != 200 {
		t.Fatal("Expected response status code to be tampered to be 200, got", ctx.Response.StatusCode)
	}
	if ctx.Response.Cookies()[0].Name != "MyNewCookie" {
		t.Fatal("Expected response cookies to contain 'MyNewCookie' but did not")
	}
}
func TestHTTPTampererSymptom_HandleEventPreDispatch(t *testing.T) {
	tamperer := HTTPTampererSymptom{
		Request: RequestConfig{
			Method: "POST",
			Body:   "my new body",
			Path:   "/my/new/path",
			Host:   "mynewhost.com",
			Headers: map[string]string{
				"MyNewHeader": "MyNewHeader",
			},
			Cookies: []http.Cookie{
				http.Cookie{
					Name:  "MyNewCookie",
					Value: "MyNewValue",
				},
			},
		},
	}
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

	tamperer.Setup()
	tamperer.HandleEvent(muxy.EventPreDispatch, ctx)

	if ctx.Request.Method != "POST" {
		t.Fatal("Expected request to be tampered with a POST request, got", ctx.Request.Method)
	}
	if ctx.Request.URL.Path != "/my/new/path" {
		t.Fatal("Expected request to be tampered to path /my/new/path, got", ctx.Request.URL.Path)
	}
	if ctx.Request.Host != "mynewhost.com" {
		t.Fatal("Expected request to be tampered to host mynewhost.com, got", ctx.Request.Host)
	}
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	if string(body) != "my new body" {
		t.Fatal("Expected request body to be tampered to 'my new body', got", string(body))
	}
	if ctx.Request.Header.Get("MyNewHeader") != "MyNewHeader" {
		t.Fatal("Expected request header to be tampered to contain 'MyNewHeader', got", ctx.Request.Header.Get("MyNewHeader"))
	}
	if _, err := ctx.Request.Cookie("MyNewCookie"); err != nil {
		t.Fatal("Expected request cookie 'MyNewCookie' to be added, got", err)
	}
}

func TestHTTPTampererSymptom_StringToDate(t *testing.T) {
	date := stringToDate("Sun, 07 May 2017 21:22:48 UTC")
	if date.Year() != 2017 {
		t.Fatal("Expected 2017, got", date.Year())
	}
	if date.Minute() != 22 {
		t.Fatal("Expected 22, got", date.Minute())
	}
}

func TestHTTPTampererSymptom_Setup(t *testing.T) {
	tamperer := HTTPTampererSymptom{}
	tamperer.Setup()

	if len(tamperer.MatchingRules) != 1 {
		t.Fatal("Expected default MatchingRule to be present")
	}

	tamperer = HTTPTampererSymptom{
		MatchingRules: []MatchingRule{
			MatchingRule{
				Path: "/foo",
			},
		},
	}

	if len(tamperer.MatchingRules) != 1 && tamperer.MatchingRules[0].Path == "/foo" {
		t.Fatal("Expected default ProxyRules not to be present")
	}
}

func TestHTTPTampererSymptom_Teardown(t *testing.T) {
	tamperer := HTTPTampererSymptom{}
	tamperer.Teardown()
}
