package core

import (
	"bytes"
	"errors"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/huguinho/contracts"
)

func TestCLIParserFixture(t *testing.T) {
	gunit.Run(new(CLIParserFixture), t)
}

type CLIParserFixture struct {
	*gunit.Fixture

	output *bytes.Buffer
	args   []string
}

func (this *CLIParserFixture) Setup() {
	this.output = new(bytes.Buffer)
}

func (this *CLIParserFixture) Parse() (contracts.Config, error) {
	parser := NewCLIParser("version", this.args)
	parser.flags.SetOutput(this.output)
	return parser.Parse()
}

func (this *CLIParserFixture) TestDefaults() {
	this.args = []string{}
	config, err := this.Parse()
	this.So(err, should.BeNil)
	this.So(config, should.Resemble, contracts.Config{
		TemplateDir: "templates",
		ContentRoot: "content",
		TargetRoot:  "rendered",
		BuildDrafts: false,
		BuildFuture: false,
	})
}

func (this *CLIParserFixture) TestCustomValues() {
	this.args = []string{
		"-templates", "other-templates",
		"-content", "other-content",
		"-target", "other-rendered",
		"-with-drafts",
		"-with-future",
	}
	config, err := this.Parse()
	this.So(err, should.BeNil)
	this.So(config, should.Resemble, contracts.Config{
		TemplateDir: "other-templates",
		ContentRoot: "other-content",
		TargetRoot:  "other-rendered",
		BuildDrafts: true,
		BuildFuture: true,
	})
}

func (this *CLIParserFixture) TestMissingTemplatesFolder() {
	this.args = []string{"-templates", ""}
	config, err := this.Parse()
	this.So(errors.Is(err, ErrInvalidConfig), should.BeTrue)
	this.So(config, should.BeZeroValue)
}

func (this *CLIParserFixture) TestMissingContentFolder() {
	this.args = []string{"-content", ""}
	config, err := this.Parse()
	this.So(errors.Is(err, ErrInvalidConfig), should.BeTrue)
	this.So(config, should.BeZeroValue)
}

func (this *CLIParserFixture) TestMissingTargetFolder() {
	this.args = []string{"-target", ""}
	config, err := this.Parse()
	this.So(errors.Is(err, ErrInvalidConfig), should.BeTrue)
	this.So(config, should.BeZeroValue)
}

func (this *CLIParserFixture) TestBogusValue() {
	this.args = []string{"-bogus"}
	config, err := this.Parse()
	this.So(errors.Is(err, ErrInvalidConfig), should.BeTrue)
	this.So(config, should.BeZeroValue)
}
