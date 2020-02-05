package core

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/huguinho/contracts"
)

func TestJSONMetadataParserFixture(t *testing.T) {
	gunit.Run(new(JSONMetadataParserFixture), t)
}

type JSONMetadataParserFixture struct {
	*gunit.Fixture

	parser *JSONMetadataParser
}

func (this *JSONMetadataParserFixture) Setup() {
	this.parser = NewJSONMetadataParser()
}

func (this *JSONMetadataParserFixture) prepareArticleForParsing(content string) *contracts.Article {
	return &contracts.Article{
		Source: contracts.ArticleSource{
			Data: content,
		},
	}
}

func (this *JSONMetadataParserFixture) prepareValidMetadataJSON() (JSONMetadata, []byte) {
	valid := JSONMetadata{
		Draft: true,
		Slug:  "slug",
		Title: "title",
		Intro: "intro",
		Tags:  []string{"tag1", "tag2"},
		Date:  time.Date(2020, 2, 4, 0, 0, 0, 0, time.UTC),
	}
	raw, err := json.Marshal(valid)
	this.So(err, should.BeNil)
	return valid, raw
}

func (this *JSONMetadataParserFixture) TestValidMetadata() {
	valid, raw := this.prepareValidMetadataJSON()
	article := this.prepareArticleForParsing(string(raw) + contracts.METADATA_CONTENT_DIVIDER + "Content")

	err := this.parser.Handle(article)

	this.So(err, should.BeNil)
	this.So(article.Metadata, should.Resemble, contracts.ArticleMetadata{
		Draft: valid.Draft,
		Slug:  valid.Slug,
		Title: valid.Title,
		Intro: valid.Intro,
		Tags:  valid.Tags,
		Date:  valid.Date,
	})
}
