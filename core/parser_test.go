package core

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/huguinho/contracts"
)

func TestPageParserFixture(t *testing.T) {
	gunit.Run(new(PageParserFixture), t)
}

type PageParserFixture struct {
	*gunit.Fixture

	parser    *PageParser
	converter *FakeContentConverter
	sample    contracts.Page
}

func date(y, m, d int) time.Time {
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
}

func (this *PageParserFixture) prepareValidPage() string {
	front, err := json.Marshal(this.sample.Metadata)
	this.So(err, should.BeNil)
	return string(front) + "\n\n+++\n\n" + this.sample.Content.Original
}
func (this *PageParserFixture) preparePage_MissingDivider() string {
	front, err := json.Marshal(this.sample.Metadata)
	this.So(err, should.BeNil)
	return string(front) + "\n\n" + /* missing divider */ "\n\n" + this.sample.Content.Original
}
func (this *PageParserFixture) preparePage_MalformedFrontMatter() string {
	return "malformed front matter" + "\n\n+++\n\n" + this.sample.Content.Original
}

func (this *PageParserFixture) assertMetadataDecoded(page contracts.Page) bool {
	return this.So(page.Metadata, should.Resemble, this.sample.Metadata)
}
func (this *PageParserFixture) assertContentConverted(page contracts.Page) {
	this.So(this.converter.original, should.Equal, this.sample.Content.Original)
	this.So(page.Content.Original, should.Equal, this.sample.Content.Original)
	this.So(page.Content.Converted, should.Equal, CONVERTED_CONTENT)
}

func (this *PageParserFixture) Setup() {
	this.converter = NewFakeContentConverter()
	this.parser = NewPageParser(this.converter)
	this.sample = contracts.Page{
		Metadata: contracts.JSONFrontMatter{
			Slug:        "/slug",
			Title:       "title",
			Description: "description",
			Date:        date(2020, 2, 1),
			Tags:        []string{"a", "b"},
			IsDraft:     true,
		},
		Content: contracts.Content{
			Original: "# H1",
		},
	}
}

func (this *PageParserFixture) TestValidPageParsed() {
	rawPage := this.prepareValidPage()

	page, err := this.parser.ParsePage(rawPage)

	this.So(err, should.BeNil)
	this.assertMetadataDecoded(page)
	this.assertContentConverted(page)
}
func (this *PageParserFixture) TestMissingFrontMatterDivider() {
	rawPage := this.preparePage_MissingDivider()

	page, err := this.parser.ParsePage(rawPage)

	this.So(err, should.NotBeNil)
	this.So(page, should.BeZeroValue)
	this.So(this.converter.original, should.BeZeroValue)
}
func (this *PageParserFixture) SkipTestMalformedFrontMatter()        {}
func (this *PageParserFixture) SkipTestInvalidSlug_URLUnsafe()       {}
func (this *PageParserFixture) SkipTestInvalidSlug_NoLeadingSlash()  {}
func (this *PageParserFixture) SkipTestInvalidSlug_NoTrailingSlash() {}
func (this *PageParserFixture) SkipTestMissingTitle()                {}
func (this *PageParserFixture) SkipTestMissingSlug()                 {}
func (this *PageParserFixture) SkipTestMissingDate()                 {}

////////////////////////////////////////////////////////

type FakeContentConverter struct {
	original string
	err      error
}

func NewFakeContentConverter() *FakeContentConverter {
	return &FakeContentConverter{}
}

func (this *FakeContentConverter) Convert(content string) (string, error) {
	this.original = content
	return CONVERTED_CONTENT, this.err
}

const CONVERTED_CONTENT = "CONVERTED"
