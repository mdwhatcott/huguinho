package core

import (
	"errors"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/huguinho/contracts"
)

func TestMetadataParserFixture(t *testing.T) {
	gunit.Run(new(MetadataParserFixture), t)
}

type MetadataParserFixture struct {
	*gunit.Fixture
	parser  *MetadataParsingHandler
	article *contracts.Article
}

func (this *MetadataParserFixture) Setup() {
	this.parser = NewMetadataParsingHandler()
	this.article = &contracts.Article{}
}

func (this *MetadataParserFixture) appendSourceLine(line string) {
	this.article.Source.Data += line + "\n"
}
func (this *MetadataParserFixture) appendDividerAndContent() {
	this.appendSourceLine("\n+++\n")
	this.appendSourceLine("article content")
}
func (this *MetadataParserFixture) appendMetadataWithContent(lines ...string) {
	for _, line := range lines {
		this.appendSourceLine(line)
	}
	this.appendDividerAndContent()
}

func (this *MetadataParserFixture) TestBlankFile_Err() {
	this.appendSourceLine(" \t \n  ")
	this.parser.Handle(this.article)
	this.So(errors.Is(this.article.Error, errMissingMetadata), should.BeTrue)
}
func (this *MetadataParserFixture) TestMissingMetadataDivider_Err() {
	this.appendSourceLine("title: This is the title")
	this.appendSourceLine("")
	this.appendSourceLine("This is the article content")

	this.parser.Handle(this.article)

	this.So(errors.Is(this.article.Error, errMissingMetadataDivider), should.BeTrue)
}
func (this *MetadataParserFixture) TestDuplicateTitle_Err() {
	this.appendMetadataWithContent(
		"title: This is the title",
		"title: This is another title",
	)

	this.parser.Handle(this.article)

	this.So(errors.Is(this.article.Error, errDuplicateMetadataTitle), should.BeTrue)
}
func (this *MetadataParserFixture) TestBlankTitle_Err() {
	this.appendMetadataWithContent("title: ")

	this.parser.Handle(this.article)

	this.So(errors.Is(this.article.Error, errBlankMetadataTitle), should.BeTrue)
}
func (this *MetadataParserFixture) TestDuplicateIntro_Err() {
	this.appendMetadataWithContent(
		"intro: This is the intro",
		"intro: This is another intro",
	)

	this.parser.Handle(this.article)

	this.So(errors.Is(this.article.Error, errDuplicateMetadataIntro), should.BeTrue)
}
func (this *MetadataParserFixture) TestDuplicateSlug_Err() {
	this.appendMetadataWithContent(
		"slug: /this/is/the/slug",
		"slug: /this/is/another/slug",
	)

	this.parser.Handle(this.article)

	this.So(errors.Is(this.article.Error, errDuplicateMetadataSlug), should.BeTrue)
}
func (this *MetadataParserFixture) TestInvalidSlug_Err() {
	this.appendMetadataWithContent("slug: /this/slug/is invalid")

	this.parser.Handle(this.article)

	this.So(errors.Is(this.article.Error, errInvalidMetadataSlug), should.BeTrue)
}
func (this *MetadataParserFixture) TestInvalidCasingSlug_Err() {
	this.appendMetadataWithContent("slug: /THIS/SLUG/IS/INVALID")

	this.parser.Handle(this.article)

	this.So(errors.Is(this.article.Error, errInvalidMetadataSlug), should.BeTrue)
}
func (this *MetadataParserFixture) TestBlankSlug_Err() {
	this.appendMetadataWithContent("slug: ")

	this.parser.Handle(this.article)

	this.So(errors.Is(this.article.Error, errBlankMetadataSlug), should.BeTrue)
}
func (this *MetadataParserFixture) TestBlankDraft_Err() {
	this.appendMetadataWithContent("draft: ")

	this.parser.Handle(this.article)

	this.So(errors.Is(this.article.Error, errBlankMetadataDraft), should.BeTrue)
}
func (this *MetadataParserFixture) TestInvalidDraft_Err() {
	this.appendMetadataWithContent("draft: not-true-or-false")

	this.parser.Handle(this.article)

	this.So(errors.Is(this.article.Error, errInvalidMetadataDraft), should.BeTrue)
}
func (this *MetadataParserFixture) TestMetadataDraftFalse_Valid() {
	this.appendMetadataWithContent("draft: false")

	this.parser.Handle(this.article)

	this.So(errors.Is(this.article.Error, errBlankMetadataDraft), should.BeFalse)
	this.So(errors.Is(this.article.Error, errInvalidMetadataDraft), should.BeFalse)
}
func (this *MetadataParserFixture) TestDuplicateDraft_Err() {
	this.appendMetadataWithContent(
		"draft: false",
		"draft: true",
	)

	this.parser.Handle(this.article)

	this.So(errors.Is(this.article.Error, errDuplicateMetadataDraft), should.BeTrue)
}
func (this *MetadataParserFixture) TestBlankDate_Err() {
	this.appendMetadataWithContent("date: ")

	this.parser.Handle(this.article)

	this.So(errors.Is(this.article.Error, errBlankMetadataDate), should.BeTrue)
}
func (this *MetadataParserFixture) TestInvalidDate_Err() {
	this.appendMetadataWithContent("date: not-a-date")

	this.parser.Handle(this.article)

	this.So(errors.Is(this.article.Error, errInvalidMetadataDate), should.BeTrue)
}
func (this *MetadataParserFixture) TestDuplicateDate_Err() {
	this.appendMetadataWithContent(
		"date: 2020-02-01",
		"date: 2020-02-02",
	)

	this.parser.Handle(this.article)

	this.So(errors.Is(this.article.Error, errDuplicateMetadataDate), should.BeTrue)
}
func (this *MetadataParserFixture) TestInvalidTags_Err() {
	this.appendMetadataWithContent("tags: invalid?!")

	this.parser.Handle(this.article)

	this.So(errors.Is(this.article.Error, errInvalidMetadataTags), should.BeTrue)
}
func (this *MetadataParserFixture) TestInvalidCasingTags_Err() {
	this.appendMetadataWithContent("tags: INVALID")

	this.parser.Handle(this.article)

	this.So(errors.Is(this.article.Error, errInvalidMetadataTags), should.BeTrue)
}
func (this *MetadataParserFixture) TestInvalidRepeatedTags_Err() {
	this.appendMetadataWithContent("tags: repeated something-else repeated")

	this.parser.Handle(this.article)

	this.So(errors.Is(this.article.Error, errInvalidMetadataTags), should.BeTrue)
}
func (this *MetadataParserFixture) TestDuplicateTags_Err() {
	this.appendMetadataWithContent(
		"tags: a b c",
		"tags: x y z",
	)

	this.parser.Handle(this.article)

	this.So(errors.Is(this.article.Error, errDuplicateMetadataTags), should.BeTrue)
}

func (this *MetadataParserFixture) TestAllValidAttributes() {
	this.appendMetadataWithContent(
		"title: This is the title ",
		"intro: This is the intro ",
		"slug:  /this/is/the/slug ",
		"draft: true              ",
		"date:  2020-02-16        ",
		"tags:  a-a b c           ",
	)

	this.parser.Handle(this.article)

	this.So(this.article.Error, should.BeNil)
	this.So(this.article.Metadata.Title, should.Equal, "This is the title")
	this.So(this.article.Metadata.Intro, should.Equal, "This is the intro")
	this.So(this.article.Metadata.Slug, should.Equal, "/this/is/the/slug")
	this.So(this.article.Metadata.Draft, should.BeTrue)
	this.So(this.article.Metadata.Date, should.Resemble, Date(2020, 2, 16))
	this.So(this.article.Metadata.Tags, should.Resemble, []string{"a-a", "b", "c"})
}
