package symptom

import (
	"testing"

	"github.com/mefellows/muxy/muxy"
)

var newRequestBody = "new request body"
var newResponseBody = "new response body"

func defaultTamperer() TCPTampererSymptom {
	tamperer := TCPTampererSymptom{
		Request: TCPRequestConfig{
			Body:      newRequestBody,
			Randomize: false,
			Truncate:  false,
		},
		Response: TCPResponseConfig{
			Body:      newResponseBody,
			Randomize: false,
			Truncate:  false,
		},
	}
	tamperer.Setup()
	return tamperer
}

func TestTCPTamperer_randStringBytesMaskImprSrc(t *testing.T) {
	r := randStringBytesMaskImprSrc(10)
	if len(r) != 10 {
		t.Fatal("Want 10 got", len(r))
	}
}

func TestTCPTamperer_Setup(t *testing.T) {
	tamperer := TCPTampererSymptom{}
	tamperer.Setup()

	if len(tamperer.MatchingRules) != 1 {
		t.Fatal("Expected default MatchingRule to be present")
	}

	tamperer = TCPTampererSymptom{
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

func TestTCPTamperer_Teardown(t *testing.T) {
	tamperer := TCPTampererSymptom{}
	tamperer.Teardown()
}

func TestTCPTamperer_HandleEventPreDispatchWithTCP(t *testing.T) {
	// Body tests
	tamperer := defaultTamperer()
	ctx := &muxy.Context{
		Bytes: []byte("this is a message"),
	}

	tamperer.HandleEvent(muxy.EventPreDispatch, ctx)
	if string(ctx.Bytes) != newRequestBody {
		t.Fatal("Want", newRequestBody, "got", string(ctx.Bytes))
	}

	// Truncate
	tamperer = defaultTamperer()
	ctx = &muxy.Context{
		Bytes: []byte("this is a message"),
	}
	tamperer.Request.Truncate = true
	tamperer.HandleEvent(muxy.EventPreDispatch, ctx)

	want := newRequestBody[:len(newRequestBody)-2]
	if string(ctx.Bytes) != want {
		t.Fatal("Want", want, "got", string(ctx.Bytes))
	}

	// Randomised
	tamperer = defaultTamperer()
	ctx = &muxy.Context{
		Bytes: []byte("this is a message"),
	}
	tamperer.Request.Randomize = true
	tamperer.HandleEvent(muxy.EventPreDispatch, ctx)

	if string(ctx.Bytes) == newRequestBody || len(ctx.Bytes) != len(newRequestBody) {
		t.Fatal("Want a random string of 15 chars, got", newRequestBody)
	}
}

func TestTCPTamperer_HandleEventPostDispatchWithTCP(t *testing.T) {
	// Body tests
	tamperer := defaultTamperer()

	ctx := &muxy.Context{
		Bytes: []byte("this is a message"),
	}

	tamperer.HandleEvent(muxy.EventPostDispatch, ctx)

	// Truncate
	tamperer = defaultTamperer()
	ctx = &muxy.Context{
		Bytes: []byte("this is a message"),
	}
	tamperer.Response.Truncate = true
	tamperer.HandleEvent(muxy.EventPostDispatch, ctx)

	want := newResponseBody[:len(newResponseBody)-2]
	if string(ctx.Bytes) != want {
		t.Fatal("Want", want, "got", string(ctx.Bytes))
	}

	// Randomised
	tamperer = defaultTamperer()
	ctx = &muxy.Context{
		Bytes: []byte("this is a message"),
	}
	tamperer.Response.Randomize = true
	tamperer.HandleEvent(muxy.EventPostDispatch, ctx)

	if string(ctx.Bytes) == newResponseBody || len(ctx.Bytes) != len(newResponseBody) {
		t.Fatal("Want a random string of 15 chars, got", newResponseBody)
	}
}

func TestTCPTamperer_HandleEventMiss(t *testing.T) {
	tamperer := defaultTamperer()
	tamperer.MatchingRules = []MatchingRule{
		MatchingRule{
			Probability: 1,
		},
	}

	ctx := &muxy.Context{
		Bytes: []byte("this is a message that will never be received"),
	}

	tamperer.HandleEvent(muxy.EventPostDispatch, ctx)
}
