package symptom

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/mefellows/muxy/muxy"
)

func TestMatchHTTPSymptom_Hit(t *testing.T) {
	ctx := muxy.Context{
		Request: &http.Request{
			URL: &url.URL{
				Path:   "/foo/bar",
				Host:   "foo.com",
				Scheme: "https",
			},
			Host:   "foo.com",
			Method: "GET",
		},
	}
	subPathHTTPMatchingRule := HTTPMatchingRule{
		Path: "/foo/",
	}
	hostHTTPMatchingRule := HTTPMatchingRule{
		Host: "foo\\.com",
	}
	methodHTTPMatchingRule := HTTPMatchingRule{
		Method: "(?i)get",
	}
	allHTTPMatchingRule := HTTPMatchingRule{
		Path:   "/foo/bar",
		Host:   ".*foo.*",
		Method: "(?i)get",
	}

	testCases := map[HTTPMatchingRule]bool{
		subPathHTTPMatchingRule: true,
		hostHTTPMatchingRule:    true,
		methodHTTPMatchingRule:  true,
		allHTTPMatchingRule:     true,
	}

	for rule, expected := range testCases {
		if MatchHTTPSymptom(rule, ctx) != expected {
			t.Fatal("Expected", expected, ", got", !expected)
		}
	}

	for rule, expected := range testCases {
		if MatchHTTPSymptoms([]HTTPMatchingRule{rule}, ctx) != expected {
			t.Fatal("Expected", expected, ", got", !expected)
		}
	}
}

func TestMatchHTTPSymptom_Miss(t *testing.T) {
	ctx := muxy.Context{
		Request: &http.Request{
			URL: &url.URL{
				Path:   "/foo/bar",
				Host:   "foo.com",
				Scheme: "https",
			},
			Host:   "foo.com",
			Method: "GET",
		},
	}
	subPathHTTPMatchingRule := HTTPMatchingRule{
		Path: "^/bar",
	}
	hostHTTPMatchingRule := HTTPMatchingRule{
		Host: "bar\\.com",
	}
	methodHTTPMatchingRule := HTTPMatchingRule{
		Method: "(?i)post",
	}
	allHTTPMatchingRule := HTTPMatchingRule{
		Path:   "^/baz",
		Host:   ".*bar.*",
		Method: "(?i)post",
	}

	testCases := map[HTTPMatchingRule]bool{
		subPathHTTPMatchingRule: false,
		hostHTTPMatchingRule:    false,
		methodHTTPMatchingRule:  false,
		allHTTPMatchingRule:     false,
	}

	for rule, expected := range testCases {
		if MatchHTTPSymptom(rule, ctx) != expected {
			t.Fatal("Expected", expected, ", got", !expected)
		}
	}
	for rule, expected := range testCases {
		if MatchHTTPSymptoms([]HTTPMatchingRule{rule}, ctx) != expected {
			t.Fatal("Expected", expected, ", got", !expected)
		}
	}
}
