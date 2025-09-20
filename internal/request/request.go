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
	HttpVersion   string
	RequestTarget string
	Method        string
}

var ErrorMethodNotFound = errors.New("HTTP method not found")
var ErrorInvalidTarget = errors.New("invalid target")
var ErrrorInvalidHttpVersion = errors.New("invalid HTTP version")
var ErrorInvalidRequestLine = errors.New("invalid request line")
var ErrorUnsupportedHTTPVersion = errors.New("unsupported HTTP version")

const SEPARATOR = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("unable to read io.ReadAll"), err)
	}

	str := string(data)
	i := strings.Index(str, SEPARATOR)
	if i == -1 {
		return nil, ErrorInvalidRequestLine
	}

	requestLine, err := parseRequestLine(str[:i])
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *requestLine,
	}, nil
}

func parseRequestLine(l string) (*RequestLine, error) {
	parts := strings.Split(l, " ")

	if len(parts) != 3 {
		return nil, ErrorInvalidRequestLine
	}

	httpParts := strings.Split(parts[2], "/")
	if len(httpParts) != 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return nil, ErrrorInvalidHttpVersion
	}

	rl := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   httpParts[1],
	}

	return rl, nil
}
