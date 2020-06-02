package validate

import "desafio-b2w/errors"

// RequiredString represents required string validation.
type RequiredString struct {
	s string
	e string
}

// NewRequiredString returns a new required string validation.
func NewRequiredString(s string, errMsg string) RequiredString {
	return RequiredString{s, errMsg}
}

// Validate executes validation.
func (rs RequiredString) Validate() error {
	if rs.s == "" {
		return errors.InputError(rs.e)
	}
	return nil
}
