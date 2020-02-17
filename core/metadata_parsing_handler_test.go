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
	err := this.parser.Handle(this.article)
	this.So(errors.Is(err, errMissingMetadata), should.BeTrue)
}
func (this *MetadataParserFixture) TestMissingMetadataDivider_Err() {
	this.appendSourceLine("title: This is the title")
	this.appendSourceLine("")
	this.appendSourceLine("This is the article content")

	err := this.parser.Handle(this.article)

	this.So(errors.Is(err, errMissingMetadataDivider), should.BeTrue)
}
func (this *MetadataParserFixture) TestDuplicateTitle_Err() {
	this.appendMetadataWithContent(
		"title: This is the title",
		"title: This is another title",
	)

	err := this.parser.Handle(this.article)

	this.So(errors.Is(err, errDuplicateMetadataTitle), should.BeTrue)
}
func (this *MetadataParserFixture) TestBlankTitle_Err() {
	this.appendMetadataWithContent("title: ")

	err := this.parser.Handle(this.article)

	this.So(errors.Is(err, errBlankMetadataTitle), should.BeTrue)
}
func (this *MetadataParserFixture) TestDuplicateIntro_Err() {
	this.appendMetadataWithContent(
		"intro: This is the intro",
		"intro: This is another intro",
	)

	err := this.parser.Handle(this.article)

	this.So(errors.Is(err, errDuplicateMetadataIntro), should.BeTrue)
}
func (this *MetadataParserFixture) TestBlankIntro_Err() {
	this.appendMetadataWithContent("intro: ")

	err := this.parser.Handle(this.article)

	this.So(errors.Is(err, errBlankMetadataIntro), should.BeTrue)
}
func (this *MetadataParserFixture) TestDuplicateSlug_Err() {
	this.appendMetadataWithContent(
		"slug: /this/is/the/slug",
		"slug: /this/is/another/slug",
	)

	err := this.parser.Handle(this.article)

	this.So(errors.Is(err, errDuplicateMetadataSlug), should.BeTrue)
}
func (this *MetadataParserFixture) TestInvalidSlug_Err() {
	this.appendMetadataWithContent("slug: /this/slug/is invalid")

	err := this.parser.Handle(this.article)

	this.So(errors.Is(err, errInvalidMetadataSlug), should.BeTrue)
}
func (this *MetadataParserFixture) TestInvalidCasingSlug_Err() {
	this.appendMetadataWithContent("slug: /THIS/SLUG/IS/INVALID")

	err := this.parser.Handle(this.article)

	this.So(errors.Is(err, errInvalidMetadataSlug), should.BeTrue)
}
func (this *MetadataParserFixture) TestBlankSlug_Err() {
	this.appendMetadataWithContent("slug: ")

	err := this.parser.Handle(this.article)

	this.So(errors.Is(err, errBlankMetadataSlug), should.BeTrue)
}
func (this *MetadataParserFixture) TestBlankDraft_Err() {
	this.appendMetadataWithContent("draft: ")

	err := this.parser.Handle(this.article)

	this.So(errors.Is(err, errBlankMetadataDraft), should.BeTrue)
}
func (this *MetadataParserFixture) TestInvalidDraft_Err() {
	this.appendMetadataWithContent("draft: not-true-or-false")

	err := this.parser.Handle(this.article)

	this.So(errors.Is(err, errInvalidMetadataDraft), should.BeTrue)
}
func (this *MetadataParserFixture) TestMetadataDraftFalse_Valid() {
	this.appendMetadataWithContent("draft: false")

	err := this.parser.Handle(this.article)

	this.So(errors.Is(err, errBlankMetadataDraft), should.BeFalse)
	this.So(errors.Is(err, errInvalidMetadataDraft), should.BeFalse)
}
func (this *MetadataParserFixture) TestDuplicateDraft_Err() {
	this.appendMetadataWithContent(
		"draft: false",
		"draft: true",
	)

	err := this.parser.Handle(this.article)

	this.So(errors.Is(err, errDuplicateMetadataDraft), should.BeTrue)
}
func (this *MetadataParserFixture) TestBlankDate_Err() {
	this.appendMetadataWithContent("date: ")

	err := this.parser.Handle(this.article)

	this.So(errors.Is(err, errBlankMetadataDate), should.BeTrue)
}
func (this *MetadataParserFixture) TestInvalidDate_Err() {
	this.appendMetadataWithContent("date: not-a-date")

	err := this.parser.Handle(this.article)

	this.So(errors.Is(err, errInvalidMetadataDate), should.BeTrue)
}
func (this *MetadataParserFixture) TestDuplicateDate_Err() {
	this.appendMetadataWithContent(
		"date: 2020-02-01",
		"date: 2020-02-02",
	)

	err := this.parser.Handle(this.article)

	this.So(errors.Is(err, errDuplicateMetadataDate), should.BeTrue)
}
func (this *MetadataParserFixture) TestBlankTags_Err() {
	this.appendMetadataWithContent("tags: ")

	err := this.parser.Handle(this.article)

	this.So(errors.Is(err, errBlankMetadataTags), should.BeTrue)
}
func (this *MetadataParserFixture) TestInvalidTags_Err() {
	this.appendMetadataWithContent("tags: invalid?!")

	err := this.parser.Handle(this.article)

	this.So(errors.Is(err, errInvalidMetadataTags), should.BeTrue)
}
func (this *MetadataParserFixture) TestInvalidCasingTags_Err() {
	this.appendMetadataWithContent("tags: INVALID")

	err := this.parser.Handle(this.article)

	this.So(errors.Is(err, errInvalidMetadataTags), should.BeTrue)
}
func (this *MetadataParserFixture) TestDuplicateTags_Err() {
	this.appendMetadataWithContent(
		"tags: a b c",
		"tags: x y z",
	)

	err := this.parser.Handle(this.article)

	this.So(errors.Is(err, errDuplicateMetadataTags), should.BeTrue)
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

	err := this.parser.Handle(this.article)

	this.So(err, should.BeNil)
	this.So(this.article.Metadata.Title, should.Equal, "This is the title")
	this.So(this.article.Metadata.Intro, should.Equal, "This is the intro")
	this.So(this.article.Metadata.Slug, should.Equal, "/this/is/the/slug")
	this.So(this.article.Metadata.Draft, should.BeTrue)
	this.So(this.article.Metadata.Date, should.Resemble, Date(2020, 2, 16))
	this.So(this.article.Metadata.Tags, should.Resemble, []string{"a-a", "b", "c"})
}
