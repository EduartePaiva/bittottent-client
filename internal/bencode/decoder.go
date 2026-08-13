package bencode

import (
	"bufio"
	"io"
	"strconv"
)

type decoder struct {
	bufio.Reader
}

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
