package core

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/mdwhatcott/huguinho/contracts"
)

type JSONMetadata struct {
	Draft bool      `json:"draft"`
	Slug  string    `json:"slug"` // TODO: schema change (all articles)
	Title string    `json:"title"`
	Intro string    `json:"intro"` // TODO: schema change (all articles)
	Tags  []string  `json:"tags"`
	Date  time.Time `json:"date"`
}

type JSONMetadataParser struct{}

func NewJSONMetadataParser() *JSONMetadataParser {
	return &JSONMetadataParser{}
}

func (this *JSONMetadataParser) Handle(article *contracts.Article) error {
	division := strings.Index(article.Source.Data, contracts.METADATA_CONTENT_DIVIDER)
	// TODO: no division err
	raw := []byte(article.Source.Data[:division])
	var metadata JSONMetadata
	_ = json.Unmarshal(raw, &metadata) // TODO: err
	article.Metadata = contracts.ArticleMetadata{
		Draft: metadata.Draft,
		Slug:  metadata.Slug,
		Title: metadata.Title,
		Intro: metadata.Intro,
		Tags:  metadata.Tags,
		Date:  metadata.Date,
	}
	// TODO: validation? (or separate handler?)
	return nil
}
