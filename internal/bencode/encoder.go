package bencode

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
)

type encoder struct {
	bytes.Buffer
}

var ErrInvalidType = errors.New("error: invalid encoded type")

func (e *encoder) encodeString(str string) error {
	_, err := fmt.Fprintf(e, "%d:%s", len(str), str)

	return err
}

func (e *encoder) encodeInt(n int64) error {
	_, err := fmt.Fprintf(e, "i%de", n)

	return err
}

func (e *encoder) encodeAnyType(anyType any) error {
	switch v := anyType.(type) {
	case string:
		return e.encodeString(v)
	case int64:
		return e.encodeInt(v)
	case []any:
		return e.encodeList(v)
	case map[string]any:
		return e.encodeDictionary(v)
	default:
		return fmt.Errorf("error: invalid type'%T'", v)
	}
}

func (e *encoder) encodeList(list []any) error {
	err := e.WriteByte('l')
	if err != nil {
		return err
	}

	for _, value := range list {
		err := e.encodeAnyType(value)
		if err != nil {
			return err
		}
	}

	err = e.WriteByte('e')
	return err
}

func (e *encoder) encodeDictionary(dictionary map[string]any) error {

	sortedKeys := make(sort.StringSlice, 0, len(dictionary))

	for key := range dictionary {
		sortedKeys = append(sortedKeys, key)
	}

	sortedKeys.Sort()

	err := e.WriteByte('d')
	if err != nil {
		return err
	}

	for _, key := range sortedKeys {
		err := e.encodeString(key)
		if err != nil {
			return err
		}
		err = e.encodeAnyType(dictionary[key])
		if err != nil {
			return err
		}
	}

	err = e.WriteByte('e')

	return err
}

func Encode(value any) ([]byte, error) {
	encoder := encoder{}
	err := encoder.encodeAnyType(value)

	return encoder.Bytes(), err
}
