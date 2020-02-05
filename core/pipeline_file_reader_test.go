package core

import (
	"errors"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/huguinho/contracts"
)

func TestFileReaderFixture(t *testing.T) {
	gunit.Run(new(FileReaderFixture), t)
}

type FileReaderFixture struct {
	*gunit.Fixture
	reader *FileReader
	files  *InMemoryFileSystem
}

func (this *FileReaderFixture) Setup() {
	this.files = NewInMemoryFileSystem()
	this.reader = NewFileReader(this.files)

	_ = this.files.WriteFile("/file1", []byte("FILE1"), 0644)
}

func (this *FileReaderFixture) TestRead() {
	article := &contracts.Article{Source: contracts.ArticleSource{Path: "/file1"}}
	err := this.reader.Handle(article)
	this.So(err, should.BeNil)
	this.So(article, should.Resemble, &contracts.Article{
		Source: contracts.ArticleSource{Path: "/file1", Data: "FILE1"},
	})
}

func (this *FileReaderFixture) TestReadError() {
	article := &contracts.Article{Source: contracts.ArticleSource{Path: "/file1"}}
	this.files.ErrReadFile["/file1"] = readError

	err := this.reader.Handle(article)

	this.So(errors.Is(err, readError), should.BeTrue)
	this.So(article, should.Resemble, &contracts.Article{
		Source: contracts.ArticleSource{Path: "/file1"},
	})
}

var readError = errors.New("read error")
