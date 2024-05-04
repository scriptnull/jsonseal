package jsonseal

// Validate could be used to perform the validations and get the validation errors if any
func Validate(v Validator) error {
	return v.Validate()
}

// Validator is the interface that wraps the Validate method
type Validator interface {
	Validate() error
}
