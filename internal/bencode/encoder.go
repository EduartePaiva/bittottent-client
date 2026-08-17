package bencode

import (
	"errors"
	"fmt"
	"io"
)

type encoder struct {
	io.Writer
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
	default:
		return fmt.Errorf("error: invalid type'%T'", v)
	}
}

func (e *encoder) encodeList(list []any) error {
	_, err := e.Write([]byte("l"))
	if err != nil {
		return err
	}

	for _, value := range list {
		err := e.encodeAnyType(value)
		if err != nil {
			return err
		}
	}

	_, err = e.Write([]byte("e"))
	return err
}

func Encode(w io.Writer, value map[string]any) error {
	return nil
}
