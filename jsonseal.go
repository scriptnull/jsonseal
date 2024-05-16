package jsonseal

import (
	"encoding/json"
)

// Validate could be used to perform the validations and get the validation errors if any
func Validate(v Validator) error {
	return v.Validate()
}

// Validator is the interface that wraps the Validate method
type Validator interface {
	Validate() error
}

// Unmarshal is a drop-in replacement for the standard library json.Unmarshal
func Unmarshal(data []byte, v Validator) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	err = v.Validate()
	if err != nil {
		return err
	}

	return nil
}
