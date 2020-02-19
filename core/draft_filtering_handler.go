package core

import (
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
	if this.enabled && article.Metadata.Draft {
		this.log.Println("[INFO] dropping draft article:", article.Metadata.Slug)
		article.Error = contracts.ErrDropArticle
	}
}
