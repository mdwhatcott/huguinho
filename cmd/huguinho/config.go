package main

import (
	"flag"
	"os"
	"strings"
)

func ParseCLI() *Config {
	config := new(Config)
	config.ParseCLI()

	if !config.Validate() {
		flag.PrintDefaults()
		os.Exit(1)
	}

	return config
}

type Config struct {
	templateDir string
	stylesDir   string
	contentRoot string
	targetRoot  string
	buildDrafts bool
	buildFuture bool
}

func (this *Config) ParseCLI() {
	stringFlag("templates", "Directory with html templates.  ", "templates", &this.templateDir)
	stringFlag("styles   ", "Directory with css stylesheets. ", "css      ", &this.stylesDir)
	stringFlag("content  ", "Directory with markdown content.", "content  ", &this.contentRoot)
	stringFlag("target   ", "Directory for rendered html.    ", "rendered ", &this.targetRoot)
	boolFlag("with-drafts", "When set, include drafts.             ", false, &this.buildDrafts)
	boolFlag("with-future", "When set, include future articles.    ", false, &this.buildFuture)
	flag.Parse()
}

func (this *Config) Validate() bool {
	if this.templateDir == "" {
		return false
	}
	if this.stylesDir == "" {
		return false
	}
	if this.contentRoot == "" {
		return false
	}
	if this.targetRoot == "" {
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
