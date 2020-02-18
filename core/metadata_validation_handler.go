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

func (this *MetadataValidationHandler) Handle(article *contracts.Article) error {
	if article.Metadata.Title == "" {
		return contracts.NewStackTraceError(errBlankMetadataTitle)
	}
	if article.Metadata.Slug == "" {
		return contracts.NewStackTraceError(errBlankMetadataSlug)
	}
	if article.Metadata.Date.IsZero() {
		return contracts.NewStackTraceError(errBlankMetadataDate)
	}
	_, found := this.slugs[article.Metadata.Slug]
	if found {
		return contracts.NewStackTraceError(errRepeatedMetadataSlug)
	}
	this.slugs[article.Metadata.Slug] = struct{}{}
	return nil
}
