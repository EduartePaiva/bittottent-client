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
	if bytes.Equal(bytesInt, []byte("-0e")) || (len(bytesInt) > 2 && bytesInt[1] == '0') {
		return 0, ErrLeadingZero
	}

	strInt := string(bytesInt[:len(bytesInt)-1])

	n, err := strconv.ParseInt(strInt, 10, 64)

	return n, err
}

func (d *decoder) readNextItem() (any, error) {
	itemType, err := d.ReadByte()
	if err != nil {
		return nil, err
	}

	switch itemType {
	case 'i':
		return d.readInt()
	case 'l':
		return d.readList()
	default:
		// default is a string type
		err = d.UnreadByte()
		if err != nil {
			return nil, err
		}

		return d.readString()
	}
}

func (d *decoder) readList() ([]any, error) {
	// we need to read the next value until next byte is 'e'
	var list []any = []any{}
	for {
		v, err := d.ReadByte()
		if err != nil {
			return nil, err
		}
		if v == 'e' {
			break
		}
		err = d.UnreadByte()
		if err != nil {
			return nil, err
		}

		nextItem, err := d.readNextItem()
		if err != nil {
			return nil, err
		}

		list = append(list, nextItem)
	}

	return list, nil
}
