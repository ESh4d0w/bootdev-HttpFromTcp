package headers

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	i := bytes.Index(data, []byte(crlf))
	if i == -1 {
		return 0, false, nil
	}

	if i == 0 {
		// empty line
		return 2, true, nil
	}

	parts := bytes.SplitN(data[:i], []byte(":"), 2)

	key, err := validateAndTrimKey(string(parts[0]))
	if err != nil {
		return 0, false, fmt.Errorf("Invalid Key: %s", err)
	}

	value := bytes.TrimSpace(parts[1])
	h[key] = string(value)

	return i + 2, false, nil
}

func validateAndTrimKey(key string) (string, error) {
	if key != strings.TrimRight(key, " ") {
		return "", fmt.Errorf("Space between key and : %s", key)
	}
	key = strings.ToLower(strings.TrimSpace(key))

	var tokenChars = []byte{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}
	for _, c := range key {
		if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') {
			continue
		}
		if slices.Contains(tokenChars, byte(c)) {
			continue
		}
		return "", fmt.Errorf("Invalid Characters: %s", key)
	}

	return key, nil
}
