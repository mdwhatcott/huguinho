package core

import (
	"testing"

	"github.com/smartystreets/gunit"
)

func TestDraftFilteringHandlerFixture(t *testing.T) {
	gunit.Run(new(DraftFilteringHandlerFixture), t)
}

type DraftFilteringHandlerFixture struct {
	*gunit.Fixture
}

func (this *DraftFilteringHandlerFixture) Setup() {
}

func (this *DraftFilteringHandlerFixture) Test() {
}
