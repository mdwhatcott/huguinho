package core

import (
	"errors"
	"testing"

	"github.com/smartystreets/gunit"
)

func TestTraceErrorFixture(t *testing.T) {
	gunit.Run(new(TraceErrorFixture), t)
}

type TraceErrorFixture struct {
	*gunit.Fixture
}

func (this *TraceErrorFixture) Test() {
	this.Println(NewStackTraceError(errors.New("gophers")))
}
