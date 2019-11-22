package site

import (
	"testing"

	"github.com/mdwhatcott/static/contracts"
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestPageFixture(t *testing.T) {
	gunit.Run(new(PageFixture), t)
}

type PageFixture struct {
	*gunit.Fixture
}

func (this *PageFixture) TestConvertEmptyFileToPage() {
	file := contracts.File("")
	page := ConvertToPage(file)
	this.So(page, should.Resemble, contracts.Page{})
}

func (this *PageFixture) TestConvertContentOnlyFileToPage() {
	file := contracts.File("I have some content")
	page := ConvertToPage(file)
	this.So(page, should.Resemble, contracts.Page{
		OriginalContent: "I have some content",
		HTMLContent:     "<p>I have some content</p>\n",
	})
}

func (this *PageFixture) TestConvertEmptyFrontMatterAndContentToPage() {
	file := contracts.File("+++\n\n+++\nI have some content")
	page := ConvertToPage(file)
	this.So(page, should.Resemble, contracts.Page{
		OriginalContent: "I have some content",
		HTMLContent:     "<p>I have some content</p>\n",
	})
}
