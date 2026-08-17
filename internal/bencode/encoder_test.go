package bencode

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeList(t *testing.T) {
	tests := []struct {
		name     string
		sample   []any
		expected string
	}{
		{"empty list", []any{}, "le"},
		{"list of integers", []any{int64(1), int64(2), int64(3), int64(123456789)}, "li1ei2ei3ei123456789ee"},
		{"list of strings", []any{"spam", "eggs"}, "l4:spam4:eggse"},
		{"list of lists with integer and string", []any{[]any{int64(123), "test"}, "teste", int64(123)}, "lli123e4:teste5:testei123ee"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			encoder := encoder{Writer: &buf}

			err := encoder.encodeList(tt.sample)
			assert.NoError(t, err)

			assert.Equal(t, tt.expected, buf.String())
		})
	}
}
