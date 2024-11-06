package util

import (
	"runtime"
	"strings"
	"unicode"
)

// Normalize process strings with unicode characters
// to handle weird behaviors on specific platforms.
// Mainly for windows right now.
func Normalize(text string) string {
	if runtime.GOOS != "windows" {
		return text
	}
	return strings.Map(func(r rune) rune {
		if unicode.IsMark(r) {
			return -1
		}
		return r
	}, text)
}
