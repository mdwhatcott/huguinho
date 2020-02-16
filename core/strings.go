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
