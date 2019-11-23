package content

import (
	"sort"
	"testing"

	"github.com/mdwhatcott/huguinho/contracts"
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

func titles(ordered []contracts.Article) (titles []string) {
	for _, page := range ordered {
		titles = append(titles, page.Title)
	}
	return titles
}
func allKeys(site contracts.Site) (keys []string) {
	for key := range site {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (this *IndexingFixture) TestOrganizePages2() {
	all := []contracts.Article{
		{FrontMatter: contracts.FrontMatter{Title: "1", Date: nu.UTCDate(2000, 1, 1), Tags: []string{"a"}}},
		{FrontMatter: contracts.FrontMatter{Title: "2", Date: nu.UTCDate(2000, 1, 2), Tags: []string{"b"}}},
		{FrontMatter: contracts.FrontMatter{Title: "3", Date: nu.UTCDate(2000, 1, 3), Tags: []string{"a", "c"}}},
		{FrontMatter: contracts.FrontMatter{Title: "4", Date: nu.UTCDate(2000, 1, 4), Tags: []string{"b"}}},
		{FrontMatter: contracts.FrontMatter{Title: "5", Date: nu.UTCDate(2000, 1, 5), Tags: []string{"a"}}},
	}

	site := organizeContent(all)

	this.So(allKeys(site), should.Resemble, []string{contracts.HomePageListingID, "a", "b", "c"})
	this.So(titles(site[contracts.HomePageListingID]), should.Resemble, []string{"5", "4", "3", "2", "1"})
	this.So(titles(site["a"]), should.Resemble, []string{"5", "3", "1"})
	this.So(titles(site["b"]), should.Resemble, []string{"4", "2"})
	this.So(titles(site["c"]), should.Resemble, []string{"3"})
}
