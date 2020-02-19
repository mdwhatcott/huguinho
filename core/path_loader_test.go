package core

import (
	"errors"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/huguinho/contracts"
)

func TestPathLoaderFixture(t *testing.T) {
	gunit.Run(new(PathLoaderFixture), t)
}

type PathLoaderFixture struct {
	*gunit.Fixture
	loader *PathLoader
	files  *InMemoryFileSystem
	output chan contracts.Article
}

func (this *PathLoaderFixture) Setup() {
	this.files = NewInMemoryFileSystem()
	this.output = make(chan contracts.Article, 10)
	this.loader = NewPathLoader(this.files, "/content", this.output)

	_ = this.files.WriteFile("/article1.md", []byte("outside of content root"), 0644)
	_ = this.files.WriteFile("/content/article1.md", []byte("article1"), 0644)
	_ = this.files.WriteFile("/content/article2.txt", []byte("not an article"), 0644)
	_ = this.files.Mkdir("/content/folder", 0577)
	_ = this.files.WriteFile("/content/folder/article3.md", []byte("article3"), 0644)
}

func (this *PathLoaderFixture) Test() {
	this.loader.Start()
	err := this.loader.Finalize()

	this.So(err, should.BeNil)
	this.So(gather(this.output), should.Resemble, []contracts.Article{
		{Source: contracts.ArticleSource{Path: "/content/article1.md"}},
		{Source: contracts.ArticleSource{Path: "/content/folder/article3.md"}},
	})
}

func (this *PathLoaderFixture) TestErrWalkFunc() {
	this.files.ErrWalkFunc["/content/folder/article3.md"] = walkFuncErr

	this.loader.Start()
	err := this.loader.Finalize()

	this.So(errors.Is(err, walkFuncErr), should.BeTrue)
	this.So(gather(this.output), should.Resemble, []contracts.Article{
		{Source: contracts.ArticleSource{Path: "/content/article1.md"}},
	})
}

var walkFuncErr = errors.New("walk func error")
