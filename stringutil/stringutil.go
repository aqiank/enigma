package stringutil

import (
	"strings"
)

// Remove unacceptable characters from message
func Sanitize(msg []byte) []byte {
	s := string(msg)
	s = strings.ToUpper(s)
	s = strings.TrimSpace(s)
	s = stripChars(s, " `1234567890-=~!@#$%^&*()_+[]\\;',./{}|:\"<>?")
	return []byte(s)
}

// Strip a set of characters from string
func stripChars(str, chr string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(chr, r) < 0 {
			return r
		}
		return -1
	}, str)
}
