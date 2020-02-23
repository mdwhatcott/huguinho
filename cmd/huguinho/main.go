package main

import (
	"log"
	"os"
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
	renderer := BuildTemplateRenderer(config, NewDisk())
	pipeline := NewPipeline(config, NewDisk(), renderer)
	reporter.ProcessStream(pipeline.Run())
	reporter.RenderFinalReport(time.Now())
	return reporter.Errors()
}

func BuildTemplateRenderer(config contracts.Config, disk contracts.FileSystem) *core.TemplateRenderer {
	loader := core.NewTemplateLoader(disk, config.TemplateDir)
	templates, err := loader.Load()
	if err != nil {
		log.Fatal(err)
	}

	renderer := core.NewTemplateRenderer(templates)
	err = renderer.Validate()
	if err != nil {
		log.Fatal(err)
	}

	return renderer
}
