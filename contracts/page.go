package contracts

import "time"

type Page struct {
	FrontMatter
	OriginalContent string
	HTMLContent     string
	ParseError      error
}

type FrontMatter struct {
	Title       string    `toml:"title"`
	Description string    `toml:"description"`
	Date        time.Time `toml:"date"`
	Tags        []string  `toml:"tags"`
	IsDraft     bool      `toml:"draft"`
}
