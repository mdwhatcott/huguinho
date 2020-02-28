package core

import (
	"github.com/smartystreets/clock"
	"github.com/smartystreets/logging"

	"github.com/mdwhatcott/huguinho/contracts"
)

type PipelineRunner struct {
	clock *clock.Clock
	log   *logging.Logger
	args  []string
	fs    contracts.FileSystem
}

func NewPipelineRunner(args []string, fs contracts.FileSystem) *PipelineRunner {
	return &PipelineRunner{args: args, fs: fs}
}

func (this *PipelineRunner) Run() (errors int) {
	start := this.clock.UTCNow()

	config, err := NewCLIParser(this.args).Parse()
	if err != nil {
		this.log.Println(err)
		return 1
	}

	loader := NewTemplateLoader(this.fs, config.TemplateDir)
	templates, err := loader.Load()
	if err != nil {
		this.log.Println(err)
		return 1
	}

	renderer := NewTemplateRenderer(templates)
	err = renderer.Validate()
	if err != nil {
		this.log.Println(err)
		return 1
	}

	pipeline := NewPipeline(config, this.fs, renderer)
	reporter := NewReporter(start)
	reporter.log = this.log
	reporter.ProcessStream(pipeline.Run())
	reporter.RenderFinalReport(this.clock.UTCNow())
	return reporter.Errors()
}
