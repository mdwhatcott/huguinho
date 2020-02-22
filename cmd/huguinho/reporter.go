package main

import (
	"errors"
	"time"

	"github.com/smartystreets/logging"

	"github.com/mdwhatcott/huguinho/contracts"
)

type Reporter struct {
	log       *logging.Logger
	started   time.Time
	errors    int
	dropped   int
	published int
}

func NewReporter(started time.Time) *Reporter {
	return &Reporter{started: started}
}

func (this *Reporter) ProcessStream(out chan contracts.Article) {
	this.log.SetFlags(0)
	for article := range out {
		this.accountFor(article)
	}
}

func (this *Reporter) accountFor(article contracts.Article) {
	if errors.Is(article.Error, contracts.ErrDroppedArticle) {
		this.log.Println("[INFO]", article.Error)
		this.dropped++
	} else if article.Error != nil {
		this.log.Println("[WARN] error:", article.Error)
		this.errors++
	} else if article.Source.Path != "" {
		this.log.Println("[INFO] published article:", article.Metadata.Slug)
		this.published++
	} else {
		this.log.Printf("[WARN] not sure what this article struct represents: %#v", article)
		this.errors++
	}
}

func (this Reporter) RenderFinalReport(finished time.Time) {
	this.log.Println("[INFO] errors encountered: ", this.errors)
	this.log.Println("[INFO] articles dropped:   ", this.dropped)
	this.log.Println("[INFO] articles published: ", this.published)
	this.log.Println("[INFO] processing duration:", finished.Sub(this.started).Round(time.Millisecond))
}
