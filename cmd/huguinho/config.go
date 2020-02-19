package main

import (
	"flag"
	"os"
	"strings"
)

type Config struct {
	TemplateDir string
	StylesDir   string
	ContentRoot string
	TargetRoot  string
	BuildDrafts bool
	BuildFuture bool
}

func parseConfig() (config Config) {
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

func Validate(config Config) bool {
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
