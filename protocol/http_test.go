package protocol

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/mefellows/muxy/muxy"
)

func TestMatchRule_Hit(t *testing.T) {
	proxy := HTTPProxy{}
	defaultProxyRule := proxy.defaultProxyRule()
	subPathProxyRule := ProxyRule{
		ProxyRequest: ProxyRequest{
			Path: "/foo/",
		},
		ProxyPass: ProxyPass{},
	}
	hostProxyRule := ProxyRule{
		ProxyRequest: ProxyRequest{
			Host: "foo\\.com",
		},
		ProxyPass: ProxyPass{},
	}
	methodProxyRule := ProxyRule{
		ProxyRequest: ProxyRequest{
			Method: "(?i)get",
		},
		ProxyPass: ProxyPass{},
	}
	allProxyRule := ProxyRule{
		ProxyRequest: ProxyRequest{
			Path:   "/foo/bar",
			Host:   ".*foo.*",
			Method: "(?i)get",
		},
		ProxyPass: ProxyPass{},
	}

	defaultRequest := http.Request{
		URL: &url.URL{
			Path:   "/foo/bar",
			Host:   "foo.com",
			Scheme: "https",
		},
		Method: "GET",
	}

	testCases := map[ProxyRule]http.Request{
		defaultProxyRule: defaultRequest,
		subPathProxyRule: defaultRequest,
		hostProxyRule:    defaultRequest,
		methodProxyRule:  defaultRequest,
		allProxyRule:     defaultRequest,
	}

	for rule, req := range testCases {
		if MatchRule(rule, req) != true {
			t.Fatal("Expected ProxyRule", rule, "to match request", req, "but did not")
		}
	}
}

func TestMatchRule_Miss(t *testing.T) {
	subPathProxyRule := ProxyRule{
		ProxyRequest: ProxyRequest{
			Path: "^/bar",
		},
		ProxyPass: ProxyPass{},
	}
	hostProxyRule := ProxyRule{
		ProxyRequest: ProxyRequest{
			Host: "bar\\.com",
		},
		ProxyPass: ProxyPass{},
	}
	methodProxyRule := ProxyRule{
		ProxyRequest: ProxyRequest{
			Method: "(?i)post",
		},
		ProxyPass: ProxyPass{},
	}
	allProxyRule := ProxyRule{
		ProxyRequest: ProxyRequest{
			Path:   "^/baz",
			Host:   ".*bar.*",
			Method: "(?i)post",
		},
		ProxyPass: ProxyPass{},
	}

	defaultRequest := http.Request{
		URL: &url.URL{
			Path:   "/foo/bar",
			Host:   "foo.com",
			Scheme: "https",
		},
		Method: "GET",
	}

	testCases := map[ProxyRule]http.Request{
		subPathProxyRule: defaultRequest,
		hostProxyRule:    defaultRequest,
		methodProxyRule:  defaultRequest,
		allProxyRule:     defaultRequest,
	}

	for rule, req := range testCases {
		if MatchRule(rule, req) != false {
			t.Fatal("Expected ProxyRule", rule, "to not match request", req, "but did")
		}
	}
}

func TestApplyProxyPassRule_Path(t *testing.T) {
	proxy := HTTPProxy{}
	subPathProxyRule := ProxyRule{
		ProxyRequest: ProxyRequest{},
		ProxyPass: ProxyPass{
			Path: "/newstart",
		},
	}
	defaultRequest := &http.Request{
		URL: &url.URL{
			Path:   "/foo/bar",
			Host:   "foo.com",
			Scheme: "https",
		},
		Method: "GET",
	}
	proxy.ApplyProxyPassRule(subPathProxyRule, defaultRequest)

	if defaultRequest.URL.Path != "/newstart/foo/bar" {
		t.Fatal("Expected URL to be translated to /newstart/foo/bar but got", defaultRequest.URL.Path)
	}
	rootProxyRule := ProxyRule{
		ProxyRequest: ProxyRequest{},
		ProxyPass:    ProxyPass{},
	}
	defaultRequest = &http.Request{
		URL: &url.URL{
			Path:   "/foo/bar",
			Host:   "foo.com",
			Scheme: "https",
		},
		Method: "GET",
	}
	proxy.ApplyProxyPassRule(rootProxyRule, defaultRequest)

	if defaultRequest.URL.Path != "/foo/bar" {
		t.Fatal("Expected URL to be unmodified at /foo/bar but got", defaultRequest.URL.Path)
	}
}

func TestApplyProxyPassRule_Method(t *testing.T) {
	proxy := HTTPProxy{}
	hostProxyRule := ProxyRule{
		ProxyRequest: ProxyRequest{},
		ProxyPass: ProxyPass{
			Method: "POST",
		},
	}
	defaultRequest := &http.Request{
		URL: &url.URL{
			Path:   "/foo/bar",
			Host:   "foo.com",
			Scheme: "https",
		},
		Method: "GET",
	}
	proxy.ApplyProxyPassRule(hostProxyRule, defaultRequest)

	if defaultRequest.Method != "POST" {
		t.Fatal("Expected request method to be POST but got", defaultRequest.Method)
	}
}

func TestApplyProxyPassRule_Scheme(t *testing.T) {
	proxy := HTTPProxy{}
	schemeProxyRule := ProxyRule{
		ProxyRequest: ProxyRequest{},
		ProxyPass: ProxyPass{
			Scheme: "http",
		},
	}
	defaultRequest := &http.Request{
		URL: &url.URL{
			Path:   "/foo/bar",
			Host:   "foo.com",
			Scheme: "https",
		},
		Method: "GET",
	}
	proxy.ApplyProxyPassRule(schemeProxyRule, defaultRequest)

	if defaultRequest.URL.Scheme != "http" {
		t.Fatal("Expected URL scheme to be http but got", defaultRequest.URL.Scheme)
	}
}

func TestApplyProxyPassRule_Host(t *testing.T) {
	proxy := HTTPProxy{}
	hostProxyRule := ProxyRule{
		ProxyRequest: ProxyRequest{},
		ProxyPass: ProxyPass{
			Host: "bar.com",
		},
	}
	defaultRequest := &http.Request{
		URL: &url.URL{
			Path:   "/foo/bar",
			Host:   "foo.com",
			Scheme: "https",
		},
		Method: "GET",
	}
	proxy.ApplyProxyPassRule(hostProxyRule, defaultRequest)

	if defaultRequest.URL.Host != "bar.com" {
		t.Fatal("Expected URL Host to be bar.com but got", defaultRequest.URL.Host)
	}
}

func TestSetup(t *testing.T) {
	proxy := HTTPProxy{}
	proxy.Setup([]muxy.Middleware{})

	if len(proxy.ProxyRules) != 1 {
		t.Fatal("Expected default ProxyRules to be present")
	}

	proxy = HTTPProxy{
		ProxyRules: []ProxyRule{
			ProxyRule{},
		},
	}
	proxy.Setup([]muxy.Middleware{})

	if len(proxy.ProxyRules) != 2 {
		t.Fatal("Expected default ProxyRules to be present")
	}
}

func TestDefaultProxyRule(t *testing.T) {
	proxy := HTTPProxy{
		ProxyHost: "foo.com",
		ProxyPort: 1234,
	}
	rule := proxy.defaultProxyRule()

	expected := "foo.com:1234"
	if rule.ProxyPass.Host != expected {
		t.Fatal("Expected host to be", expected, "but got", rule.ProxyPass.Host)
	}
}
