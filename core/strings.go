package core

import "strings"

func divide(full string, split string) (before, after string) {
	before, after, ok := strings.Cut(full, split)
	if !ok {
		return "", ""
	}
	return strings.TrimSpace(before), strings.TrimSpace(after)
}

func isSpace(c rune) bool      { return c == ' ' }
func isDash(c rune) bool       { return c == '-' }
func isLowerAlpha(c rune) bool { return c >= 'a' && c <= 'z' }
func isNumber(c rune) bool     { return c >= '0' && c <= '9' }
