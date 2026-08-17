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
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, str)
		})
	}

}

func TestDecodeList(t *testing.T) {
	tests := []struct {
		name     string
		sample   []byte
		expected any
	}{
		{"empty list", []byte("le"), []any{}},
		{"list of integers", []byte("li1ei2ei3ei123456789ee"), []any{int64(1), int64(2), int64(3), int64(123456789)}},
		{"list of strings", []byte("l4:spam4:eggse"), []any{"spam", "eggs"}},
		{"list of lists with integer and string", []byte("lli123e4:teste5:testei123ee"), []any{[]any{int64(123), "test"}, "teste", int64(123)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoder := decoder{*bufio.NewReader(bytes.NewReader(tt.sample[1:]))}

			list, err := decoder.readList()
			t.Log(list)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, list)
		})
	}
}

func TestDecodeDictionary(t *testing.T) {
	tests := []struct {
		name     string
		sample   []byte
		expected any
	}{
		{"empty dictionary", []byte("de"), make(map[string]any)},
		{"key value of strings", []byte("d3:cow3:moo4:spam4:eggse"), map[string]any{"cow": "moo", "spam": "eggs"}},
		{"key value of list", []byte("d4:spaml1:a1:bee"), map[string]any{"spam": []any{"a", "b"}}},
		{
			"big key values of string",
			[]byte("d9:publisher3:bob17:publisher-webpage15:www.example.com18:publisher.location4:homee"),
			map[string]any{
				"publisher":          "bob",
				"publisher-webpage":  "www.example.com",
				"publisher.location": "home",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoder := decoder{*bufio.NewReader(bytes.NewReader(tt.sample[1:]))}

			list, err := decoder.readDictionary()
			t.Log(list)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, list)
		})
	}
}
