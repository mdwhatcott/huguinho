package core

import (
	"errors"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/huguinho/contracts"
)

func TestArticleRenderingHandlerFixture(t *testing.T) {
	gunit.Run(new(ArticleRenderingHandlerFixture), t)
}

type ArticleRenderingHandlerFixture struct {
	*gunit.Fixture
	handler  *ArticleRenderingHandler
	renderer *FakeRenderer
	disk     *InMemoryFileSystem
	article  *contracts.Article
}

func (this *ArticleRenderingHandlerFixture) Setup() {
	this.renderer = NewFakeRenderer()
	this.disk = NewInMemoryFileSystem()
	this.handler = NewArticleRenderingHandler(this.disk, this.renderer, "output/folder")

	this.article = &contracts.Article{
		Metadata: contracts.ArticleMetadata{
			Draft: false,
			Slug:  "/slug",
			Title: "Title",
			Intro: "Intro",
			Tags:  []string{"A", "B"},
			Date:  Date(2020, 2, 8),
		},
		Content: contracts.ArticleContent{
			Converted: "CONTENT",
		},
	}
}

func (this *ArticleRenderingHandlerFixture) TestFileTemplateRenderedAndWrittenToDisk() {
	this.renderer.result = "RENDERED"

	err := this.handler.Handle(this.article)

	this.So(err, should.BeNil)
	this.So(this.renderer.rendered, should.Resemble, this.article)
	this.So(this.disk.Files, should.ContainKey, "output/folder/slug")
	if this.So(this.disk.Files, should.ContainKey, "output/folder/slug/index.html") {
		file := this.disk.Files["output/folder/slug/index.html"]
		this.So(file.content, should.Resemble, []byte("RENDERED"))
	}
}

func (this *ArticleRenderingHandlerFixture) TestRenderErrorReturned() {
	renderErr := errors.New("boink")
	this.renderer.err = renderErr

	err := this.handler.Handle(this.article)

	this.So(errors.Is(err, renderErr), should.BeTrue)
	this.So(this.renderer.rendered, should.Resemble, this.article)
	this.So(this.disk.Files, should.BeEmpty)
}

func (this *ArticleRenderingHandlerFixture) TestMkdirAllErrorReturned() {
	this.renderer.result = "RENDERED"
	mkdirErr := errors.New("boink")
	this.disk.ErrMkdirAll["output/folder/slug"] = mkdirErr

	err := this.handler.Handle(this.article)

	this.So(errors.Is(err, mkdirErr), should.BeTrue)
	this.So(this.renderer.rendered, should.Resemble, this.article)
	this.So(this.disk.Files, should.BeEmpty)
}

func (this *ArticleRenderingHandlerFixture) TestWriteFileErrorReturned() {
	this.renderer.result = "RENDERED"
	writeFileErr := errors.New("boink")
	this.disk.ErrWriteFile["output/folder/slug/index.html"] = writeFileErr

	err := this.handler.Handle(this.article)

	this.So(errors.Is(err, writeFileErr), should.BeTrue)
	this.So(this.renderer.rendered, should.Resemble, this.article)
	this.So(this.disk.Files, should.NotContainKey, "output/folder/slug/index.html")
}

///////////////////////////////////////////////////////////////////

type FakeRenderer struct {
	rendered interface{}
	result   string
	err      error
}

func NewFakeRenderer() *FakeRenderer {
	return &FakeRenderer{}
}

func (this *FakeRenderer) Render(rendered interface{}) (string, error) {
	this.rendered = rendered
	return this.result, this.err
}
