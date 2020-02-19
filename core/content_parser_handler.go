package core

import (
	"github.com/mdwhatcott/huguinho/contracts"
)

type ContentParsingHandler struct {
	inner contracts.ContentConverter
}

func NewContentParsingHandler(inner contracts.ContentConverter) *ContentParsingHandler {
	return &ContentParsingHandler{inner: inner}
}

func (this *ContentParsingHandler) Handle(article *contracts.Article) {
	_, original := divide(article.Source.Data, contracts.METADATA_CONTENT_DIVIDER)
	converted, err := this.inner.Convert(original)
	if err != nil {
		article.Error = contracts.NewStackTraceError(err)
		return
	}

	article.Content.Original = original
	article.Content.Converted = converted
}
