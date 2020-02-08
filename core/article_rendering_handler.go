package core

import (
	"path/filepath"

	"github.com/mdwhatcott/huguinho/contracts"
)

type RenderingFileSystem interface {
	contracts.MkdirAll
	contracts.WriteFile
}

type ArticleRenderingHandler struct {
	disk     RenderingFileSystem
	renderer contracts.Renderer
	output   string
}

func NewArticleRenderingHandler(
	disk RenderingFileSystem,
	renderer contracts.Renderer,
	output string,
) *ArticleRenderingHandler {
	return &ArticleRenderingHandler{
		disk:     disk,
		renderer: renderer,
		output:   output,
	}
}

func (this *ArticleRenderingHandler) Handle(article *contracts.Article) (err error) {
	rendered, err := this.renderer.Render(article)
	if err != nil {
		return NewStackTraceError(err)
	}

	folder := filepath.Join(this.output, article.Metadata.Slug)
	err = this.disk.MkdirAll(folder, 0755)
	if err != nil {
		return NewStackTraceError(err)
	}

	err = this.disk.WriteFile(filepath.Join(folder, "index.html"), []byte(rendered), 0644)
	if err != nil {
		return NewStackTraceError(err)
	}

	return nil
}
