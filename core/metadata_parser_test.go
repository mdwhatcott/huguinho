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

func (this *MetadataParserFixture) TestBlank_Err() {
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

func (this *MetadataParserFixture) TestRepeatedMetadataTitle_Err() {
	this.appendMetadataWithContent(
		"title: This is the title",
		"title: This is another title",
	)

	err := this.parser.Handle(this.article)

	this.So(errors.Is(err, errDuplicateMetadataTitle), should.BeTrue)
}

func (this *MetadataParserFixture) TestRepeatedMetadataIntro_Err() {
	this.appendMetadataWithContent(
		"intro: This is the intro",
		"intro: This is another intro",
	)

	err := this.parser.Handle(this.article)

	this.So(errors.Is(err, errDuplicateMetadataIntro), should.BeTrue)
}

func (this *MetadataParserFixture) TestRepeatedMetadataSlug_Err() {
	this.appendMetadataWithContent(
		"slug: /this/is/the/slug",
		"slug: /this/is/another/slug",
	)

	err := this.parser.Handle(this.article)

	this.So(errors.Is(err, errDuplicateMetadataSlug), should.BeTrue)
}
func (this *MetadataParserFixture) TestInvalidMetadataSlug_Err() {
	this.appendMetadataWithContent("slug: /this/slug/is invalid")

	err := this.parser.Handle(this.article)

	this.So(errors.Is(err, errInvalidMetadataSlug), should.BeTrue)
}
func (this *MetadataParserFixture) TestBlankMetadataSlug_Err() {
	this.appendMetadataWithContent("slug: ")

	err := this.parser.Handle(this.article)

	this.So(errors.Is(err, errBlankMetadataSlug), should.BeTrue)
}

func (this *MetadataParserFixture) TestValidMetadata() {
	this.appendMetadataWithContent( // TODO: test showing that trailing/leading whitespace is trimmed
		"title: This is the title",
		"intro: This is the intro",
		"slug:  /this/is/the/slug",
		// TODO: draft
		// TODO: date
		// TODO: tags
	)

	err := this.parser.Handle(this.article)

	this.So(err, should.BeNil)
	this.So(this.article.Metadata.Title, should.Equal, "This is the title")
	this.So(this.article.Metadata.Intro, should.Equal, "This is the intro")
	this.So(this.article.Metadata.Slug, should.Equal, "/this/is/the/slug")
}
