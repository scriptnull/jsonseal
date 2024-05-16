package jsonseal_test

import (
	"strings"
	"testing"

	"github.com/scriptnull/jsonseal"
)

func TestDecoder(t *testing.T) {
	tt := []struct {
		input     string
		decodedAt any
		expected  any
	}{
		{
			input: `{ "balance": 50 }`,
			decodedAt: struct {
				Balance int
			}{},
			expected: struct {
				Balance int
			}{
				Balance: 50,
			},
		},
	}

	for _, tc := range tt {
		jsonseal.NewDecoder(strings.NewReader(tc.input)).Decode(&tc.decodedAt)
	}
}
