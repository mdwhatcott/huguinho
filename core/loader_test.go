package core

import (
	"errors"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestContentLoaderFixture(t *testing.T) {
	gunit.Run(new(ContentLoaderFixture), t)
}

type ContentLoaderFixture struct {
	*gunit.Fixture

	loader *ContentLoader
	files  *InMemoryFileSystem
}

func (this *ContentLoaderFixture) Setup() {
	this.files = NewInMemoryFileSystem()
	this.loader = NewContentLoader(this.files, "/content")
}

func (this *ContentLoaderFixture) TestLoadContent() {
	_ = this.files.WriteFile("/article1.md", []byte("outside of content root"), 0644)
	_ = this.files.WriteFile("/content/article1.md", []byte("article1"), 0644)
	_ = this.files.WriteFile("/content/article2.txt", []byte("not an article"), 0644)
	_ = this.files.Mkdir("/content/folder", 0577)
	_ = this.files.WriteFile("/content/folder/article3.md", []byte("article3"), 0644)

	content, err := this.loader.LoadContent()

	this.So(err, should.BeNil)
	this.So(content, should.Resemble, []string{
		"article1",
		"article3",
	})
}

func (this *ContentLoaderFixture) TestLoadContent_Err() {
	_ = this.files.WriteFile("/content/article1.md", []byte("article1"), 0644)
	_ = this.files.WriteFile("/content/article2.md", []byte("article2"), 0644)
	this.files.ErrReadFile["/content/article2.md"] = errors.New("gophers")

	content, err := this.loader.LoadContent()

	this.So(err, should.NotBeNil)
	this.So(content, should.Resemble, []string{"article1"})
}
