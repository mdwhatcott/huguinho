package core

import (
	"strings"

	"github.com/mdwhatcott/huguinho/contracts"
)

type MetadataParsingHandler struct{}

func NewMetadataParsingHandler() *MetadataParsingHandler {
	return &MetadataParsingHandler{}
}

func (this *MetadataParsingHandler) Handle(article *contracts.Article) error {
	if strings.TrimSpace(article.Source.Data) == "" {
		return contracts.NewStackTraceError(errMissingMetadata)
	}

	metadata, _ := divide(article.Source.Data, contracts.METADATA_CONTENT_DIVIDER)
	if len(metadata) == 0 {
		return contracts.NewStackTraceError(errMissingMetadataDivider)
	}

	parser := NewMetadataParser(strings.Split(metadata, "\n"))
	err := parser.Parse()
	if err != nil {
		return err
	}

	article.Metadata = parser.Parsed()
	return nil
}
