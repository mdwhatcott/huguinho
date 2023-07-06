package core

import (
	"errors"
	"testing"
	"time"

	"github.com/mdwhatcott/testing/should"

	"github.com/mdwhatcott/huguinho/contracts"
)

func TestHomePageRenderingHandlerFixture(t *testing.T) {
	should.Run(&HomePageRenderingHandlerFixture{T: should.New(t)}, should.Options.UnitTests())
}

type HomePageRenderingHandlerFixture struct {
	*should.T

	now      time.Time
	handler  *HomePageRenderingHandler
	disk     *InMemoryFileSystem
	renderer *FakeRenderer
}

func (this *HomePageRenderingHandlerFixture) Now() time.Time {
	return this.now
}
func (this *HomePageRenderingHandlerFixture) Setup() {
	this.now = Date(2021, 2, 1)
	this.disk = NewInMemoryFileSystem()
	this.renderer = NewFakeRenderer()
	this.handler = NewHomePageRenderingHandler(this.Now, this.disk, this.renderer, "output/folder")
	this.handleArticles()
}
func (this *HomePageRenderingHandlerFixture) handleArticles() {
	this.handler.Handle(&contracts.Article{
		Metadata: contracts.ArticleMetadata{
			Slug:   "/slug1",
			Title:  "title1",
			Intro:  "intro1",
			Date:   Date(2020, 1, 1),
			Topics: []string{"topic-a", "topic-b"},
		},
	})
	this.handler.Handle(&contracts.Article{
		Metadata: contracts.ArticleMetadata{
			Slug:   "/slug2",
			Title:  "title2",
			Intro:  "intro2",
			Date:   Date(2020, 2, 2),
			Topics: []string{"topic-b", "topic-c"},
			Draft:  true,
		},
	})
	this.handler.Handle(&contracts.Article{
		Metadata: contracts.ArticleMetadata{
			Slug:   "/slug3",
			Title:  "title3",
			Intro:  "intro3",
			Date:   Date(2020, 3, 3),
			Topics: []string{"topic-c", "topic-d"},
		},
	})
}

func (this *HomePageRenderingHandlerFixture) assertHandledArticlesRendered() {
	this.So(this.renderer.rendered, should.Equal, contracts.RenderedHomePage{
		Pages: []contracts.RenderedArticleSummary{
			{
				Slug:   "/slug3",
				Title:  "title3",
				Intro:  "intro3",
				Date:   Date(2020, 3, 3),
				Topics: []string{"topic-c", "topic-d"},
				Draft:  false,
			},
			{
				Slug:   "/slug2",
				Title:  "title2",
				Intro:  "intro2",
				Date:   Date(2020, 2, 2),
				Topics: []string{"topic-b", "topic-c"},
				Draft:  true,
			},
			// TOO OLD:
			//{
			//	Slug:   "/slug1",
			//	Title:  "title1",
			//	Intro:  "intro1",
			//	Date:   Date(2020, 1, 1),
			//	Topics: []string{"topic-a", "topic-b"},
			//	Draft:  false,
			//},
		},
	})
}

func (this *HomePageRenderingHandlerFixture) TestFileTemplateRenderedAndWrittenToDisk() {
	this.renderer.result = "RENDERED"

	err := this.handler.Finalize()

	this.So(err, should.BeNil)
	this.assertHandledArticlesRendered()
	this.So(this.disk.Files, should.Contain, "output/folder")
	if this.So(this.disk.Files, should.Contain, "output/folder/index.html") {
		file := this.disk.Files["output/folder/index.html"]
		this.So(file.Content(), should.Equal, "RENDERED")
	}
}

func (this *HomePageRenderingHandlerFixture) TestRenderErrorReturned() {
	renderErr := errors.New("boink")
	this.renderer.err = renderErr

	err := this.handler.Finalize()

	this.So(err, should.WrapError, renderErr)
	this.So(this.disk.Files, should.BeEmpty)
}

func (this *HomePageRenderingHandlerFixture) TestMkdirAllErrorReturned() {
	this.renderer.result = "RENDERED"
	mkdirErr := errors.New("boink")
	this.disk.ErrMkdirAll["output/folder"] = mkdirErr

	err := this.handler.Finalize()

	this.So(err, should.WrapError, mkdirErr)
	this.So(this.disk.Files, should.BeEmpty)
}

func (this *HomePageRenderingHandlerFixture) TestWriteFileErrorReturned() {
	this.renderer.result = "RENDERED"
	writeFileErr := errors.New("boink")
	this.disk.ErrWriteFile["output/folder/index.html"] = writeFileErr

	err := this.handler.Finalize()

	this.So(err, should.WrapError, writeFileErr)
	this.So(this.disk.Files, should.NOT.Contain, "output/folder/index.html")
}
