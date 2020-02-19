package core

import (
	"path/filepath"
	"sort"

	"github.com/mdwhatcott/huguinho/contracts"
)

type TagPageRenderingHandler struct {
	disk     RenderingFileSystem
	renderer contracts.Renderer
	output   string
	tags     map[string]contracts.RenderedTagListing
}

func NewTagPageRenderingHandler(
	disk RenderingFileSystem,
	renderer contracts.Renderer,
	output string,
) *TagPageRenderingHandler {
	return &TagPageRenderingHandler{
		disk:     disk,
		renderer: renderer,
		output:   output,
		tags:     make(map[string]contracts.RenderedTagListing),
	}
}

func (this *TagPageRenderingHandler) Handle(article *contracts.Article) error {
	for _, tag := range article.Metadata.Tags {
		listing := this.tags[tag]
		listing.Title = tag
		listing.Name = tag
		listing.Pages = append(listing.Pages, contracts.RenderedTagEntry{
			Slug:  article.Metadata.Slug,
			Title: article.Metadata.Title,
			Date:  article.Metadata.Date,
		})
		this.tags[tag] = listing
	}
	return nil
}

func (this *TagPageRenderingHandler) Finalize() error {
	for _, listing := range this.tags {
		sort.Slice(listing.Pages, func(i, j int) bool {
			return listing.Pages[i].Date.After(listing.Pages[j].Date)
		})

		rendered, err := this.renderer.Render(listing)
		if err != nil {
			return contracts.NewStackTraceError(err)
		}

		folder := filepath.Join(this.output, "tags", listing.Name)

		err = this.disk.MkdirAll(folder, 0755)
		if err != nil {
			return contracts.NewStackTraceError(err)
		}

		err = this.disk.WriteFile(filepath.Join(folder, "index.html"), []byte(rendered), 0644)
		if err != nil {
			return contracts.NewStackTraceError(err)
		}
	}
	return nil
}
