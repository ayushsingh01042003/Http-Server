package headers

import (
	"bytes"
	"errors"
)

type Headers map[string]string

var CRLF = []byte("\r\n")

func NewHeaders() Headers {
	return Headers{}
}

func parserHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", errors.New("incorrect header data")
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", errors.New("whitespace in header key")
	}

	return string(name), string(value), nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false

	for {
		idx := bytes.Index(data[read:], CRLF)
		if idx == -1 {
			break
		}

		// END OF HEADERS
		if idx == 0 {
			done = true
			break
		}

		name, value, err := parserHeader(data[:idx])
		if err != nil {
			return 0, false, err
		}

		read += idx + len(CRLF)
		h[name] = value
	}

	return read, done, nil
}
