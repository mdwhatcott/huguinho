package core

import (
	"errors"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/huguinho/contracts"
)

func TestHomePageRenderingHandlerFixture(t *testing.T) {
	gunit.Run(new(HomePageRenderingHandlerFixture), t)
}

type HomePageRenderingHandlerFixture struct {
	*gunit.Fixture

	handler  *HomePageRenderingHandler
	disk     *InMemoryFileSystem
	renderer *FakeRenderer
}

func (this *HomePageRenderingHandlerFixture) Setup() {
	this.disk = NewInMemoryFileSystem()
	this.renderer = NewFakeRenderer()
	this.handler = NewHomePageRenderingHandler(this.disk, this.renderer, "output/folder")
	this.handleArticles()
}

func (this *HomePageRenderingHandlerFixture) handleArticles() {
	this.So(this.handler.Handle(&contracts.Article{
		Metadata: contracts.ArticleMetadata{
			Slug:  "/slug1",
			Title: "title1",
			Intro: "intro1",
			Date:  Date(2020, 1, 1),
		},
	}), should.BeNil)
	this.So(this.handler.Handle(&contracts.Article{
		Metadata: contracts.ArticleMetadata{
			Slug:  "/slug2",
			Title: "title2",
			Intro: "intro2",
			Date:  Date(2020, 2, 2),
		},
	}), should.BeNil)
	this.So(this.handler.Handle(&contracts.Article{
		Metadata: contracts.ArticleMetadata{
			Slug:  "/slug3",
			Title: "title3",
			Intro: "intro3",
			Date:  Date(2020, 3, 3),
		},
	}), should.BeNil)
}

func (this *HomePageRenderingHandlerFixture) assertHandledArticlesRendered() bool {
	return this.So(this.renderer.rendered, should.Resemble, contracts.RenderedHomePage{
		Pages: []contracts.RenderedHomePageEntry{
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
			{
				Slug:  "/slug1",
				Title: "title1",
				Intro: "intro1",
				Date:  Date(2020, 1, 1),
			},
		},
	})
}

func (this *HomePageRenderingHandlerFixture) TestFileTemplateRenderedAndWrittenToDisk() {
	this.renderer.result = "RENDERED"

	err := this.handler.Finalize()

	this.So(err, should.BeNil)
	this.assertHandledArticlesRendered()
	this.So(this.disk.Files, should.ContainKey, "output/folder")
	if this.So(this.disk.Files, should.ContainKey, "output/folder/index.html") {
		file := this.disk.Files["output/folder/index.html"]
		this.So(file.content, should.Resemble, []byte("RENDERED"))
	}
}

func (this *HomePageRenderingHandlerFixture) TestRenderErrorReturned() {
	renderErr := errors.New("boink")
	this.renderer.err = renderErr

	err := this.handler.Finalize()

	this.So(errors.Is(err, renderErr), should.BeTrue)
	this.So(this.disk.Files, should.BeEmpty)
}

func (this *HomePageRenderingHandlerFixture) TestMkdirAllErrorReturned() {
	this.renderer.result = "RENDERED"
	mkdirErr := errors.New("boink")
	this.disk.ErrMkdirAll["output/folder"] = mkdirErr

	err := this.handler.Finalize()

	this.So(errors.Is(err, mkdirErr), should.BeTrue)
	this.So(this.disk.Files, should.BeEmpty)
}

func (this *HomePageRenderingHandlerFixture) TestWriteFileErrorReturned() {
	this.renderer.result = "RENDERED"
	writeFileErr := errors.New("boink")
	this.disk.ErrWriteFile["output/folder/index.html"] = writeFileErr

	err := this.handler.Finalize()

	this.So(errors.Is(err, writeFileErr), should.BeTrue)
	this.So(this.disk.Files, should.NotContainKey, "output/folder/index.html")
}
