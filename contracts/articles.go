package contracts

import (
	"path/filepath"
	"strings"
	"time"
)

type Article struct {
	FrontMatter
	Path            Path
	Permalink       string // TODO: populate
	OriginalContent string
	HTMLContent     string
}

func (this Article) TargetPath(root string) string {
	folder := strings.TrimSuffix(string(this.Path), ".md")
	return filepath.Join(root, folder, "index.html")
}

type FrontMatter struct {
	ParseError  error     `toml:"-"`
	Title       string    `toml:"title"`
	Description string    `toml:"description"`
	Date        time.Time `toml:"date"`
	Tags        []string  `toml:"tags"`
	IsDraft     bool      `toml:"draft"`
}

type ContentListing struct {
	All   []Article
	ByTag map[string][]Article
}
