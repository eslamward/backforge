package generator

import "strings"

func GenerateValidator() string {
	var b strings.Builder

	b.WriteString(`
	package validate

import (
	"regexp"
	"slices"
)

	`)
	b.WriteString("\nvar EmailRxp = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$`)\n")
	b.WriteString(`
	
type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

func (v *Validator) Valid() bool {

	return len(v.Errors) == 0

}

func (v *Validator) AddError(key, error string) {

	if _, exist := v.Errors[key]; !exist {
		v.Errors[key] = error
	}
}

func (v *Validator) Check(ok bool, key, error string) {
	if ok {
		v.AddError(key, error)
	}
}

func PermittedValue[T comparable](value T, permittedVaues ...T) bool {
	return slices.Contains(permittedVaues, value)
}

func Match(value string, rxp *regexp.Regexp) bool {
	return rxp.MatchString(value)
}

	`)
	return b.String()
}
