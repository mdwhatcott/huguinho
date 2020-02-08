package core

import (
	"errors"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestTraceErrorFixture(t *testing.T) {
	gunit.Run(new(TraceErrorFixture), t)
}

type TraceErrorFixture struct {
	*gunit.Fixture
}

func (this *TraceErrorFixture) Test() {
	gopherErr := errors.New("gophers")
	err := NewStackTraceError(gopherErr)
	this.So(errors.Is(err, gopherErr), should.BeTrue)
}

func (this *TraceErrorFixture) TestNil() {
	var err error
	err = NewStackTraceError(err)
	this.So(err, should.BeNil)
}
