package validators

import (
	"regexp"
	"strings"
)

var (
	EmailRX = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	UppercaseRX = regexp.MustCompile(`[A-Z]`)
	LowercaseRX = regexp.MustCompile(`[a-z]`)
	NumberRX    = regexp.MustCompile(`[0-9]`)
	SpecialRX   = regexp.MustCompile(`[!@#\$%\^&\*]`)
)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Require(field, value string) {
	if strings.TrimSpace(value) == "" {
		v.Errors[field] = "This field is required"
	}
}

func (v *Validator) MatchPattern(field, value string, pattern *regexp.Regexp, msg string) {
	if !pattern.MatchString(value) {
		v.Errors[field] = msg
	}
}

func (v *Validator) Check(field string, ok bool, msg string) {
	if !ok {
		v.Errors[field] = msg
	}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}
