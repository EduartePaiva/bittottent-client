package bencode

import "bufio"

type decoder struct {
	bufio.Reader
}
