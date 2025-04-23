package request

import (
	"errors"
	"io"
	"regexp"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type ProcessedState int

const (
	Initialized = iota
	Done
)

const BufferSize int = 8

type RequestLine struct {
	HttpVersion    string
	RequestTarget  string
	Method         string
	ProcessedState ProcessedState
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := Request{
		RequestLine: RequestLine{
			Method:         "",
			RequestTarget:  "",
			HttpVersion:    "",
			ProcessedState: Initialized,
		},
	}

	data := make([]byte, 0, BufferSize)
	readToIndex := 0

	for {
		dataBuffer := make([]byte, BufferSize)
		dataMaxSize := cap(data)

		if request.RequestLine.ProcessedState != Done {
			n, err := reader.Read(dataBuffer)
			if err != nil {
				if err == io.EOF {
					request.RequestLine.ProcessedState = Done
					break
				}
				return nil, err
			}
			if n == 0 {
				return nil, errors.New("error: no data read")
			}
			readToIndex += n
			data = append(data, dataBuffer[:n]...)

			if countNonZeroBytes(data)+BufferSize >= dataMaxSize {
				data = append(data, make([]byte, 0, BufferSize)...)
			}

			parsedAmount, err := request.parse(data)
			if err != nil {
				return nil, err
			}
			readToIndex -= parsedAmount
		} else {
			break
		}
	}
	return &request, nil
}

func parseRequestLine(request string) (int, []string, error) {
	parts := []string{}
	requestLength := len(request)
	if requestLength == 0 {
		return requestLength, nil, errors.New("Request is empty")
	}
	lines := strings.Split(request, "\r\n")
	if len(lines) <= 1 {
		return 0, nil, nil
	}

	if len(lines) > 1 {
		requestLine := lines[0]
		requestLineParts := strings.Split(requestLine, " ")
		if len(requestLineParts) != 3 {
			return requestLength, nil, nil
		}

		validVerb, err := regexp.MatchString(`^([A-Z])`, requestLineParts[0])
		if err != nil {
			return requestLength, nil, errors.New("Verb is incorrect")
		}
		validHttpType, err := regexp.MatchString(`^HTTP/1.1`, requestLineParts[2])
		if err != nil {
			return requestLength, nil, errors.New("Only HTTP/1.1 is supported")
		}
		httpTypePart := strings.Split(requestLineParts[2], "/")[1]

		parts = append(parts, requestLineParts[0])
		parts = append(parts, requestLineParts[1])
		parts = append(parts, httpTypePart)

		if validVerb && validHttpType {
			return requestLength, parts, nil
		} else {
			return requestLength, nil, errors.New("invalid verb or http type")
		}
	}
	return requestLength, nil, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.RequestLine.ProcessedState == Done {
		return 0, errors.New("error: trying to read data in a Done state")
	}
	if r.RequestLine.ProcessedState == Initialized {
		parsedAmount, parts, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}
		if parsedAmount == 0 {
			return 0, nil
		}

		r.RequestLine = RequestLine{
			Method:         parts[0],
			RequestTarget:  parts[1],
			HttpVersion:    parts[2],
			ProcessedState: Done,
		}

		return parsedAmount, nil
	}
	return 0, errors.New("error: unknown state")
}

func countNonZeroBytes(data []byte) int {
	count := 0
	for _, b := range data {
		if b != 0 { // Check if the byte is non-zero
			count++
		}
	}
	return count
}
