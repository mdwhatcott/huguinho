package core

import (
	"path/filepath"
	"sort"

	"github.com/mdwhatcott/huguinho/contracts"
)

type ListRenderingHandler struct {
	listing  []contracts.RenderedArticleSummary
	filter   contracts.Filter
	sorter   contracts.Sorter
	renderer contracts.Renderer
	disk     RenderingFileSystem
	output   string
}

func NewListRenderingHandler(
	filter contracts.Filter,
	sorter contracts.Sorter,
	renderer contracts.Renderer,
	disk RenderingFileSystem,
	output string,
) *ListRenderingHandler {
	return &ListRenderingHandler{
		filter:   filter,
		sorter:   sorter,
		renderer: renderer,
		disk:     disk,
		output:   output,
	}
}
func (this *ListRenderingHandler) Handle(article *contracts.Article) {
	if !this.filter(article) {
		return
	}
	this.listing = append(this.listing, contracts.RenderedArticleSummary{
		Slug:   article.Metadata.Slug,
		Title:  article.Metadata.Title,
		Intro:  article.Metadata.Intro,
		Date:   article.Metadata.Date,
		Topics: article.Metadata.Topics,
		Draft:  article.Metadata.Draft,
	})
}
func (this *ListRenderingHandler) Finalize() error {
	sort.SliceStable(this.listing, func(i, j int) bool {
		return this.sorter(this.listing[i], this.listing[j]) < 0
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
