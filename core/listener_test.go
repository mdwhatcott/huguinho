package core

import (
	"fmt"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/huguinho/contracts"
)

func TestListenerFixture(t *testing.T) {
	gunit.Run(new(ListenerFixture), t)
}

type ListenerFixture struct {
	*gunit.Fixture
	input   chan contracts.Article
	output  chan contracts.Article
	handler *FakeHandler
}

func (this *ListenerFixture) Setup() {
	this.input = make(chan contracts.Article, 10)
	this.output = make(chan contracts.Article, 10)
	this.handler = NewFakeHandler()
}

func (this *ListenerFixture) TestEachArticleHandledIfNotErrantAndPassedOn() {
	this.input <- contracts.Article{Content: contracts.ArticleContent{Original: "A"}}
	this.input <- contracts.Article{Content: contracts.ArticleContent{Original: "B"}, Error: contracts.ErrDropArticle}
	this.input <- contracts.Article{Content: contracts.ArticleContent{Original: "C"}}
	close(this.input)

	Listen(this.input, this.output, this.handler)

	this.So(gather(this.output), should.Resemble, []contracts.Article{
		{Content: contracts.ArticleContent{Original: "A", Converted: "A1"}},
		{Content: contracts.ArticleContent{Original: "B", Converted: ""}, Error: contracts.ErrDropArticle},
		{Content: contracts.ArticleContent{Original: "C", Converted: "C2"}},
	})
}

///////////////////////////////////////////////////////////////

type FakeHandler struct {
	calls int
}

func NewFakeHandler() *FakeHandler {
	return &FakeHandler{}
}

func (this *FakeHandler) Handle(article *contracts.Article) {
	this.calls++
	article.Content.Converted = article.Content.Original + fmt.Sprint(this.calls)
}
