package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/mdwhatcott/huguinho/contracts"
	"github.com/mdwhatcott/huguinho/core"
	"github.com/mdwhatcott/huguinho/shell"
)

func main() {
	log.SetFlags(0)
	reporter := NewReporter(time.Now())
	reporter.ProcessStream(stream())
	reporter.RenderFinalReport(time.Now())
	os.Exit(reporter.errors)
}

func stream() chan contracts.Article {
	config := parseConfig()
	template := shell.ParseTemplates(filepath.Join(config.TemplateDir, "*.tmpl"))
	pipeline := NewPipeline(config, shell.NewDisk(), core.NewTemplateRenderer(template))
	return pipeline.Run()
}
