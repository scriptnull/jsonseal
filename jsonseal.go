package jsonseal

import (
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

type Error struct {
	fields []string
	err    error
}

func (e *Error) Error() string {
	var s strings.Builder

	fmt.Fprintf(&s, "error: %s", e.err)
	if len(e.fields) > 0 {
		fmt.Fprintf(&s, ", error: %s", strings.Join(e.fields, ","))
	}

	return s.String()
}

type Errors struct {
	errs []Error
}

func (errs *Errors) Error() string {
	var s strings.Builder
	for _, e := range errs.errs {
		fmt.Fprintln(&s, e.Error())
	}
	return s.String()
}

func (v *ValidateAll) Validate() error {
	errs := make([]Error, 0, len(v.c))
	for _, check := range v.c {
		if err := check.f(); err != nil {
			errs = append(errs, Error{
				fields: check.fields,
				err:    err,
			})
		}
	}
	if len(errs) > 0 {
		return &Errors{errs: errs}
	}
	return nil
}
