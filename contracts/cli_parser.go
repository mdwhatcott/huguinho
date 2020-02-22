package contracts

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"strings"
)

type CLIParser struct {
	args  []string
	flags *flag.FlagSet
}

func NewCLIParser(args []string) *CLIParser {
	return &CLIParser{
		args:  args,
		flags: flag.NewFlagSet("huguinho", flag.ContinueOnError),
	}
}

func (this *CLIParser) Parse() (config Config, err error) {
	this.stringFlag("templates", "Directory with html templates.  ", "templates", &config.TemplateDir)
	this.stringFlag("content  ", "Directory with markdown content.", "content  ", &config.ContentRoot)
	this.stringFlag("target   ", "Directory for rendered html.    ", "rendered ", &config.TargetRoot)
	this.boolFlag("with-drafts", "When set, include drafts.             ", false, &config.BuildDrafts)
	this.boolFlag("with-future", "When set, include future articles.    ", false, &config.BuildFuture)

	err = this.flags.Parse(this.args)
	if err != nil {
		return Config{}, this.composeError(err)
	}

	err = config.Validate()
	if err != nil {
		return Config{}, this.composeError(err)
	}

	return config, nil
}

func (this *CLIParser) composeError(err error) error {
	buffer := new(bytes.Buffer)
	this.flags.SetOutput(buffer)
	this.flags.PrintDefaults()
	return fmt.Errorf("%w: %v\n%s", ErrInvalidConfig, err, buffer.String())
}

var ErrInvalidConfig = errors.New("invalid config")

func (this *CLIParser) stringFlag(name, description, value string, s *string) {
	this.flags.StringVar(s,
		strings.TrimSpace(name),
		strings.TrimSpace(value),
		strings.TrimSpace(description),
	)
}

func (this *CLIParser) boolFlag(name, description string, value bool, b *bool) {
	this.flags.BoolVar(b,
		strings.TrimSpace(name),
		value,
		strings.TrimSpace(description),
	)
}
