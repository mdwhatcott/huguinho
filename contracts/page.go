package contracts

import "time"

type Page struct {
	FrontMatter
	Path            Path
	OriginalContent string
	HTMLContent     string
}

type FrontMatter struct {
	ParseError  error     `toml:"-"`
	Title       string    `toml:"title"`
	Description string    `toml:"description"`
	Date        time.Time `toml:"date"`
	Tags        []string  `toml:"tags"`
	IsDraft     bool      `toml:"draft"`
}
