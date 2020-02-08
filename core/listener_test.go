package core

import (
	"errors"
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
	input    chan contracts.Article
	output   chan contracts.Article
	handler  *FakeHandler
	listener *Listener
}

func (this *ListenerFixture) Setup() {
	this.input = make(chan contracts.Article, 10)
	this.output = make(chan contracts.Article, 10)
	this.handler = NewFakeHandler()
	this.listener = NewListener(this.input, this.output, this.handler)

	this.input <- contracts.Article{Content: contracts.ArticleContent{Original: "A"}}
	this.input <- contracts.Article{Content: contracts.ArticleContent{Original: "B"}}
	this.input <- contracts.Article{Content: contracts.ArticleContent{Original: "C"}}
	close(this.input)
}

func (this *ListenerFixture) setupWithFakeFinalizingHandler() {
	handler := NewFakeFinalizingHandler()
	handler.final = errors.New("final")
	this.listener = NewListener(this.input, this.output, handler)
}

func (this *ListenerFixture) TestEachArticleHandledAndPassedOn() {
	err := this.listener.Listen()

	this.So(err, should.BeNil)
	this.So(gather(this.output), should.Resemble, []contracts.Article{
		{Content: contracts.ArticleContent{Original: "A", Converted: "A1"}},
		{Content: contracts.ArticleContent{Original: "B", Converted: "B2"}},
		{Content: contracts.ArticleContent{Original: "C", Converted: "C3"}},
	})
}

func (this *ListenerFixture) TestSomeArticlesMightBeDropped() {
	this.handler.errs[2] = ErrDropArticle

	err := this.listener.Listen()

	this.So(err, should.BeNil)
	this.So(gather(this.output), should.Resemble, []contracts.Article{
		{Content: contracts.ArticleContent{Original: "A", Converted: "A1"}},
		{Content: contracts.ArticleContent{Original: "C", Converted: "C3"}},
	})
}

func (this *ListenerFixture) TestFirstHandlerToErrEndsHandling() {
	handlerError := errors.New("handler error")
	this.handler.errs[2] = handlerError

	err := this.listener.Listen()

	this.So(errors.Is(err, handlerError), should.BeTrue)
	this.So(gather(this.output), should.Resemble, []contracts.Article{
		{Content: contracts.ArticleContent{Original: "A", Converted: "A1"}},
	})
}

func (this *ListenerFixture) TestFinalizingHandlersAreTreatedAsSuch() {
	this.setupWithFakeFinalizingHandler()

	err := this.listener.Listen()

	this.So(err, should.Resemble, errors.New("final"))
	this.So(gather(this.output), should.HaveLength, 3)
}

///////////////////////////////////////////////////////////////

type FakeHandler struct {
	articles []contracts.Article
	calls    int
	errs     map[int]error
}

func NewFakeHandler() *FakeHandler {
	return &FakeHandler{errs: make(map[int]error)}
}

func (this *FakeHandler) Handle(article *contracts.Article) error {
	this.calls++
	article.Content.Converted = article.Content.Original + fmt.Sprint(this.calls)
	return this.errs[this.calls]
}

///////////////////////////////////////////////////////////////

type FakeFinalizingHandler struct {
	*FakeHandler
	final error
}

func NewFakeFinalizingHandler() *FakeFinalizingHandler {
	return &FakeFinalizingHandler{FakeHandler: NewFakeHandler()}
}

func (this *FakeFinalizingHandler) Finalize() error {
	return this.final
}

///////////////////////////////////////////////////////////////
