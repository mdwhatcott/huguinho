package core

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/mdwhatcott/huguinho/contracts"
)

type MetadataParsingHandler struct{}

func NewMetadataParsingHandler() *MetadataParsingHandler {
	return &MetadataParsingHandler{}
}

func (this *MetadataParsingHandler) Handle(article *contracts.Article) error {
	if strings.TrimSpace(article.Source.Data) == "" {
		return NewStackTraceError(errMissingMetadata)
	}

	metadata, _ := divide(article.Source.Data, contracts.METADATA_CONTENT_DIVIDER)
	if len(metadata) == 0 {
		return NewStackTraceError(errMissingMetadataDivider)
	}

	parser := NewMetadataParser(strings.Split(metadata, "\n"))
	err := parser.Parse()
	if err != nil {
		return err
	}

	article.Metadata = parser.Parsed()
	return nil
}

type MetadataParser struct {
	lines  []string
	parsed contracts.ArticleMetadata

	parsedTitle bool
	parsedIntro bool
	parsedSlug  bool
	parsedDraft bool
}

func NewMetadataParser(lines []string) *MetadataParser {
	return &MetadataParser{lines: lines}
}

func (this *MetadataParser) Parse() error {
	for _, line := range this.lines {
		key, value := divide(line, ":")

		switch key {
		case "title":
			err := this.parseTitle(value)
			if err != nil {
				return err
			}
		case "intro":
			err := this.parseIntro(value)
			if err != nil {
				return err
			}
		case "slug":
			err := this.parseSlug(value)
			if err != nil {
				return err
			}
		case "draft":
			err := this.parseDraft(value)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (this *MetadataParser) parseTitle(value string) error {
	if this.parsedTitle {
		return NewStackTraceError(errDuplicateMetadataTitle)
	}
	if value == "" {
		return errBlankMetadataTitle
	}
	this.parsed.Title = value
	this.parsedTitle = true
	return nil
}
func (this *MetadataParser) parseIntro(value string) error {
	if this.parsedIntro {
		return NewStackTraceError(errDuplicateMetadataIntro)
	}
	if value == "" {
		return NewStackTraceError(errBlankMetadataIntro)
	}
	this.parsed.Intro = value
	this.parsedIntro = true
	return nil
}
func (this *MetadataParser) parseSlug(value string) error {
	if this.parsedSlug {
		return NewStackTraceError(errDuplicateMetadataSlug)
	}
	if value == "" {
		return NewStackTraceError(errBlankMetadataSlug)
	}
	parsed, _ := url.Parse(value)
	if parsed.Path != parsed.EscapedPath() {
		return NewStackTraceError(fmt.Errorf("%w: [%s]", errInvalidMetadataSlug, value))
	}
	this.parsed.Slug = value
	this.parsedSlug = true
	return nil
}
func (this *MetadataParser) parseDraft(value string) error {
	if this.parsedDraft {
		return NewStackTraceError(errDuplicateMetadataDraft)
	}

	switch value {
	case "true":
		this.parsed.Draft = true
		this.parsedDraft = true
	case "false":
		this.parsed.Draft = false
		this.parsedDraft = true
	case "":
		return NewStackTraceError(errBlankMetadataDraft)
	default:
		return NewStackTraceError(fmt.Errorf("%w: [%s]", errInvalidMetadataDraft, value))
	}
	return nil
}

func (this *MetadataParser) Parsed() contracts.ArticleMetadata {
	return this.parsed
}

var (
	errMissingMetadata        = errors.New("article lacks metadata")
	errMissingMetadataDivider = errors.New("article lacks metadata divider")

	errDuplicateMetadataTitle = errors.New("duplicate metadata title")
	errDuplicateMetadataIntro = errors.New("duplicate metadata intro")
	errDuplicateMetadataSlug  = errors.New("duplicate metadata slug")
	errDuplicateMetadataDraft = errors.New("duplicate metadata draft")

	errInvalidMetadataSlug  = errors.New("invalid metadata slug")
	errInvalidMetadataDraft = errors.New("invalid metadata draft")

	errBlankMetadataSlug  = errors.New("blank metadata slug")
	errBlankMetadataDraft = errors.New("blank metadata draft")
	errBlankMetadataTitle = errors.New("blank metadata title")
	errBlankMetadataIntro = errors.New("blank metadata intro")
)
