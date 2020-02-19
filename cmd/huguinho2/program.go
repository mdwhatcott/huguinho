package main

import (
	"log"
	"path/filepath"
	"time"

	"github.com/mdwhatcott/huguinho/contracts"
	"github.com/mdwhatcott/huguinho/shell"
)

type Program struct {
	disk   *shell.Disk
	config contracts.Config
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
	this.exportCSS()
	pipeline := this.buildPipeline()
	articles, errs := pipeline.Run()
	this.report(articles, errs)
	return errs
}

func (this *Program) exportCSS() {
	err := this.disk.CopyFile(
		filepath.Join(this.config.StylesDir, "custom.css"),
		filepath.Join(this.config.TargetRoot, "css", "custom.css"),
	)
	if err != nil {
		log.Fatal("[WARN] failed to copy css:", err)
	}
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
