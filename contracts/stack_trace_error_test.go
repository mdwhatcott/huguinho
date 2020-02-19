package contracts

import (
	"errors"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestTraceErrorFixture(t *testing.T) {
	gunit.Run(new(StackTraceErrorFixture), t)
}

type StackTraceErrorFixture struct {
	*gunit.Fixture
}

func (this *StackTraceErrorFixture) Test() {
	gopherErr := errors.New("gophers")
	err := StackTraceError(gopherErr)
	this.So(errors.Is(err, gopherErr), should.BeTrue)
	this.So(err.Error(), should.ContainSubstring, "gophers")
	this.So(err.Error(), should.ContainSubstring, "stack:")
}

func (this *StackTraceErrorFixture) TestNil() {
	var err error
	err = StackTraceError(err)
	this.So(err, should.BeNil)
}
