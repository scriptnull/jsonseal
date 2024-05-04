package jsonseal

import (
	"fmt"
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
