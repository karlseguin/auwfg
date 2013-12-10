package validation

import (
	"regexp"
)

type PatternRule struct {
	pattern *regexp.Regexp
}

func Pattern(pattern *regexp.Regexp) Rule {
	return &PatternRule{pattern}
}

var emailPatternRule = Pattern(regexp.MustCompile("\\S+@\\S+\\.\\S+"))

func Email() Rule { return emailPatternRule }

func (r *PatternRule) Verify(value interface{}) bool {
	switch value := value.(type) {
	case string:
		return r.pattern.MatchString(value)
	default:
		return false
	}
}
