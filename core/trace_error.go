package core

import (
	"fmt"
	"runtime/debug"
)

type StackTraceError struct {
	inner error
	stack []byte
}

func NewStackTraceError(inner error) *StackTraceError {
	return &StackTraceError{
		inner: inner,
		stack: debug.Stack(),
	}
}

func (this *StackTraceError) Error() string {
	return fmt.Sprintf("error: %q\nstack:\n%s", this.inner, this.stack)
}

func (this *StackTraceError) Unwrap() error {
	return this.inner
}
