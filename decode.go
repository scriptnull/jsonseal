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

	if vv, ok := v.(Validator); ok {
		err = vv.Validate()
		if err != nil {
			return err
		}
	}

	return err
}

func (dec *Decoder) DecodeValidate(v Validator) error {
	return dec.Decode(v)
}

func (dec *Decoder) Buffered() io.Reader { return dec.d.Buffered() }

func (dec *Decoder) InputOffset() int64 { return dec.d.InputOffset() }

func (dec *Decoder) More() bool { return dec.d.More() }

func (dec *Decoder) Token() (json.Token, error) { return dec.d.Token() }

func (dec *Decoder) WithUnknownFieldSuggestion() *Decoder {
	dec.unknownFieldSuggestion = true
	dec.DisallowUnknownFields()
	return dec
}

func jsonFields(v any) []string {
	val := reflect.ValueOf(v)
	fieldsCount := reflect.Indirect(val).NumField()
	fields := make([]string, 0, fieldsCount)
	for i := 0; i < fieldsCount; i++ {
		field := reflect.Indirect(val).Type().Field(i)
		fieldTag, jsonTagPresent := field.Tag.Lookup("json")
		if !jsonTagPresent && field.IsExported() {
			fields = append(fields, field.Name)
			continue
		}
		if fieldTag == "" || fieldTag == "-" {
			continue
		}
		if split := strings.Split(fieldTag, ","); len(split) > 1 {
			fields = append(fields, split[0])
			continue
		}
		fields = append(fields, fieldTag)
	}
	return fields
}
