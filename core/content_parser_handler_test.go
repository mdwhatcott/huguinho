package core

import (
	"errors"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/huguinho/contracts"
)

func TestContentParserFixture(t *testing.T) {
	gunit.Run(new(ContentParserFixture), t)
}

type ContentParserFixture struct {
	*gunit.Fixture

	converter *ContentParsingHandler
	inner     *FakeConverter
	input     chan contracts.Article
	output    chan contracts.Article
}

func (this *ContentParserFixture) Setup() {
	this.inner = NewFakeConverter()
	this.converter = NewContentParsingHandler(this.inner)
}

func (this *ContentParserFixture) formatSourceData(original string) string {
	return "\n" + contracts.METADATA_CONTENT_DIVIDER + "\n" + original
}

func (this *ContentParserFixture) TestValidContentParsedAndConverted() {
	article := &contracts.Article{Source: contracts.ArticleSource{Data: this.formatSourceData("content1")}}

	this.converter.Handle(article)

	this.So(article, should.Resemble, &contracts.Article{
		Error:   nil,
		Source:  contracts.ArticleSource{Data: this.formatSourceData("content1")},
		Content: contracts.ArticleContent{Original: "content1", Converted: "content1 (CONVERTED)"},
	})
}

func (this *ContentParserFixture) TestInvalidContentElicitsError() {
	article := &contracts.Article{Source: contracts.ArticleSource{Data: this.formatSourceData("content1")}}
	conversionError := errors.New("conversion error")
	this.inner.err = conversionError

	this.converter.Handle(article)

	this.So(errors.Is(article.Error, conversionError), should.BeTrue)
	this.So(article.Content, should.BeZeroValue)
}

////////////////////////////////////////////////////

type FakeConverter struct {
	original string
	err      error
}

func NewFakeConverter() *FakeConverter {
	return &FakeConverter{}
}

func (this *FakeConverter) Convert(content string) (string, error) {
	this.original = content
	return content + " (CONVERTED)", this.err
}
