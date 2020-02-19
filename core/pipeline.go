package core

import (
	"log"
	"time"

	"github.com/mdwhatcott/huguinho/contracts"
	"github.com/mdwhatcott/huguinho/shell"
)

type Pipeline struct {
	config     contracts.Config
	disk       contracts.FileSystem
	renderer   contracts.Renderer
	finalizers []contracts.Finalizer
	published  int
	errors     int
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
	}
}
func (this *Pipeline) Run() (published, errors int) {
	this.drain(this.startAll())
	this.runFinalizer()
	return this.published, this.errors
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
	this.finalizers = append(this.finalizers, loader)
	go loader.Start()
	return out
}
func (this *Pipeline) goListen(in chan contracts.Article, handler contracts.Handler) (out chan contracts.Article) {
	finalizer, ok := handler.(contracts.Finalizer)
	if ok {
		this.finalizers = append(this.finalizers, finalizer)
	}
	out = make(chan contracts.Article)
	go Listen(in, out, handler)
	return out
}

func (this *Pipeline) runFinalizer() {
	for _, finalizer := range this.finalizers {
		err := finalizer.Finalize()
		if err != nil {
			log.Println("[WARN] handler finalization error:", err)
			this.errors++
		}
	}
}

func (this *Pipeline) drain(out chan contracts.Article) {
	for article := range out {
		if article.Error != nil {
			log.Println("[WARN] article handling error:", article.Error)
			this.errors++
		} else {
			log.Println("[INFO] published article:", article.Metadata.Slug)
			this.published++
		}
	}
}
