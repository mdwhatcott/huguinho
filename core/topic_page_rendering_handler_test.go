package core

import (
	"errors"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/huguinho/contracts"
)

func TestTopicPageRenderingHandlerFixture(t *testing.T) {
	gunit.Run(new(TopicPageRenderingHandlerFixture), t)
}

type TopicPageRenderingHandlerFixture struct {
	*gunit.Fixture

	handler  *TopicPageRenderingHandler
	disk     *InMemoryFileSystem
	renderer *FakeRenderer
}

func (this *TopicPageRenderingHandlerFixture) Setup() {
	this.disk = NewInMemoryFileSystem()
	this.renderer = NewFakeRenderer()
	this.handler = NewTopicPageRenderingHandler(this.disk, this.renderer, "output/folder")
	this.handleArticles()
}

func (this *TopicPageRenderingHandlerFixture) handleArticles() {
	this.handler.Handle(&contracts.Article{
		Metadata: contracts.ArticleMetadata{
			Slug:   "/slug1",
			Title:  "title1",
			Intro:  "intro1",
			Date:   Date(2020, 1, 1),
			Topics: []string{"a", "b"},
		},
	})
	this.handler.Handle(&contracts.Article{
		Metadata: contracts.ArticleMetadata{
			Slug:   "/slug2",
			Title:  "title2",
			Intro:  "intro2",
			Date:   Date(2020, 2, 2),
			Topics: []string{"b", "c"},
		},
	})
	this.handler.Handle(&contracts.Article{
		Metadata: contracts.ArticleMetadata{
			Slug:   "/slug3",
			Title:  "title3",
			Intro:  "intro3",
			Date:   Date(2020, 3, 3),
			Topics: []string{"c"},
		},
	})
}

func (this *TopicPageRenderingHandlerFixture) assertHandledArticlesRendered() bool {
	return this.So(this.renderer.rendered, should.Resemble, contracts.RenderedTopicsListing{
		Topics: []contracts.RenderedTopicListing{
			{
				Topic: "a",
				Articles: []contracts.RenderedArticleSummary{
					{
						Slug:  "/slug1",
						Title: "title1",
						Intro: "intro1",
						Date:  Date(2020, 1, 1),
					},
				},
			},
			{
				Topic: "b",
				Articles: []contracts.RenderedArticleSummary{
					{
						Slug:  "/slug2",
						Title: "title2",
						Intro: "intro2",
						Date:  Date(2020, 2, 2),
					},
					{
						Slug:  "/slug1",
						Title: "title1",
						Intro: "intro1",
						Date:  Date(2020, 1, 1),
					},
				},
			},
			{
				Topic: "c",
				Articles: []contracts.RenderedArticleSummary{
					{
						Slug:  "/slug3",
						Title: "title3",
						Intro: "intro3",
						Date:  Date(2020, 3, 3),
					},
					{
						Slug:  "/slug2",
						Title: "title2",
						Intro: "intro2",
						Date:  Date(2020, 2, 2),
					},
				},
			},
		},
	})
}

func (this *TopicPageRenderingHandlerFixture) TestFileTemplateRenderedAndWrittenToDisk() {
	this.renderer.result = "RENDERED"

	err := this.handler.Finalize()

	this.So(err, should.BeNil)
	this.assertHandledArticlesRendered()
	this.So(this.disk.Files, should.ContainKey, "output/folder")
	if this.So(this.disk.Files, should.ContainKey, "output/folder/topics/index.html") {
		file := this.disk.Files["output/folder/topics/index.html"]
		this.So(file.content, should.Resemble, []byte("RENDERED"))
	}
}

func (this *TopicPageRenderingHandlerFixture) TestRenderErrorReturned() {
	renderErr := errors.New("boink")
	this.renderer.err = renderErr

	err := this.handler.Finalize()

	this.So(errors.Is(err, renderErr), should.BeTrue)
	this.So(this.disk.Files, should.BeEmpty)
}

func (this *TopicPageRenderingHandlerFixture) TestMkdirAllErrorReturned() {
	this.renderer.result = "RENDERED"
	mkdirErr := errors.New("boink")
	this.disk.ErrMkdirAll["output/folder/topics"] = mkdirErr

	err := this.handler.Finalize()

	this.So(errors.Is(err, mkdirErr), should.BeTrue)
	this.So(this.disk.Files, should.BeEmpty)
}

func (this *TopicPageRenderingHandlerFixture) TestWriteFileErrorReturned() {
	this.renderer.result = "RENDERED"
	writeFileErr := errors.New("boink")
	this.disk.ErrWriteFile["output/folder/topics/index.html"] = writeFileErr

	err := this.handler.Finalize()

	this.So(errors.Is(err, writeFileErr), should.BeTrue)
	this.So(this.disk.Files, should.NotContainKey, "output/folder/topics/index.html")
}
