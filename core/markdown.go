package core

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

type GoldmarkMarkdownConverter struct {
	buffer    *bytes.Buffer
	converter goldmark.Markdown
}

func NewGoldmarkMarkdownConverter() *GoldmarkMarkdownConverter {
	return &GoldmarkMarkdownConverter{
		buffer: new(bytes.Buffer),
		converter: goldmark.New(
			goldmark.WithRendererOptions(
				html.WithUnsafe(),
			),
			goldmark.WithExtensions(
				extension.GFM,
			),
		),
	}
}

func (this *GoldmarkMarkdownConverter) Convert(content string) (string, error) {
	this.buffer.Reset()
	err := this.converter.Convert([]byte(content), this.buffer)
	return this.buffer.String(), err
}
