package contracts

import (
	"path/filepath"
	"strings"
	"time"
)

type Article struct {
	FrontMatter
	Path            Path   `json:"Path"`
	OriginalContent string `json:"-"`
	Content         string `json:"Content"`
}

func (this Article) TargetPath(root string) string {
	folder := strings.TrimSuffix(string(this.Path), ".md")
	return filepath.Join(root, folder, "index.html")
}

type FrontMatter struct {
	ParseError  error     `json:"-"           toml:"-"`
	Title       string    `json:"Title"       toml:"title"`
	Description string    `json:"Description" toml:"description"`
	Date        time.Time `json:"Date"        toml:"date"`
	Tags        []string  `json:"Tags"        toml:"tags"`
	IsDraft     bool      `json:"-"           toml:"draft"`
}

type Site map[string][]Article

type ContentListing struct {
	Pages []Article
	ByTag map[string][]Article
}

const HomePageListingID = "00000000-ba05-4d31-97b3-57d8d80b0dda"
