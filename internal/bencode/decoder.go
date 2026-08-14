package bencode

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strconv"
)

type decoder struct {
	bufio.Reader
}

var (
	ErrLeadingZero = errors.New("error: int have leading zero")
)

func (d *decoder) readString() (string, error) {
	bytesNum, err := d.ReadSlice(':')
	if err != nil {
		return "", err
	}

	str := string(bytesNum[:len(bytesNum)-1])

	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return "", err
	}

	buf := make([]byte, n)
	_, err = io.ReadFull(d, buf)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func (d *decoder) readInt() (int64, error) {
	bytesInt, err := d.ReadSlice('e')
	if err != nil {
		return 0, err
	}

	// check leading zeros
	if bytes.Equal(bytesInt, []byte("i-0e")) || (len(bytesInt) > 3 && bytesInt[1] == '0') {
		return 0, ErrLeadingZero
	}

	strInt := string(bytesInt[1 : len(bytesInt)-1])

	n, err := strconv.ParseInt(strInt, 10, 64)

	return n, err
}
