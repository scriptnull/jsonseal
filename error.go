package jsonseal

import (
	"encoding/json"
	"fmt"
	"strings"
)

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
