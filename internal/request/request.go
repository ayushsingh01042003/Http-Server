package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

type Request struct {
	RequestLine RequestLine
	State       parserState
}

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

type parserState string

const (
	StateInit  parserState = "initialized"
	StateDone  parserState = "done"
	StateError parserState = "error"
)

func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := min(cr.pos+cr.numBytesPerRead, len(cr.data))
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n

	return n, nil
}

func newRequest() *Request {
	return &Request{
		State: StateInit,
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	buf := make([]byte, 1024)
	bufLen := 0

	for !request.done() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			return nil, err
		}

		bufLen += n
		bytesParsed, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[bytesParsed:bufLen])
		bufLen -= bytesParsed
	}

	return request, nil
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	separator := []byte("\r\n")

	index := bytes.Index(b, separator)
	if index == -1 {
		return nil, 0, nil
	}

	var rl RequestLine
	startLine := b[:index]
	bytesRead := index

	rlArr := bytes.Split(startLine, []byte(" "))
	if len(rlArr) != 3 {
		return nil, 0, errors.New("invalid request line size")
	}

	for _, c := range rlArr[0] {
		if c < 65 || c > 96 {
			return nil, 0, errors.New("invalid request line method")
		}
	}

	httpVersionParts := bytes.Split(rlArr[2], []byte("/"))
	if string(httpVersionParts[0]) != "HTTP" || string(httpVersionParts[1]) != "1.1" || len(httpVersionParts) != 2 {
		return nil, 0, errors.New("invalid http version")
	}

	rl.Method = string(rlArr[0])
	rl.RequestTarget = string(rlArr[1])
	rl.HttpVersion = string(httpVersionParts[1])

	return &rl, bytesRead, nil
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		switch r.State {
		case StateError:
			return 0, errors.New("error state")
		case StateInit:
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				r.State = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n
			r.State = StateDone
		case StateDone:
			break outer
		}
	}

	return read, nil
}

func (r *Request) done() bool {
	return r.State == StateDone || r.State == StateError
}

func (r *RequestLine) PrintData() {
	fmt.Println("METHOD:", r.Method)
	fmt.Println("VERSION:", r.HttpVersion)
	fmt.Println("TARGET:", r.RequestTarget)
}