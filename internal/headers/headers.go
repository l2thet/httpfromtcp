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
		return len(parts), true, nil
	}

	containsCRLF, err := regexp.MatchString(`\r\n`, parts)
	if err != nil {
		return 0, false, errors.New("unable to parse headers")
	}
	if !containsCRLF {
		return 0, false, nil
	}

	validKey, err := regexp.MatchString(`\s*[A-Za-z]*:`, parts)
	if err != nil {
		return 0, false, errors.New("unable to parse headers")
	}
	if !validKey {
		return 0, false, errors.New("invalid header key")
	}

	before, after, found := strings.Cut(parts, ":")
	if !found {
		return 0, false, errors.New("invalid header key")
	}
	before = strings.TrimSpace(before)
	after = strings.TrimSpace(after)
	if before == "" || after == "" {
		return 0, false, errors.New("invalid header key")
	}
	h[before] = after
	parsedData := before + ": " + after + "\r\n"
	n = len(parsedData)

	return n, true, nil
}
