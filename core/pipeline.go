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

	published int
	errors    int
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
	go NewPathLoader(this.disk, this.config.ContentRoot, out).Start()
	return out
}
func (this *Pipeline) goListen(in chan contracts.Article, handler contracts.Handler) (out chan contracts.Article) {
	out = make(chan contracts.Article)
	go Listen(in, out, handler)
	return out
}

func (this *Pipeline) drain(out chan contracts.Article) {
	for article := range out {
		if article.Error != nil {
			log.Println("[WARN] error:", article.Error)
			this.errors++
		} else if article.Source.Path != "" {
			log.Println("[INFO] published article:", article.Metadata.Slug)
			this.published++
		} else {
			log.Printf("[INFO] not sure what this article struct represents: %#v", article)
		}
	}
}
