package core

import (
	"errors"
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/huguinho/contracts"
)

func TestFutureFilteringHandlerFixture(t *testing.T) {
	gunit.Run(new(FutureFilteringHandlerFixture), t)
}

type FutureFilteringHandlerFixture struct {
	*gunit.Fixture

	present time.Time
	past    time.Time
	future  time.Time
}

func (this *FutureFilteringHandlerFixture) Setup() {
	this.present = time.Now()
	this.past = this.present.Add(-time.Second)
	this.future = this.present.Add(time.Second)
}

func (this *FutureFilteringHandlerFixture) article(date time.Time) *contracts.Article {
	return &contracts.Article{Metadata: contracts.ArticleMetadata{Date: date}}
}

func (this *FutureFilteringHandlerFixture) TestDisabled_LetEverythingThrough() {
	disabled := NewFutureFilteringHandler(this.present, false)

	past := this.article(this.past)
	disabled.Handle(past)
	this.So(past.Error, should.BeNil)

	present := this.article(this.present)
	disabled.Handle(present)
	this.So(present.Error, should.BeNil)

	future := this.article(this.future)
	disabled.Handle(future)
	this.So(future.Error, should.BeNil)
}

func (this *FutureFilteringHandlerFixture) TestEnabled_AnythingAfterNowDropped() {
	enabled := NewFutureFilteringHandler(this.present, true)

	past := this.article(this.past)
	enabled.Handle(past)
	this.So(past.Error, should.BeNil)

	present := this.article(this.present)
	enabled.Handle(present)
	this.So(present.Error, should.BeNil)

	future := this.article(this.future)
	enabled.Handle(future)
	this.So(errors.Is(future.Error, contracts.ErrDropArticle), should.BeTrue)
}
