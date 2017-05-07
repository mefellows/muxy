package symptom

import (
	"testing"

	"github.com/mefellows/muxy/muxy"
)

func simpleDelay() HTTPDelaySymptom {
	return HTTPDelaySymptom{
		RequestDelay:  1,
		ResponseDelay: 1,
	}
}
func TestHTTPDelay_Muck(t *testing.T) {
	ctx := &muxy.Context{}
	delay := simpleDelay()
	delay.Muck(ctx, 1)
}

func TestHTTPDelayHandleEvent_Hit(t *testing.T) {
	oldMatchHTTPSymptoms := MatchHTTPSymptoms
	MatchHTTPSymptoms = func(rules []HTTPMatchingRule, ctx muxy.Context) bool {
		return true
	}
	defer func() {
		MatchHTTPSymptoms = oldMatchHTTPSymptoms
	}()

	ctx := &muxy.Context{}
	delay := simpleDelay()
	delay.HandleEvent(muxy.EventPreDispatch, ctx)
	delay.HandleEvent(muxy.EventPostDispatch, ctx)

	oneSecondInMillis = 1
	delay = HTTPDelaySymptom{
		Delay: 1,
	}
	delay.HandleEvent(muxy.EventPostDispatch, ctx)
}

func TestHTTPDelayHandleEvent_Miss(t *testing.T) {
	oldMatchHTTPSymptoms := MatchHTTPSymptoms
	MatchHTTPSymptoms = func(rules []HTTPMatchingRule, ctx muxy.Context) bool {
		return false
	}
	defer func() {
		MatchHTTPSymptoms = oldMatchHTTPSymptoms
	}()

	ctx := &muxy.Context{}
	delay := simpleDelay()
	delay.HandleEvent(muxy.EventPreDispatch, ctx)
}

func TestHTTPDelay_Setup(t *testing.T) {
	delay := simpleDelay()
	delay.Setup()

	if len(delay.MatchingRules) != 1 {
		t.Fatal("Expected default HTTPMatchingRule to be present")
	}

	delay = HTTPDelaySymptom{
		MatchingRules: []HTTPMatchingRule{
			HTTPMatchingRule{
				Path: "/foo",
			},
		},
	}

	if len(delay.MatchingRules) != 1 && delay.MatchingRules[0].Path == "/foo" {
		t.Fatal("Expected default ProxyRules not to be present")
	}
}

func TestHTTPDelay_Teardown(t *testing.T) {
	delay := HTTPDelaySymptom{}
	delay.Teardown()
}
