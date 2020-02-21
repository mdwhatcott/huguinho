package contracts

import "time"

type Renderer interface {
	Render(interface{}) (string, error)
}

const (
	HomePageTemplateName = "home.tmpl"
	ArticleTemplateName  = "article.tmpl"
	TopicsTemplateName   = "topics.tmpl"
)

type (
	RenderedHomePage struct {
		Pages []RenderedArticleSummary
	}

	RenderedArticle struct {
		Slug    string
		Title   string
		Intro   string
		Date    time.Time
		Topics  []string
		Content string
	}

	RenderedArticleSummary struct {
		Slug  string
		Title string
		Intro string
		Date  time.Time
	}

	RenderedTopicsListing struct {
		Topics []RenderedTopicListing
	}

	RenderedTopicListing struct {
		Topic    string
		Articles []RenderedArticleSummary
	}
)
