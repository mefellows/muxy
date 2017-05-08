package symptom

import (
	"math"
	"regexp"

	"math/rand"

	"github.com/mefellows/muxy/log"
	"github.com/mefellows/muxy/muxy"
)

// MatchingRule describes the fields to match on an HTTP request
type MatchingRule struct {
	Method      string
	Path        string
	Host        string
	Probability float64
}

// MatchSymptom takes a matching rule and a Muxy context and determines
// if there is a match
func MatchSymptom(rule MatchingRule, ctx muxy.Context) bool {
	log.Trace("MatchSymptom testing rule %v", rule)

	// HTTP only matching
	// TODO: Rules should be abstracted better so that we can pass them around
	//       without awareness of individual protocols, like TCP or HTTP
	if ctx.Request != nil {

		if rule.Path != "" {
			log.Debug("MatchingRule matching path '%s' with '%s'", rule.Path, ctx.Request.URL.Path)
			if match, _ := regexp.MatchString(rule.Path, ctx.Request.URL.Path); !match {
				return false
			}
		}

		if rule.Host != "" {
			log.Debug("MatchingRule matching host '%s' with '%s'", rule.Host, ctx.Request.Host)
			if match, _ := regexp.MatchString(rule.Host, ctx.Request.Host); !match {
				return false
			}
		}

		if rule.Method != "" {
			log.Debug("MatchingRule matching method '%s' with '%s'", rule.Method, ctx.Request.Method)
			if match, _ := regexp.MatchString(rule.Method, ctx.Request.Method); !match {
				return false
			}
		}
	}

	// All protocols
	if rule.Probability > 0 {
		random := rand.Intn(100)
		log.Debug("MatchingRule assessing probability %.2f against computed %d", rule.Probability, random)
		if random > int(math.Min(rule.Probability, 100)) {
			return false
		}
	}

	return true
}

// MatchSymptoms takes a set of matching rules and a Muxy context and determines
// if there is a match
var MatchSymptoms = func(rules []MatchingRule, ctx muxy.Context) bool {
	for _, rule := range rules {
		if MatchSymptom(rule, ctx) {
			return true
		}
	}
	return false
}
