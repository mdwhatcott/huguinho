package static

import (
	"testing"

	"github.com/smartystreets/gunit"
)

func TestNothingFixture(t *testing.T) {
	gunit.Run(new(NothingFixture), t)
}

type NothingFixture struct {
	*gunit.Fixture
}

func (this *NothingFixture) Setup() {
}

func (this *NothingFixture) Test() {
}
