package contracts

import (
	"path/filepath"
	"strings"
	"time"
)

type Article struct {
	FrontMatter
	Path            Path
	OriginalContent string
	Content         string
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

type Site map[string][]Article

type ContentListing struct {
	Pages []Article
	ByTag map[string][]Article
}

const HomePageListingID = "00000000-ba05-4d31-97b3-57d8d80b0dda"
