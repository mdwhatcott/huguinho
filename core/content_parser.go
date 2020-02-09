package core

import (
	"strings"

	"github.com/mdwhatcott/huguinho/contracts"
)

type ContentParser struct {
	inner contracts.ContentConverter
}

func NewContentParser(inner contracts.ContentConverter) *ContentParser {
	return &ContentParser{inner: inner}
}

func (this *ContentParser) Handle(article *contracts.Article) (err error) {
	marker := contracts.METADATA_CONTENT_DIVIDER
	divider := strings.Index(article.Source.Data, marker)
	original := strings.TrimSpace(article.Source.Data[divider+len(marker):])
	converted, err := this.inner.Convert(original)
	if err != nil {
		return NewStackTraceError(err)
	}

	article.Content.Original = original
	article.Content.Converted = converted
	return nil
}
