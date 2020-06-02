package validate

// Conditional represents a conditional validation.
// Only executes another validation when condition is satisfied.
type Conditional struct {
	cond bool
	vals []Validater
}

// NewConditional returns new conditional validation.
func NewConditional(condition bool, validators ...Validater) Conditional {
	return Conditional{condition, validators}
}

// Validate executes validation.
func (c Conditional) Validate() error {
	if c.cond {
		for _, val := range c.vals {
			return val.Validate()
		}
	}
	return nil
}
