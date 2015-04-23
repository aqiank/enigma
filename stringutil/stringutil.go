package stringutil

import (
	"strings"
)

// Remove unacceptable characters from message
func Sanitize(msg string) string {
	msg = strings.ToUpper(msg)
	msg = strings.TrimSpace(msg)
	msg = stripChars(msg)
	return msg
}

// Strip unknown characters
func stripChars(str string) string {
	return strings.Map(func(r rune) rune {
		if r < 65 || r > 90 {
			return -1
		}
		return r
	}, str)
}
