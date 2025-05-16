package headers

import (
	"errors"
	"regexp"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	parts := string(data)

	startsWithCRLF, err := regexp.MatchString(`^\r\n`, parts)
	if err != nil {
		return 0, false, errors.New("unable to parse headers")
	}
	if startsWithCRLF {
		return 0, true, nil
	}

	crlfIndex := strings.Index(parts, "\r\n")
	if crlfIndex == -1 {
		return 0, false, nil // No error, just need more data
	}

	parts = parts[:crlfIndex]

	validHeaderSplit, err := regexp.MatchString(`[:]`, parts)
	if err != nil {
		return 0, false, errors.New("unable to parse headers")
	}
	if !validHeaderSplit {
		return 0, false, errors.New("invalid header, header does not containt ':'")
	}

	before, after, found := strings.Cut(parts, ":")
	if !found {
		return 0, false, errors.New("invalid header key")
	}

	validKey, err := regexp.MatchString(`^[a-zA-Z0-9\-!#\$%&'\*\+\.\^_`+"`"+`|\~]+$`, before)
	if err != nil {
		return 0, false, errors.New("unable to parse headers")
	}
	if !validKey {
		return 0, false, errors.New("invalid header key")
	}

	trimmedBefore := strings.TrimSpace(before)
	timmedAfter := strings.TrimSpace(after)
	if trimmedBefore == "" || timmedAfter == "" {
		return 0, false, errors.New("invalid header key")
	}

	if len(before) > 0 && before[len(before)-1] == ' ' {
		return 0, false, errors.New("invalid spacing in header key")
	}

	lowerKey := strings.ToLower(trimmedBefore)

	h[lowerKey] = timmedAfter
	n = len(parts) + 2 // +2 for CRLF

	return n, false, nil
}
