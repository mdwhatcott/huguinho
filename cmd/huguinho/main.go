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

	os.Exit(Run(ParseConfig(os.Args[1:])))
}

func ParseConfig(args []string) contracts.Config {
	config, err := contracts.NewCLIParser(args).Parse()
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func Run(config contracts.Config) int {
	reporter := core.NewReporter(time.Now())
	reporter.ProcessStream(stream(config))
	reporter.RenderFinalReport(time.Now())
	return reporter.Errors()
}

func stream(config contracts.Config) chan contracts.Article {
	templates := parseTemplates(filepath.Join(config.TemplateDir, "*.tmpl"))
	pipeline := NewPipeline(config, NewDisk(), templates)
	return pipeline.Run()
}

// TEST (see core.TemplateLoader)
func parseTemplates(glob string) *template.Template {
	templates, err := template.ParseGlob(glob)
	if err != nil {
		log.Fatalln("Could not parse templates:", err)
	}
	//templates.Lookup(contracts.HomePageTemplateName)
	return templates
}
