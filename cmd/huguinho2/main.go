package main

import (
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/mdwhatcott/huguinho/core"
	"github.com/mdwhatcott/huguinho/shell"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config := ParseCLI()
	templates := parseTemplates(filepath.Join(config.TemplateDir, "*.tmpl"))
	pipeline := core.NewPipeline(
		config,
		shell.NewDisk(),
		core.NewTemplateRenderer(templates),
	)
	os.Exit(pipeline.Run())
}
func parseTemplates(glob string) *template.Template {
	templates, err := template.ParseGlob(glob)
	if err != nil {
		log.Fatalln("Could not parse templates:", err)
	}
	return templates
}
