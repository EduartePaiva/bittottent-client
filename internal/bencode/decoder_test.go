package bencode

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeString(t *testing.T) {
	tests := []struct {
		name     string
		sample   []byte
		expected string
	}{
		{"test ascii world", []byte("7:example"), "example"},
		{"test emoji world", []byte("11:example😀"), "example😀"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoder := decoder{*bufio.NewReader(bytes.NewReader(tt.sample))}

			str, err := decoder.readString()
			t.Log(str)
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, str)
		})
	}

}
