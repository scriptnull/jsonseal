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
// But performs jsonseal validations if the input implements the [jsonseal.Validator] interface
func Unmarshal(data []byte, v any) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	if vv, ok := v.(Validator); ok {
		err = vv.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

// UnmarshalValidate is like [jsonseal.Unmarshal] but helps to ensure that the input
// implements the [jsonseal.Validator] interface at compile time
func UnmarshalValidate(data []byte, v Validator) error {
	return Unmarshal(data, v)
}
