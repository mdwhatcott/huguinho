package core

import (
	"errors"
	"testing"

	"github.com/mdwhatcott/testing/should"
	"github.com/mdwhatcott/testing/suite"
)

func TestTraceErrorFixture(t *testing.T) {
	suite.Run(&StackTraceErrorFixture{T: suite.New(t)}, suite.Options.UnitTests())
}

type StackTraceErrorFixture struct {
	*suite.T
}

func (this *StackTraceErrorFixture) Test() {
	gopherErr := errors.New("gophers")
	err := StackTraceError(gopherErr)
	this.FatalSo(err, should.WrapError, gopherErr)
	this.So(err.Error(), should.Contain, "gophers")
	this.So(err.Error(), should.Contain, "stack:")
}

func (this *StackTraceErrorFixture) TestNil() {
	var err error
	err = StackTraceError(err)
	this.So(err, should.BeNil)
}
