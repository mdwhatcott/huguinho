package main

import (
	"flag"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mdwhatcott/huguinho/contracts"
	"github.com/mdwhatcott/huguinho/core"
	"github.com/mdwhatcott/huguinho/shell"
)

func main() {
	start := time.Now()
	log.SetFlags(0)
	log.SetPrefix("(stderr) ")

	disk := shell.NewDisk()
	config := parseConfig()
	exportCSS(disk, config)
	templates := parseTemplates(filepath.Join(config.TemplateDir, "*.tmpl"))
	renderer := core.NewTemplateRenderer(templates)
	pipeline := core.NewPipeline(config, disk, renderer)
	errs := pipeline.Run()
	log.Println("[INFO] duration:", time.Since(start))
	os.Exit(errs)
}

func parseConfig() (config contracts.Config) {
	stringFlag("templates", "Directory with html templates.  ", "templates", &config.TemplateDir)
	stringFlag("styles   ", "Directory with css stylesheets. ", "css      ", &config.StylesDir)
	stringFlag("content  ", "Directory with markdown content.", "content  ", &config.ContentRoot)
	stringFlag("target   ", "Directory for rendered html.    ", "rendered ", &config.TargetRoot)
	boolFlag("with-drafts", "When set, include drafts.             ", false, &config.BuildDrafts)
	boolFlag("with-future", "When set, include future articles.    ", false, &config.BuildFuture)
	flag.Parse()

	if !Validate(config) {
		flag.PrintDefaults()
		os.Exit(1)
	}

	return config
}

func Validate(config contracts.Config) bool {
	if config.TemplateDir == "" {
		return false
	}
	if config.StylesDir == "" {
		return false
	}
	if config.ContentRoot == "" {
		return false
	}
	if config.TargetRoot == "" {
		return false
	}
	return true
}

func stringFlag(name, description, value string, s *string) {
	flag.StringVar(s,
		strings.TrimSpace(name),
		strings.TrimSpace(value),
		strings.TrimSpace(description),
	)
}

func boolFlag(name, description string, value bool, b *bool) {
	flag.BoolVar(b,
		strings.TrimSpace(name),
		value,
		strings.TrimSpace(description),
	)
}

func parseTemplates(glob string) *template.Template {
	templates, err := template.ParseGlob(glob)
	if err != nil {
		log.Fatalln("Could not parse templates:", err)
	}
	return templates
}

func exportCSS(disk contracts.FileSystem, config contracts.Config) error {
	_, err := disk.Stat(config.StylesDir)
	if os.IsNotExist(err) {
		return nil
	}
	return disk.Walk(config.StylesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return contracts.NewStackTraceError(err)
		}

		src, err := disk.Open(path)
		if err != nil {
			return contracts.NewStackTraceError(err)
		}

		defer func() { _ = src.Close() }()

		dst, err := disk.Create(filepath.Join(config.TargetRoot, "css", info.Name()))
		if err != nil {
			return contracts.NewStackTraceError(err)
		}

		defer func() { _ = dst.Close() }()

		_, err = io.Copy(dst, src)
		return err
	})
}
