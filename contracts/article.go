package contracts

import "time"

type Article struct {
	Error    error
	Source   ArticleSource
	Metadata ArticleMetadata
	Content  ArticleContent
}

type ArticleSource struct {
	Path string
	Data string
}

type ArticleMetadata struct {
	Draft  bool
	Slug   string
	Title  string
	Intro  string
	Topics []string
	Date   time.Time
}

const METADATA_CONTENT_DIVIDER = "\n+++\n"

type ArticleContent struct {
	Original  string
	Converted string
}

type (
	Filter func(*Article) bool
	Sorter func(i, j RenderedArticleSummary) int
)
