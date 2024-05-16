package jsonseal_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/scriptnull/jsonseal"
)

func TestDecoder(t *testing.T) {
	tt := []struct {
		input    string
		decoded  Data
		expected Data
	}{
		{
			input:   `{ "balance": 50 }`,
			decoded: Data{},
			expected: Data{
				Balance: 50,
			},
		},
	}

	for _, tc := range tt {
		err := jsonseal.NewDecoder(strings.NewReader(tc.input)).Decode(&tc.decoded)
		if err != nil {
			t.Error(err)
		}
	}
}

type Data struct {
	ExpiresIn      int    `json:"expires_in"`
	Balance        int    `json:"balance,omitempty"`
	unexposedField string //nolint:unused
	PrivateField   string `json:"-"`
}

func (d *Data) Validate() error {
	var grp jsonseal.CheckGroup

	grp.Check(func() error {
		if d.ExpiresIn < 0 {
			return errors.New("should have non-zero expiry option")
		}
		return nil
	})

	return grp.Validate()
}

func TestDecoderWithUnknownFieldSuggestion(t *testing.T) {
	t.Run("simple json struct tag", func(t *testing.T) {

		rawData := `{ "expires": 50 }`
		expectedError := `{"errors":[{"fields":["expires"],"error":"unknown field. Did you mean \"expires_in\""}]}`
		var d Data
		err := jsonseal.NewDecoder(strings.NewReader(rawData)).WithUnknownFieldSuggestion().Decode(&d)
		if jsonseal.JSONFormat(err) != expectedError {
			t.Errorf("expected: %s, got %s", expectedError, jsonseal.JSONFormat(err))
		}
	})

	t.Run("Exported field in struct", func(t *testing.T) {
		type Data2 struct {
			Data
			ExportedField string
		}
		rawData := `{ "exported": 50 }`
		expectedError := `{"errors":[{"fields":["exported"],"error":"unknown field. Did you mean \"ExportedField\""}]}`
		var d Data
		err := jsonseal.NewDecoder(strings.NewReader(rawData)).WithUnknownFieldSuggestion().Decode(&d)
		if jsonseal.JSONFormat(err) != expectedError {
			t.Errorf("expected: %s, got %s", expectedError, jsonseal.JSONFormat(err))
		}
	})
}
