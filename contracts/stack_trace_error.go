package contracts

import (
	"fmt"
	"runtime/debug"
)

func StackTraceError(inner error) error {
	if inner == nil {
		return nil
	}
	return fmt.Errorf("error: %w\nstack:\n%s", inner, debug.Stack())
}
