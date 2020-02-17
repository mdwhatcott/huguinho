package core

import (
	"time"

	"github.com/smartystreets/logging"

	"github.com/mdwhatcott/huguinho/contracts"
)

type FutureFilteringHandler struct {
	log     *logging.Logger
	now     time.Time
	enabled bool
}

func NewFutureFilteringHandler(now time.Time, enabled bool) *FutureFilteringHandler {
	return &FutureFilteringHandler{now: now, enabled: enabled}
}

func (this *FutureFilteringHandler) Handle(article *contracts.Article) error {
	if this.enabled && article.Metadata.Date.After(this.now) {
		this.log.Println("[INFO] dropping future article:", article.Metadata.Slug)
		return ErrDropArticle
	}
	return nil
}
