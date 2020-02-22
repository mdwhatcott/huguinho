package core

import (
	"fmt"

	"github.com/smartystreets/logging"

	"github.com/mdwhatcott/huguinho/contracts"
)

type DraftFilteringHandler struct {
	log     *logging.Logger
	enabled bool
}

func NewDraftFilteringHandler(enabled bool) *DraftFilteringHandler {
	return &DraftFilteringHandler{enabled: enabled}
}

func (this *DraftFilteringHandler) Handle(article *contracts.Article) {
	if !this.enabled {
		return
	}
	if !article.Metadata.Draft {
		return
	}
	article.Error = fmt.Errorf(
		"%w: %s (DRAFT)",
		contracts.ErrDroppedArticle,
		article.Metadata.Slug,
	)
}
