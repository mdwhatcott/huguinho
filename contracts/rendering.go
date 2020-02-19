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
	Pages []RenderedHomePageEntry
}

type RenderedHomePageEntry struct {
	Slug  string
	Title string
	Intro string
	Date  time.Time
}

//////////////////////////////////////////////

type RenderedArticle struct {
	Slug    string
	Title   string
	Intro   string
	Date    time.Time
	Tags    []string
	Content string
}

//////////////////////////////////////////////

type RenderedTagListing struct {
	Title string
	Name  string
	Pages []RenderedTagEntry
}

type RenderedTagEntry struct {
	Slug  string
	Title string
	Date  time.Time
}

//////////////////////////////////////////////

type RenderedAllTagsListing struct {
	Tags []RenderedAllTagsEntry
}

type RenderedAllTagsEntry struct {
	Name  string
	Path  string
	Count int
}

//////////////////////////////////////////////
