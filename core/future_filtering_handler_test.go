package core

import (
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/logging"

	"github.com/mdwhatcott/huguinho/contracts"
)

func TestFutureFilteringHandlerFixture(t *testing.T) {
	gunit.Run(new(FutureFilteringHandlerFixture), t)
}

type FutureFilteringHandlerFixture struct {
	*gunit.Fixture

	now    time.Time
	before time.Time
	after  time.Time
}

func (this *FutureFilteringHandlerFixture) Setup() {
	this.now = time.Now()
	this.before = this.now.Add(-time.Second)
	this.after = this.now.Add(time.Second)
}

func (this *FutureFilteringHandlerFixture) article(date time.Time) *contracts.Article {
	return &contracts.Article{Metadata: contracts.ArticleMetadata{Date: date}}
}

func (this *FutureFilteringHandlerFixture) buildHandler(enabled bool) *FutureFilteringHandler {
	handler := NewFutureFilteringHandler(this.now, enabled)
	handler.log = logging.Capture()
	return handler
}

func (this *FutureFilteringHandlerFixture) TestDisabled_LetEverythingThrough() {
	handler := this.buildHandler(false)

	this.So(handler.Handle(this.article(this.before)), should.BeNil)
	this.So(handler.Handle(this.article(this.now)), should.BeNil)
	this.So(handler.Handle(this.article(this.after)), should.BeNil)
}

func (this *FutureFilteringHandlerFixture) TestEnabled_AnythingAfterNowDropped() {
	handler := this.buildHandler(true)

	this.So(handler.Handle(this.article(this.before)), should.BeNil)
	this.So(handler.Handle(this.article(this.now)), should.BeNil)
	this.So(handler.Handle(this.article(this.after)), should.Resemble, ErrDropArticle)
}
