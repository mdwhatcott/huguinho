package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/mdwhatcott/huguinho/contracts"
)

type PageParser struct {
	converter contracts.ContentConverter
}

func NewPageParser(converter contracts.ContentConverter) *PageParser {
	return &PageParser{converter: converter}
}

func (this *PageParser) ParsePage(file ContentFile) (page contracts.Page, err error) {
	divider := contracts.FRONT_MATTER_DIVIDER
	frontMatterEnd := strings.Index(file.Content, divider)
	contentStart := frontMatterEnd + len(divider)

	if frontMatterEnd < 0 {
		return page, fmt.Errorf("%w: [%s]", errMissingFrontMatter, file.Path)
	}
	rawFrontMatter := strings.TrimSpace(file.Content[:frontMatterEnd])
	rawContent := strings.TrimSpace(file.Content[contentStart:])

	_ = json.Unmarshal([]byte(rawFrontMatter), &page.Metadata)
	page.Content.Original = rawContent
	page.Content.Converted, _ = this.converter.Convert(rawContent)
	return page, err
}

var errMissingFrontMatter = errors.New("missing front matter divider")
