package contracts

import "time"

type Page struct {
	SourcePath string
	Metadata   JSONFrontMatter
	Content    Content
}

type JSONFrontMatter struct {
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Tags        []string  `json:"tags"`
	IsDraft     bool      `json:"draft"`
}

type Content struct {
	Original  string
	Converted string
}

const FRONT_MATTER_DIVIDER = "+++"
