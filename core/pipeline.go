package core

import (
	"fmt"
	"path/filepath"

	"github.com/mdwhatcott/huguinho/contracts"
)

type Pipeline struct {
	clock    contracts.Clock
	config   contracts.Config
	disk     contracts.FileSystem
	renderer contracts.Renderer
}

func NewPipeline(
	clock contracts.Clock,
	config contracts.Config,
	disk contracts.FileSystem,
	renderer contracts.Renderer,
) *Pipeline {
	return &Pipeline{
		clock:    clock,
		config:   config,
		disk:     disk,
		renderer: renderer,
	}
}
func (this *Pipeline) Run() (out chan contracts.Article) {
	home := NewListRenderingHandler(
		this.filterLastYear,
		this.sortByDateDescending,
		this.renderer,
		this.disk,
		this.config.TargetRoot,
		this.config.Author,
		"Here's what I've been working on lately:",
	)
	archive := NewListRenderingHandler(
		this.filterAll,
		this.sortByDateDescending,
		this.renderer,
		this.disk,
		filepath.Join(this.config.TargetRoot, "archives"),
		fmt.Sprintf("%s - Archives", this.config.Author),
		"Here's a complete history of my writings:",
	)
	var years []*ListRenderingHandler
	for year := 2000; year <= this.clock().Year(); year++ {
		years = append(years, NewListRenderingHandler(
			this.filterCalendarYear(year),
			this.sortByDate,
			this.renderer,
			this.disk,
			filepath.Join(this.config.TargetRoot, fmt.Sprint(year)),
			fmt.Sprintf("%s - %d", this.config.Author, year),
			fmt.Sprintf("Here's what I wrote in %d:", year),
		))
	}
	out = this.goLoad()
	out = this.goListen(out, NewFileReadingHandler(this.disk))
	out = this.goListen(out, NewMetadataParsingHandler())
	out = this.goListen(out, NewMetadataValidationHandler())
	out = this.goListen(out, NewDraftFilteringHandler(!this.config.BuildDrafts))
	out = this.goListen(out, NewFutureFilteringHandler(this.clock(), !this.config.BuildFuture))
	out = this.goListen(out, NewContentConversionHandler(NewGoldmarkMarkdownConverter()))
	out = this.goListen(out, NewArticleRenderingHandler(this.disk, this.renderer, this.config.TargetRoot))
	out = this.goListen(out, NewTopicPageRenderingHandler(this.disk, this.renderer, this.config.TargetRoot))
	out = this.goListen(out, home)
	out = this.goListen(out, archive)
	for _, handler := range years {
		out = this.goListen(out, handler)
	}
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
func (this *Pipeline) filterAll(*contracts.Article) bool { return true }
func (this *Pipeline) filterLastYear(a *contracts.Article) bool {
	return a.Metadata.Date.After(this.clock().AddDate(-1, 0, 0))
}
func (this *Pipeline) sortByDateDescending(i, j contracts.RenderedArticleSummary) int {
	return int(j.Date.UnixNano() - i.Date.UnixNano())
}
func (this *Pipeline) sortByDate(i, j contracts.RenderedArticleSummary) int {
	return int(i.Date.UnixNano() - j.Date.UnixNano())
}
func (this *Pipeline) filterCalendarYear(year int) contracts.Filter {
	return func(article *contracts.Article) bool {
		return article.Metadata.Date.Year() == year
	}
}
