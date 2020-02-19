package core

import (
	"errors"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/huguinho/contracts"
)

func TestTagPageRenderingHandlerFixture(t *testing.T) {
	gunit.Run(new(TagPageRenderingHandlerFixture), t)
}

type TagPageRenderingHandlerFixture struct {
	*gunit.Fixture

	handler  *TagPageRenderingHandler
	disk     *InMemoryFileSystem
	renderer *FakeRenderer
}

func (this *TagPageRenderingHandlerFixture) Setup() {
	this.disk = NewInMemoryFileSystem()
	this.renderer = NewFakeRenderer()
	this.handler = NewTagPageRenderingHandler(this.disk, this.renderer, "output/folder")
	this.handleArticles()
}

func (this *TagPageRenderingHandlerFixture) handleArticles() {
	this.So(this.handler.Handle(&contracts.Article{
		Metadata: contracts.ArticleMetadata{
			Slug:  "/slug1",
			Title: "title1",
			Intro: "intro1",
			Date:  Date(2020, 1, 1),
			Tags:  []string{"a", "b"},
		},
	}), should.BeNil)
	this.So(this.handler.Handle(&contracts.Article{
		Metadata: contracts.ArticleMetadata{
			Slug:  "/slug2",
			Title: "title2",
			Intro: "intro2",
			Date:  Date(2020, 2, 2),
			Tags:  []string{"b", "c"},
		},
	}), should.BeNil)
	this.So(this.handler.Handle(&contracts.Article{
		Metadata: contracts.ArticleMetadata{
			Slug:  "/slug3",
			Title: "title3",
			Intro: "intro3",
			Date:  Date(2020, 3, 3),
			Tags:  []string{"c"},
		},
	}), should.BeNil)
}

func (this *TagPageRenderingHandlerFixture) assertHandledArticlesRenderedInTagListings() {
	if !this.So(this.renderer.all, should.HaveLength, 3) {
		return
	}
	this.So(this.renderer.all, should.Contain, contracts.RenderedTagListing{
		Title: "a",
		Name:  "a",
		Pages: []contracts.RenderedTagEntry{
			{
				Slug:  "/slug1",
				Title: "title1",
				Date:  Date(2020, 1, 1),
			},
		},
	})
	this.So(this.renderer.all, should.Contain, contracts.RenderedTagListing{
		Title: "b",
		Name:  "b",
		Pages: []contracts.RenderedTagEntry{
			{
				Slug:  "/slug2",
				Title: "title2",
				Date:  Date(2020, 2, 2),
			},
			{
				Slug:  "/slug1",
				Title: "title1",
				Date:  Date(2020, 1, 1),
			},
		},
	})
	this.So(this.renderer.all, should.Contain, contracts.RenderedTagListing{
		Title: "c",
		Name:  "c",
		Pages: []contracts.RenderedTagEntry{
			{
				Slug:  "/slug3",
				Title: "title3",
				Date:  Date(2020, 3, 3),
			},
			{
				Slug:  "/slug2",
				Title: "title2",
				Date:  Date(2020, 2, 2),
			},
		},
	})
}

func (this *TagPageRenderingHandlerFixture) TestTagTemplateRenderedAndWrittenToDisk() {
	this.renderer.result = "RENDERED"

	err := this.handler.Finalize()

	this.So(err, should.BeNil)
	this.assertHandledArticlesRenderedInTagListings()
	this.So(this.disk.Files, should.ContainKey, "output/folder/tags/a")
	this.So(this.disk.Files, should.ContainKey, "output/folder/tags/b")
	this.So(this.disk.Files, should.ContainKey, "output/folder/tags/c")
	if this.So(this.disk.Files, should.ContainKey, "output/folder/tags/a/index.html") {
		file := this.disk.Files["output/folder/tags/a/index.html"]
		this.So(file.content, should.Resemble, []byte("RENDERED"))
	}
	if this.So(this.disk.Files, should.ContainKey, "output/folder/tags/b/index.html") {
		file := this.disk.Files["output/folder/tags/b/index.html"]
		this.So(file.content, should.Resemble, []byte("RENDERED"))
	}
	if this.So(this.disk.Files, should.ContainKey, "output/folder/tags/c/index.html") {
		file := this.disk.Files["output/folder/tags/c/index.html"]
		this.So(file.content, should.Resemble, []byte("RENDERED"))
	}
}

func (this *TagPageRenderingHandlerFixture) TestRenderErrorReturned() {
	renderErr := errors.New("boink")
	this.renderer.err = renderErr

	err := this.handler.Finalize()

	this.So(errors.Is(err, renderErr), should.BeTrue)
	this.So(this.disk.Files, should.BeEmpty)
}

func (this *TagPageRenderingHandlerFixture) TestMkdirAllErrorReturned() {
	this.renderer.result = "RENDERED"
	mkdirErr := errors.New("boink")
	this.disk.ErrMkdirAll["output/folder/tags/a"] = mkdirErr

	err := this.handler.Finalize()

	this.So(errors.Is(err, mkdirErr), should.BeTrue)
	this.So(this.disk.Files, should.NotContainKey, "output/folder/tags/a")
}

func (this *TagPageRenderingHandlerFixture) TestWriteFileErrorReturned() {
	this.renderer.result = "RENDERED"
	writeFileErr := errors.New("boink")
	this.disk.ErrWriteFile["output/folder/tags/a/index.html"] = writeFileErr

	err := this.handler.Finalize()

	this.So(errors.Is(err, writeFileErr), should.BeTrue)
	this.So(this.disk.Files, should.NotContainKey, "output/folder/tags/a/index.html")
}
