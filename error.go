package jsonseal

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Error struct {
	Fields []string `json:"fields,omitempty"`
	Err    error    `json:"error"`
}

func (e *Error) Error() string {
	if len(e.Fields) > 0 {
		return fmt.Sprintf("%s %s", strings.Join(e.Fields, ","), e.Err.Error())
	}

	return e.Err.Error()
}

func (e *Error) String() string {
	return e.Error()
}

func (e *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Fields []string `json:"fields,omitempty"`
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

func JSONFormat(e error) string {
	errContent, err := json.Marshal(e)
	if err != nil {
		errContent, _ = json.Marshal(&Errors{
			Errs: []Error{
				{
					Err: err,
				},
			},
		})
	}

	return string(errContent)
}

func JSONIndentFormat(e error, prefix string, indent string) string {
	errContent, err := json.MarshalIndent(e, prefix, indent)
	if err != nil {
		errContent, _ = json.MarshalIndent(&Errors{
			Errs: []Error{
				{
					Err: err,
				},
			},
		}, prefix, indent)
	}

	return string(errContent)
}
