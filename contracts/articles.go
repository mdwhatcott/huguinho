package contracts

import (
	"path/filepath"
	"strings"
	"time"
)

type Article__DEPRECATED struct {
	FrontMatter__DEPRECATED
	Path            Path
	OriginalContent string
	Content         string
}

func (this Article__DEPRECATED) TargetPath(root string) string {
	folder := strings.TrimSuffix(string(this.Path), ".md")
	return filepath.Join(root, folder, "index.html")
}

type FrontMatter__DEPRECATED struct {
	ParseError  error     `toml:"-"`
	Title       string    `toml:"title"`
	Description string    `toml:"description"`
	Date        time.Time `toml:"date"`
	Tags        []string  `toml:"tags"`
	IsDraft     bool      `toml:"draft"`
}

type Site__DEPRECATED map[string][]Article__DEPRECATED

type ContentListing__DEPRECATED struct {
	Pages []Article__DEPRECATED
	ByTag map[string][]Article__DEPRECATED
}

const HomePageListingID__DEPRECATED = "00000000-ba05-4d31-97b3-57d8d80b0dda"
