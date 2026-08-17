package bencode

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

type decoder struct {
	bufio.Reader
}

var (
	ErrLeadingZero       = errors.New("error: int have leading zero")
	ErrMustBeDictionary  = errors.New("error: bencode must start as a dictionary")
	ErrInvalidStringSize = errors.New("error: string size is invalid")
)

const maxStringSize = 10 * 1024 * 1024 // 10MiB

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

	if n < 0 || n > maxStringSize {
		return "", ErrInvalidStringSize
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

	strInt := string(bytesInt[:len(bytesInt)-1])
	// check leading zeros
	if len(strInt) > 1 && strInt[0] == '0' {
		return 0, ErrLeadingZero
	}

	if len(strInt) > 1 && strInt[0] == '-' && strInt[1] == '0' {
		return 0, ErrLeadingZero
	}

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
	case 'd':
		return d.readDictionary()
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

func (d *decoder) readDictionary() (map[string]any, error) {
	// we need to read the next value until next byte is 'e'
	dictionary := make(map[string]any)
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

		key, err := d.readString()
		if err != nil {
			return nil, err
		}

		value, err := d.readNextItem()
		if err != nil {
			return nil, err
		}

		dictionary[key] = value
	}

	return dictionary, nil
}

func Decode(reader io.Reader) (map[string]any, error) {
	d := decoder{Reader: *bufio.NewReader(reader)}

	firstByte, err := d.ReadByte()
	if err != nil {
		return nil, err
	}
	if firstByte != 'd' {
		return nil, ErrMustBeDictionary
	}

	return d.readDictionary()
}
