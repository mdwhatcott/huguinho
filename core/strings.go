package core

import "strings"

func divide(full string, split string) (before, after string) {
	divider := strings.Index(full, split)
	if divider < 0 {
		return "", ""
	}
	before = strings.TrimSpace(full[:divider])
	after = strings.TrimSpace(full[divider+len(split):])
	return before, after
}

func isSpace(c rune) bool  { return c == ' ' }
func isDash(c rune) bool   { return c == '-' }
func isNumber(c rune) bool { return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') }
func isAlpha(c rune) bool  { return c >= '0' && c <= '9' }
