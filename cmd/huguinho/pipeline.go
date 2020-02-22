package main

import (
	"time"

	"github.com/mdwhatcott/huguinho/contracts"
	"github.com/mdwhatcott/huguinho/core"
	"github.com/mdwhatcott/huguinho/shell"
)

// TEST (integration)
type Pipeline struct {
	config   Config
	disk     contracts.FileSystem
	renderer contracts.Renderer
}

func NewPipeline(
	config Config,
	disk contracts.FileSystem,
	renderer contracts.Renderer,
) *Pipeline {
	return &Pipeline{
		config:   config,
		disk:     disk,
		renderer: renderer,
	}
}
func (this *Pipeline) Run() (out chan contracts.Article) {
	out = this.goLoad()
	out = this.goListen(out, core.NewFileReadingHandler(this.disk))
	out = this.goListen(out, core.NewMetadataParsingHandler())
	out = this.goListen(out, core.NewMetadataValidationHandler())
	out = this.goListen(out, core.NewDraftFilteringHandler(!this.config.BuildDrafts))
	out = this.goListen(out, core.NewFutureFilteringHandler(time.Now(), !this.config.BuildFuture))
	out = this.goListen(out, core.NewContentParsingHandler(shell.NewGoldmarkMarkdownConverter()))
	out = this.goListen(out, core.NewArticleRenderingHandler(this.disk, this.renderer, this.config.TargetRoot))
	out = this.goListen(out, core.NewTopicPageRenderingHandler(this.disk, this.renderer, this.config.TargetRoot))
	out = this.goListen(out, core.NewHomePageRenderingHandler(this.disk, this.renderer, this.config.TargetRoot))
	return out
}
func (this *Pipeline) goLoad() (out chan contracts.Article) {
	out = make(chan contracts.Article)
	go core.NewPathLoader(this.disk, this.config.ContentRoot, out).Start()
	return out
}
func (this *Pipeline) goListen(in chan contracts.Article, handler contracts.Handler) (out chan contracts.Article) {
	out = make(chan contracts.Article)
	go core.Listen(in, out, handler)
	return out
}
