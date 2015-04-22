package stringutil

import (
	"strings"
)

// Remove unacceptable characters from message
func Sanitize(msg string) string {
	msg = strings.ToUpper(msg)
	msg = strings.TrimSpace(msg)
	msg = stripChars(msg, " `1234567890-=~!@#$%^&*()_+[]\\;',./{}|:\"<>?")
	return msg
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
