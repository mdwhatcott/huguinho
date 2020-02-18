package core

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/mdwhatcott/huguinho/contracts"
)

type MetadataParser struct {
	lines  []string
	parsed contracts.ArticleMetadata

	parsedTitle bool
	parsedIntro bool
	parsedSlug  bool
	parsedDraft bool
	parsedDate  bool
	parsedTags  bool
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
		case "date":
			err := this.parseDate(value)
			if err != nil {
				return err
			}
		case "tags":
			err := this.parseTags(value)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (this *MetadataParser) parseTitle(value string) error {
	if this.parsedTitle {
		return contracts.NewStackTraceError(errDuplicateMetadataTitle)
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
		return contracts.NewStackTraceError(errDuplicateMetadataIntro)
	}
	if value == "" {
		return contracts.NewStackTraceError(errBlankMetadataIntro)
	}
	this.parsed.Intro = value
	this.parsedIntro = true
	return nil
}
func (this *MetadataParser) parseSlug(value string) error {
	if this.parsedSlug {
		return contracts.NewStackTraceError(errDuplicateMetadataSlug)
	}
	if value == "" {
		return contracts.NewStackTraceError(errBlankMetadataSlug)
	}
	if strings.ToLower(value) != value {
		return contracts.NewStackTraceError(errInvalidMetadataSlug)
	}
	parsed, _ := url.Parse(value)
	if parsed.Path != parsed.EscapedPath() {
		return contracts.NewStackTraceError(fmt.Errorf("%w: [%s]", errInvalidMetadataSlug, value))
	}
	this.parsed.Slug = value
	this.parsedSlug = true
	return nil
}
func (this *MetadataParser) parseDraft(value string) error {
	if this.parsedDraft {
		return contracts.NewStackTraceError(errDuplicateMetadataDraft)
	}

	switch value {
	case "true":
		this.parsed.Draft = true
		this.parsedDraft = true
	case "false":
		this.parsed.Draft = false
		this.parsedDraft = true
	case "":
		return contracts.NewStackTraceError(errBlankMetadataDraft)
	default:
		return contracts.NewStackTraceError(fmt.Errorf("%w: [%s]", errInvalidMetadataDraft, value))
	}
	return nil
}
func (this *MetadataParser) parseDate(value string) error {
	if this.parsedDate {
		return contracts.NewStackTraceError(errDuplicateMetadataDate)
	}
	if value == "" {
		return contracts.NewStackTraceError(errBlankMetadataDate)
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return contracts.NewStackTraceError(fmt.Errorf("%w with value: [%s] err: %v", errInvalidMetadataDate, value, err))
	}
	this.parsed.Date = parsed
	this.parsedDate = true
	return nil
}
func (this *MetadataParser) parseTags(value string) error {
	if this.parsedTags {
		return contracts.NewStackTraceError(errDuplicateMetadataTags)
	}
	if value == "" {
		return contracts.NewStackTraceError(errBlankMetadataTags)
	}
	unique := make(map[string]struct{})
	tags := strings.Fields(value)
	for _, tag := range tags {
		if !isValidTag(tag) {
			return contracts.NewStackTraceError(fmt.Errorf("%w: [%s]", errInvalidMetadataTags, value))
		}
		unique[tag] = struct{}{}
	}
	if len(unique) != len(tags) {
		return contracts.NewStackTraceError(fmt.Errorf("%w: [%s] (repeated values)", errInvalidMetadataTags, value))
	}
	this.parsed.Tags = tags
	this.parsedTags = true
	return nil
}

func isValidTag(tag string) bool {
	for _, c := range tag {
		if !(isSpace(c) || isDash(c) || isNumber(c) || isLowerAlpha(c)) {
			return false
		}
	}
	return true
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
	errDuplicateMetadataDate  = errors.New("duplicate metadata date")
	errDuplicateMetadataTags  = errors.New("duplicate metadata tags")

	errInvalidMetadataSlug  = errors.New("invalid metadata slug")
	errInvalidMetadataDraft = errors.New("invalid metadata draft")
	errInvalidMetadataDate  = errors.New("invalid metadata date")
	errInvalidMetadataTags  = errors.New("invalid metadata tags")

	errRepeatedMetadataSlug = errors.New("repeated metadata slug")

	errBlankMetadataSlug  = errors.New("blank metadata slug")
	errBlankMetadataDraft = errors.New("blank metadata draft")
	errBlankMetadataTitle = errors.New("blank metadata title")
	errBlankMetadataIntro = errors.New("blank metadata intro")
	errBlankMetadataDate  = errors.New("blank metadata date")
	errBlankMetadataTags  = errors.New("blank metadata tags")
)
