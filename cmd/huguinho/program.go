package main

import (
	"log"
	"path/filepath"
	"time"

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
	articles, errs := pipeline.Run()
	this.report(articles, errs)
	return errs
}

func (this *Program) buildPipeline() *Pipeline {
	glob := filepath.Join(this.config.TemplateDir, "*.tmpl")
	renderer := shell.ParseTemplates(glob)
	return NewPipeline(this.config, this.disk, renderer)
}

func (this *Program) report(articles int, errs int) {
	log.Printf(
		"[INFO] published %d articles, encountered %d errors, all in %v.",
		articles, errs, time.Since(this.start).Round(time.Millisecond),
	)
}
