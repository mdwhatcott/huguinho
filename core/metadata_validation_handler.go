package core

import (
	"github.com/mdwhatcott/huguinho/contracts"
)

type MetadataValidationHandler struct {
	slugs map[string]struct{}
}

func NewMetadataValidationHandler() *MetadataValidationHandler {
	return &MetadataValidationHandler{slugs: make(map[string]struct{})}
}

func (this *MetadataValidationHandler) Handle(article *contracts.Article) {
	if article.Metadata.Title == "" {
		article.Error = contracts.NewStackTraceError(errBlankMetadataTitle)
		return
	}

	if article.Metadata.Slug == "" {
		article.Error = contracts.NewStackTraceError(errBlankMetadataSlug)
		return
	}

	if article.Metadata.Date.IsZero() {
		article.Error = contracts.NewStackTraceError(errBlankMetadataDate)
		return
	}

	_, found := this.slugs[article.Metadata.Slug]
	if found {
		article.Error = contracts.NewStackTraceError(errRepeatedMetadataSlug)
		return
	}

	this.slugs[article.Metadata.Slug] = struct{}{}
}
