package site

import (
	"testing"

	"github.com/mdwhatcott/static/contracts"
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/nu"
)

func TestIndexingFixture(t *testing.T) {
	gunit.Run(new(IndexingFixture), t)
}

type IndexingFixture struct {
	*gunit.Fixture
}

func (this *IndexingFixture) TestByDate() {
	unordered := []contracts.Page{
		{FrontMatter: contracts.FrontMatter{Title: "2", Date: nu.UTCDate(2000, 1, 2)}},
		{FrontMatter: contracts.FrontMatter{Title: "1", Date: nu.UTCDate(2000, 1, 1)}},
	}
	ordered := OrderByDate(unordered)
	this.So(titles(ordered), should.Resemble, []string{"1", "2"})
}

func titles(ordered []contracts.Page) (titles []string) {
	for _, page := range ordered {
		titles = append(titles, page.Title)
	}
	return titles
}
