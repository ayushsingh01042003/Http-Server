package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	reqInBytes, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println("Error reading request", err)
	}
	
	requestString := string(reqInBytes)
	rl, _, err := parseRequestLine(requestString)
	if err != nil {
		fmt.Println("Error Parsing RequsetLine", err)
		return nil, err
	}

	return &Request{
		RequestLine: *rl,
	}, nil
}

func parseRequestLine(b string) (*RequestLine, string, error) {
	separator := "\r\n"

	index := strings.Index(b, separator)
	if index == -1 {
		return nil, "", errors.New("invalid request line format")
	}

	var rl RequestLine
	startLine := b[:index]
	restOfMsg := b[index + len(separator):]

	rlArr := strings.Split(startLine, " ")
	if len(rlArr) != 3 {
		return nil, "", errors.New("invalid request line size")
	}

	httpVersionParts := strings.Split(rlArr[2], "/")
	if httpVersionParts[0] != "HTTP" || httpVersionParts[1] != "1.1" || len(httpVersionParts) != 2 {
		return nil, "", errors.New("invalid http version")
	}

	rl.Method = rlArr[0]
	rl.RequestTarget = rlArr[1]
	rl.HttpVersion = httpVersionParts[1]

	return &rl, restOfMsg, nil
}
