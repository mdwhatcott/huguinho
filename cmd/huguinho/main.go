package main

import (
	"log"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/mdwhatcott/huguinho/contracts"
	"github.com/mdwhatcott/huguinho/core"
)

func main() {
	log.SetFlags(0)
	reporter := core.NewReporter(time.Now())
	reporter.ProcessStream(stream())
	reporter.RenderFinalReport(time.Now())
	os.Exit(reporter.Errors())
}

func stream() chan contracts.Article {
	config := parseConfig()
	templates := parseTemplates(filepath.Join(config.TemplateDir, "*.tmpl"))
	pipeline := NewPipeline(config, NewDisk(), templates)
	return pipeline.Run()
}

func parseTemplates(glob string) *template.Template {
	templates, err := template.ParseGlob(glob)
	if err != nil {
		log.Fatalln("Could not parse templates:", err)
	}
	return templates
}
