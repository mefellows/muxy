package symptom

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/mefellows/muxy/muxy"
)

func TestMatchSymptom_Hit(t *testing.T) {
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
	subPathMatchingRule := MatchingRule{
		Path: "/foo/",
	}
	hostMatchingRule := MatchingRule{
		Host: "foo\\.com",
	}
	methodMatchingRule := MatchingRule{
		Method: "(?i)get",
	}
	allMatchingRule := MatchingRule{
		Path:   "/foo/bar",
		Host:   ".*foo.*",
		Method: "(?i)get",
	}

	testCases := map[MatchingRule]bool{
		subPathMatchingRule: true,
		hostMatchingRule:    true,
		methodMatchingRule:  true,
		allMatchingRule:     true,
	}

	for rule, expected := range testCases {
		if MatchSymptom(rule, ctx) != expected {
			t.Fatal("Expected", expected, ", got", !expected)
		}
	}

	for rule, expected := range testCases {
		if MatchSymptoms([]MatchingRule{rule}, ctx) != expected {
			t.Fatal("Expected", expected, ", got", !expected)
		}
	}
}

func TestMatchSymptom_Miss(t *testing.T) {
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
	subPathMatchingRule := MatchingRule{
		Path: "^/bar",
	}
	hostMatchingRule := MatchingRule{
		Host: "bar\\.com",
	}
	methodMatchingRule := MatchingRule{
		Method: "(?i)post",
	}
	allMatchingRule := MatchingRule{
		Path:   "^/baz",
		Host:   ".*bar.*",
		Method: "(?i)post",
	}

	testCases := map[MatchingRule]bool{
		subPathMatchingRule: false,
		hostMatchingRule:    false,
		methodMatchingRule:  false,
		allMatchingRule:     false,
	}

	for rule, expected := range testCases {
		if MatchSymptom(rule, ctx) != expected {
			t.Fatal("Expected", expected, ", got", !expected)
		}
	}
	for rule, expected := range testCases {
		if MatchSymptoms([]MatchingRule{rule}, ctx) != expected {
			t.Fatal("Expected", expected, ", got", !expected)
		}
	}
}

func TestProbability(t *testing.T) {
	rand.Seed(time.Now().Unix())
	likelihood := rand.Intn(100)
	fmt.Println(likelihood)
	fmt.Println(int(math.Min(65, 100)))
}
