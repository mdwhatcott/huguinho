package main

import (
	"errors"
	"log"
	"path/filepath"
	"time"

	"github.com/mdwhatcott/huguinho/contracts"
	"github.com/mdwhatcott/huguinho/shell"
)

// TEST (integration)
type Program struct {
	disk   *shell.Disk
	config Config
	start  time.Time
}

func NewProgram(start time.Time) *Program {
	log.SetFlags(0)
	return &Program{
		disk:   shell.NewDisk(),
		config: parseConfig(),
		start:  start,
	}
}

func (this *Program) Run() int {
	pipeline := this.buildPipeline()
	return this.receive(pipeline.Run())
}

func (this *Program) buildPipeline() *Pipeline {
	glob := filepath.Join(this.config.TemplateDir, "*.tmpl")
	renderer := shell.ParseTemplates(glob)
	return NewPipeline(this.config, this.disk, renderer)
}

func (this *Program) receive(out chan contracts.Article) int {
	var (
		errs      int
		published int
		dropped   int
	)

	for article := range out {
		if errors.Is(article.Error, contracts.ErrDropArticle) {
			log.Println("[INFO]", article.Error)
			dropped++
		} else if article.Error != nil {
			log.Println("[WARN] error:", article.Error)
			errs++
		} else if article.Source.Path != "" {
			log.Println("[INFO] published article:", article.Metadata.Slug)
			published++
		} else {
			log.Printf("[WARN] not sure what this article struct represents: %#v", article)
			errs++
		}
	}
	this.report(dropped, published, errs)
	return errs
}
func (this *Program) report(dropped, published, errs int) {
	log.Println("[INFO] errors encountered: ", errs)
	log.Println("[INFO] articles dropped:   ", dropped)
	log.Println("[INFO] articles published: ", published)
	log.Println("[INFO] processing duration:", time.Since(this.start).Round(time.Millisecond))
}
