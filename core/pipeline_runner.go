package core

import (
	"github.com/mdwhatcott/huguinho/contracts"
)

type PipelineRunner struct {
	version string
	args    []string
	fs      contracts.FileSystem
	now     contracts.Clock
	log     contracts.Logger
}

func NewPipelineRunner(
	version string,
	args []string,
	fs contracts.FileSystem,
	now contracts.Clock,
	log contracts.Logger,
) *PipelineRunner {
	return &PipelineRunner{
		version: version,
		args:    args,
		fs:      fs,
		now:     now,
		log:     log,
	}
}

func (this *PipelineRunner) Run() (errors int) {
	start := this.now()

	config, err := NewCLIParser(this.version, this.args).Parse()
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

	pipeline := NewPipeline(this.now, config, this.fs, renderer)
	reporter := NewReporter(start, this.log)
	reporter.ProcessStream(pipeline.Run())
	reporter.RenderFinalReport(this.now())
	return reporter.Errors()
}
