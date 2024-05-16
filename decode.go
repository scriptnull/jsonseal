package jsonseal

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/agnivade/levenshtein"
)

const (
	unknownFieldErrPrefix string = "json: unknown field "
)

// Decoder is a drop-in replacement for the standard library json.Decoder
type Decoder struct {
	d                      *json.Decoder
	unknownFieldSuggestion bool
}

func NewDecoder(r io.Reader) *Decoder {
	d := json.NewDecoder(r)
	return &Decoder{
		d: d,
	}
}

func (dec *Decoder) UseNumber() { dec.d.UseNumber() }

func (dec *Decoder) DisallowUnknownFields() { dec.d.DisallowUnknownFields() }

func (dec *Decoder) Decode(v any) error {
	err := dec.d.Decode(v)
	if err != nil {
		if dec.unknownFieldSuggestion && strings.HasPrefix(err.Error(), unknownFieldErrPrefix) {
			unknownField := strings.TrimPrefix(err.Error(), unknownFieldErrPrefix)
			unknownField = strings.Trim(unknownField, `"`)
			var fieldSuggesion string
			var minDistance *int
			for _, knownField := range jsonFields(v) {
				levDistance := levenshtein.ComputeDistance(unknownField, knownField)
				if minDistance == nil {
					minDistance = &levDistance
				}
				if levDistance <= *minDistance {
					minDistance = &levDistance
					fieldSuggesion = knownField
				}
			}

			if fieldSuggesion != "" {
				return &Errors{
					Errs: []Error{
						{
							Fields: []string{unknownField},
							Err:    fmt.Errorf(`unknown field. Did you mean "%s"`, fieldSuggesion),
						},
					},
				}
			}
		}
	}

	return err
}

func (dec *Decoder) Buffered() io.Reader { return dec.d.Buffered() }

func (dec *Decoder) InputOffset() int64 { return dec.d.InputOffset() }

func (dec *Decoder) More() bool { return dec.d.More() }

func (dec *Decoder) Token() (json.Token, error) { return dec.d.Token() }

func (dec *Decoder) WithUknownFieldSuggestion() *Decoder {
	dec.unknownFieldSuggestion = true
	dec.DisallowUnknownFields()
	return dec
}

func jsonFields(v any) []string {
	val := reflect.ValueOf(v)
	fieldsCount := reflect.Indirect(val).NumField()
	fields := make([]string, 0, fieldsCount)
	for i := 0; i < fieldsCount; i++ {
		field := reflect.Indirect(val).Type().Field(i).Tag.Get("json")
		if field == "" || field == "-" {
			continue
		}
		if split := strings.Split(field, ","); len(split) > 1 {
			fields = append(fields, split[0])
			continue
		}
		fields = append(fields, field)
	}
	return fields
}
