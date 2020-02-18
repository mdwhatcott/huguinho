package contracts

import "time"

//////////////////////////////////////////////

type Renderer interface {
	Render(interface{}) (string, error)
}

const (
	HomePageTemplateName = "home.tmpl"
	ArticleTemplateName  = "article.tmpl"
	AllTagsTemplateName  = "all-tags.tmpl"
	TagTemplateName      = "tag.tmpl"
)

//////////////////////////////////////////////

type RenderedHomePage struct {
	Path        string // TODO: populate?
	Title       string // TODO: populate?
	Description string // TODO: populate
	Pages       []RenderedHomePageEntry
}

type RenderedHomePageEntry struct {
	Path        string
	Title       string
	Description string
	Date        time.Time
}

//////////////////////////////////////////////

type RenderedArticle struct {
	Path        string // TODO: populate from slug
	Title       string
	Description string // TODO: rename to Intro
	Date        time.Time
	Tags        []string
	Content     string
}

//////////////////////////////////////////////

type RenderedTagListing struct {
	Path        string // TODO: populate from tag
	Title       string
	Name        string
	Description string // TODO: populate?
	Pages       []RenderedTagEntry
}

type RenderedTagEntry struct {
	Path  string
	Title string
	Date  time.Time
}

//////////////////////////////////////////////

type RenderedAllTagsListing struct {
	Path string // TODO: populate
	Tags []RenderedAllTagsEntry
}

type RenderedAllTagsEntry struct {
	Name  string
	Path  string
	Count int
}

//////////////////////////////////////////////
