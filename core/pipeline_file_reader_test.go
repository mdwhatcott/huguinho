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
	input  chan contracts.Article
	output chan contracts.Article
}

func (this *FileReaderFixture) Setup() {
	this.input = make(chan contracts.Article, 10)
	this.output = make(chan contracts.Article, 10)
	this.files = NewInMemoryFileSystem()
	this.reader = NewFileReader(this.files, this.input, this.output)

	_ = this.files.WriteFile("/file1", []byte("FILE1"), 0644)
	_ = this.files.WriteFile("/file2", []byte("FILE2"), 0644)
	_ = this.files.WriteFile("/file3", []byte("FILE3"), 0644)

	this.input <- contracts.Article{Source: contracts.ArticleSource{Path: "/file1"}}
	this.input <- contracts.Article{Source: contracts.ArticleSource{Path: "/file2"}}
	this.input <- contracts.Article{Source: contracts.ArticleSource{Path: "/file3"}}
	close(this.input)
}

func (this *FileReaderFixture) Test() {
	this.reader.Listen()
	err := this.reader.Finalize()

	this.So(err, should.BeNil)
	this.So(gather(this.output), should.Resemble, []contracts.Article{
		{Source: contracts.ArticleSource{Path: "/file1", Data: "FILE1"}},
		{Source: contracts.ArticleSource{Path: "/file2", Data: "FILE2"}},
		{Source: contracts.ArticleSource{Path: "/file3", Data: "FILE3"}},
	})
}

func (this *FileReaderFixture) TestReadError() {
	this.files.ErrReadFile["/file2"] = readError

	this.reader.Listen()
	err := this.reader.Finalize()

	this.So(errors.Is(err, readError), should.BeTrue)
	this.So(gather(this.output), should.Resemble, []contracts.Article{
		{Source: contracts.ArticleSource{Path: "/file1", Data: "FILE1"}},
	})
}

var readError = errors.New("read error")
