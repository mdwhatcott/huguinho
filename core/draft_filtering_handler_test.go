package core

import (
	"testing"

	"github.com/mdwhatcott/huguinho/contracts"
	"github.com/mdwhatcott/testing/should"
	"github.com/mdwhatcott/testing/suite"
)

func TestDraftFilteringHandlerFixture(t *testing.T) {
	suite.Run(&DraftFilteringHandlerFixture{T: suite.New(t)}, suite.Options.UnitTests())
}

type DraftFilteringHandlerFixture struct {
	*suite.T
}

func (this *DraftFilteringHandlerFixture) article(draft bool) *contracts.Article {
	return &contracts.Article{Metadata: contracts.ArticleMetadata{Draft: draft}}
}

func (this *DraftFilteringHandlerFixture) TestDisabled_LetEverythingThrough() {
	handler := NewDraftFilteringHandler(false)

	draft := this.article(true)
	handler.Handle(draft)
	this.So(draft.Error, should.BeNil)

	nonDraft := this.article(false)
	handler.Handle(nonDraft)
	this.So(nonDraft.Error, should.BeNil)
}

func (this *DraftFilteringHandlerFixture) TestEnabled_AnyDraftsDropped() {
	handler := NewDraftFilteringHandler(true)

	nonDraft := this.article(false)
	handler.Handle(nonDraft)
	this.So(nonDraft.Error, should.BeNil)

	draft := this.article(true)
	handler.Handle(draft)
	this.So(draft.Error, should.WrapError, contracts.ErrDroppedArticle)
}
