package core

import (
	"io"
	"io/ioutil"
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestInMemoryFilesSystemFixture(t *testing.T) {
	gunit.Run(new(InMemoryFilesSystemFixture), t)
}

type InMemoryFilesSystemFixture struct {
	*gunit.Fixture
	fs *InMemoryFileSystem
}

func (this *InMemoryFilesSystemFixture) Setup() {
	this.fs = NewInMemoryFileSystem()
	this.fs.ModTime = time.Now()
}

const helloWorld = "Hello, world!"

var helloWorldBytes = []byte(helloWorld)

const rootPath = "/path"

func (this *InMemoryFilesSystemFixture) Test_WriteFile_ReadFile() {
	err1 := this.fs.WriteFile(rootPath, helloWorldBytes, 0644)
	this.So(err1, should.BeNil)

	raw, err2 := this.fs.ReadFile(rootPath)
	this.So(err2, should.BeNil)
	this.So(raw, should.Resemble, helloWorldBytes)
}

func (this *InMemoryFilesSystemFixture) Test_CreateFile_OpenFile() {
	writer, err1 := this.fs.Create(rootPath)
	this.So(err1, should.BeNil)

	written, err2 := io.WriteString(writer, helloWorld)
	this.So(written, should.Equal, len(helloWorld))
	this.So(err2, should.BeNil)

	err3 := writer.Close()
	this.So(err3, should.BeNil)

	reader, err4 := this.fs.Open(rootPath)
	this.So(err4, should.BeNil)

	all, err5 := ioutil.ReadAll(reader)
	this.So(err5, should.BeNil)
	this.So(all, should.Resemble, helloWorldBytes)
}

func (this *InMemoryFilesSystemFixture) Test_Stat() {
	_ = this.fs.WriteFile(rootPath, helloWorldBytes, 0644)
	info, err := this.fs.Stat(rootPath)
	this.So(err, should.BeNil)
	this.So(info.Name(), should.Equal, "path")
	this.So(info.IsDir(), should.BeFalse)
	this.So(info.Mode(), should.Equal, 0644)
	this.So(info.ModTime(), should.Equal, this.fs.ModTime)
	this.So(info.Size(), should.Equal, len(helloWorld))
	this.So(func() { _ = info.Sys() }, should.Panic)
}
