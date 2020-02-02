package shell

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

type GoldmarkMarkdownConverter struct {
	converter goldmark.Markdown
}

func NewGoldmarkMarkdownConverter() *GoldmarkMarkdownConverter {
	return &GoldmarkMarkdownConverter{
		converter: goldmark.New(goldmark.WithRendererOptions(html.WithUnsafe())),
	}
}

func (this *GoldmarkMarkdownConverter) Convert(content string) (string, error) {
	buffer := bytes.NewBufferString("")
	err := this.converter.Convert([]byte(content), buffer)
	return buffer.String(), err
}
