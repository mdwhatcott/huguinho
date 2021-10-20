package core

import (
	"path/filepath"
	"sort"

	"github.com/mdwhatcott/huguinho/contracts"
)

type HomePageRenderingHandler struct {
	disk     RenderingFileSystem
	renderer contracts.Renderer
	output   string
	listing  []contracts.RenderedArticleSummary
}

func NewHomePageRenderingHandler(
	disk RenderingFileSystem,
	renderer contracts.Renderer,
	output string,
) *HomePageRenderingHandler {
	return &HomePageRenderingHandler{
		disk:     disk,
		renderer: renderer,
		output:   output,
	}
}

func (this *HomePageRenderingHandler) Handle(article *contracts.Article) {
	this.listing = append(this.listing, contracts.RenderedArticleSummary{
		Slug:   article.Metadata.Slug,
		Title:  article.Metadata.Title,
		Intro:  article.Metadata.Intro,
		Date:   article.Metadata.Date,
		Topics: article.Metadata.Topics,
		Draft:  article.Metadata.Draft,
	})
}

func (this *HomePageRenderingHandler) Finalize() error {
	sort.Slice(this.listing, func(i, j int) bool {
		return this.listing[i].Date.After(this.listing[j].Date)
	})

	rendered, err := this.renderer.Render(contracts.RenderedHomePage{Pages: this.listing})
	if err != nil {
		return StackTraceError(err)
	}

	err = this.disk.MkdirAll(this.output, 0755)
	if err != nil {
		return StackTraceError(err)
	}

	err = this.disk.WriteFile(filepath.Join(this.output, "index.html"), []byte(rendered), 0644)
	if err != nil {
		return StackTraceError(err)
	}

	return nil
}
