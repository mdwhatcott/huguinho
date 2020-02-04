package core

import (
	"testing"

	"github.com/smartystreets/gunit"
)

func TestMetadataParserFixture(t *testing.T) {
	gunit.Run(new(MetadataParserFixture), t)
}

type MetadataParserFixture struct {
	*gunit.Fixture
	parser *JSONMetadataParser
}

func (this *MetadataParserFixture) Setup() {
}

func (this *MetadataParserFixture) Test() {
}
