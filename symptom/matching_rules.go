package symptom

import (
	"regexp"

	"github.com/mefellows/muxy/log"
	"github.com/mefellows/muxy/muxy"
)

// HTTPMatchingRule describes the fields to match on an HTTP request
type HTTPMatchingRule struct {
	Method string
	Path   string
	Host   string
}

// MatchHTTPSymptom takes a matching rule and a Muxy context and determines
// if there is a match
func MatchHTTPSymptom(rule HTTPMatchingRule, ctx muxy.Context) bool {
	log.Trace("MatchHTTPSymptom testing rule %v", rule)

	if rule.Path != "" {
		log.Debug("HTTPMatchingRule matching path '%s' with '%s'", rule.Path, ctx.Request.URL.Path)
		if match, _ := regexp.MatchString(rule.Path, ctx.Request.URL.Path); !match {
			return false
		}
	}

	if rule.Host != "" {
		log.Debug("HTTPMatchingRule matching host '%s' with '%s'", rule.Host, ctx.Request.Host)
		if match, _ := regexp.MatchString(rule.Host, ctx.Request.Host); !match {
			return false
		}
	}

	if rule.Method != "" {
		log.Debug("HTTPMatchingRule matching method '%s' with '%s'", rule.Method, ctx.Request.Method)
		if match, _ := regexp.MatchString(rule.Method, ctx.Request.Method); !match {
			return false
		}
	}
	return true
}

// MatchHTTPSymptoms takes a set of matching rules and a Muxy context and determines
// if there is a match
var MatchHTTPSymptoms = func(rules []HTTPMatchingRule, ctx muxy.Context) bool {
	for _, rule := range rules {
		if MatchHTTPSymptom(rule, ctx) {
			return true
		}
	}
	return false
}
