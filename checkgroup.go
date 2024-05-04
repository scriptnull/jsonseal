package jsonseal

import "fmt"

// CheckGroup is collection of checker functions that contain validaton rules
type CheckGroup struct {
	c []checker
}

func (v *CheckGroup) Check(validate func() error) {
	v.c = append(v.c, checker{
		f: validate,
	})
}

func (v *CheckGroup) Field(f string) *FieldChain {
	return &FieldChain{
		validator: v,
		fields:    []string{f},
	}
}

func (v *CheckGroup) Fieldf(f string, a ...any) *FieldChain {
	return &FieldChain{
		validator: v,
		fields:    []string{fmt.Sprintf(f, a...)},
	}
}

func (v *CheckGroup) Validate() error {
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

type checker struct {
	f      func() error
	fields []string
}

type FieldChain struct {
	validator *CheckGroup
	fields    []string
}

func (fc *FieldChain) Check(validate func() error) {
	fc.validator.c = append(fc.validator.c, checker{
		f:      validate,
		fields: fc.fields,
	})
}
