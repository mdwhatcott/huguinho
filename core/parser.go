package core

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/mdwhatcott/huguinho/contracts"
)

type PageParser struct {
	converter contracts.ContentConverter
}

func NewPageParser(converter contracts.ContentConverter) *PageParser {
	return &PageParser{converter: converter}
}

func (this *PageParser) ParsePage(content string) (page contracts.Page, err error) {
	divider := contracts.FRONT_MATTER_DIVIDER
	if !strings.Contains(content, divider) {
		return page, errors.New("missing front matter divider")
	}
	frontMatterEnd := strings.Index(content, divider)
	contentStart := frontMatterEnd + len(divider)
	rawFrontMatter := strings.TrimSpace(content[:frontMatterEnd])
	rawContent := strings.TrimSpace(content[contentStart:])
	_ = json.Unmarshal([]byte(rawFrontMatter), &page.Metadata)
	page.Content.Original = strings.TrimSpace(rawContent)
	page.Content.Converted, _ = this.converter.Convert(rawContent)
	return page, err
}
