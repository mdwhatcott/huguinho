package core

import (
	"github.com/mdwhatcott/huguinho/contracts"
)

type ContentParser struct {
	inner contracts.ContentConverter
}

func NewContentParser(inner contracts.ContentConverter) *ContentParser {
	return &ContentParser{inner: inner}
}

func (this *ContentParser) Handle(article *contracts.Article) (err error) {
	_, original := divide(article.Source.Data, contracts.METADATA_CONTENT_DIVIDER)
	converted, err := this.inner.Convert(original)
	if err != nil {
		return NewStackTraceError(err)
	}

	article.Content.Original = original
	article.Content.Converted = converted
	return nil
}
