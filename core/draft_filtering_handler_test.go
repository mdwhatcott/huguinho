package core

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/logging"

	"github.com/mdwhatcott/huguinho/contracts"
)

func TestDraftFilteringHandlerFixture(t *testing.T) {
	gunit.Run(new(DraftFilteringHandlerFixture), t)
}

type DraftFilteringHandlerFixture struct {
	*gunit.Fixture
}

func (this *DraftFilteringHandlerFixture) article(draft bool) *contracts.Article {
	return &contracts.Article{Metadata: contracts.ArticleMetadata{Draft: draft}}
}

func (this *DraftFilteringHandlerFixture) buildHandler(enabled bool) *DraftFilteringHandler {
	handler := NewDraftFilteringHandler(enabled)
	handler.log = logging.Capture()
	return handler
}

func (this *DraftFilteringHandlerFixture) TestDisabled_LetEverythingThrough() {
	handler := this.buildHandler(false)

	this.So(handler.Handle(this.article(true)), should.BeNil)
	this.So(handler.Handle(this.article(false)), should.BeNil)
}

func (this *DraftFilteringHandlerFixture) TestEnabled_AnyDraftsDropped() {
	handler := this.buildHandler(true)

	this.So(handler.Handle(this.article(false)), should.BeNil)
	this.So(handler.Handle(this.article(true)), should.Resemble, ErrDropArticle)
}
