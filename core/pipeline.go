package core

import (
	"log"
	"time"

	"github.com/mdwhatcott/huguinho/contracts"
	"github.com/mdwhatcott/huguinho/shell"
)

type Pipeline struct {
	config   contracts.Config
	disk     contracts.FileSystem
	renderer contracts.Renderer
	errs     chan error
}

func NewPipeline(
	config contracts.Config,
	disk contracts.FileSystem,
	renderer contracts.Renderer,
) *Pipeline {
	return &Pipeline{
		config:   config,
		disk:     disk,
		renderer: renderer,
		errs:     make(chan error),
	}
}
func (this *Pipeline) Run() int {
	out := this.startAll()
	go this.terminate(out)
	return this.errCount()
}
func (this *Pipeline) startAll() (out chan contracts.Article) {
	out = this.goLoad()
	out = this.goListen(out, NewFileReadingHandler(this.disk))
	out = this.goListen(out, NewMetadataParsingHandler())
	out = this.goListen(out, NewMetadataValidationHandler())
	out = this.goListen(out, NewDraftFilteringHandler(!this.config.BuildDrafts))
	out = this.goListen(out, NewFutureFilteringHandler(time.Now(), !this.config.BuildFuture))
	out = this.goListen(out, NewContentParsingHandler(shell.NewGoldmarkMarkdownConverter()))
	out = this.goListen(out, NewArticleRenderingHandler(this.disk, this.renderer, this.config.TargetRoot))
	out = this.goListen(out, NewTagPageRenderingHandler(this.disk, this.renderer, this.config.TargetRoot))
	out = this.goListen(out, NewHomePageRenderingHandler(this.disk, this.renderer, this.config.TargetRoot))
	// out = this.goListen(out, core.NewAllTagsRenderer(...)) // TODO
	return out
}
func (this *Pipeline) goLoad() (out chan contracts.Article) {
	out = make(chan contracts.Article)
	loader := NewPathLoader(this.disk, this.config.ContentRoot, out)
	go func() { this.errs <- loader.Start() }()
	return out
}
func (this *Pipeline) goListen(in chan contracts.Article, handler contracts.Handler) (out chan contracts.Article) {
	out = make(chan contracts.Article)
	listener := NewListener(in, out, handler)
	go func() { this.errs <- listener.Listen() }()
	return out
}
func (this *Pipeline) terminate(out chan contracts.Article) {
	for item := range out {
		log.Println("[INFO] published article:", item.Metadata.Slug)
	}
	close(this.errs)
}
func (this *Pipeline) errCount() (errCount int) {
	for err := range this.errs {
		if err != nil {
			errCount++
			log.Println("[WARN]", err)
		}
	}
	return errCount
}
