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
	config := parseConfig()
	reporter := core.NewReporter(time.Now())
	reporter.ProcessStream(stream(config))
	reporter.RenderFinalReport(time.Now())
	os.Exit(reporter.Errors())
}

func stream(config Config) chan contracts.Article {
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
