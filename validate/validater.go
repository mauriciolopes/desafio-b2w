package validate

// Validater is the interface that wraps the basic Validate method.
type Validater interface {
	Validate() error
}

// Exec returns the first found error or nil if anything.
func Exec(v ...Validater) error {
	for _, validator := range v {
		if err := validator.Validate(); err != nil {
			return err
		}
	}
	return nil
}
