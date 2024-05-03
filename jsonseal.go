package jsonseal

import (
	"encoding/json"
	"fmt"
	"strings"
)

type checker struct {
	f      func() error
	fields []string
}

type ValidateAll struct {
	c []checker
}

func (v *ValidateAll) Check(validate func() error) {
	v.c = append(v.c, checker{
		f: validate,
	})
}

func (v *ValidateAll) Field(f string) *FieldChain {
	return &FieldChain{
		validator: v,
		fields:    []string{f},
	}
}

func (v *ValidateAll) Fieldf(f string, a ...any) *FieldChain {
	return &FieldChain{
		validator: v,
		fields:    []string{fmt.Sprintf(f, a...)},
	}
}

type FieldChain struct {
	validator *ValidateAll
	fields    []string
}

func (fc *FieldChain) Check(validate func() error) {
	fc.validator.c = append(fc.validator.c, checker{
		f:      validate,
		fields: fc.fields,
	})
}

type Error struct {
	Fields []string `json:"fields"`
	Err    error    `json:"error"`
}

func (e *Error) Error() string {
	var s strings.Builder

	fmt.Fprintf(&s, "error: %s", e.Err)
	if len(e.Fields) > 0 {
		fmt.Fprintf(&s, ", check: %s", strings.Join(e.Fields, ","))
	}

	return s.String()
}

func (e *Error) String() string {
	return e.Error()
}

func (e *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Fields []string `json:"fields"`
		Err    string   `json:"error"`
	}{
		Fields: e.Fields,
		Err:    e.Err.Error(),
	})
}

type Errors struct {
	Errs []Error `json:"errors"`
}

func (errs *Errors) Error() string {
	var s strings.Builder
	for _, e := range errs.Errs {
		fmt.Fprintln(&s, e.Error())
	}
	return s.String()
}

func (errs *Errors) String() string {
	return errs.Error()
}

func (v *ValidateAll) Validate() error {
	errs := make([]Error, 0, len(v.c))
	for _, check := range v.c {
		if err := check.f(); err != nil {
			errs = append(errs, Error{
				Fields: check.fields,
				Err:    err,
			})
		}
	}
	if len(errs) > 0 {
		return &Errors{Errs: errs}
	}
	return nil
}
