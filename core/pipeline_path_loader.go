package core

import (
	"os"
	"strings"

	"github.com/mdwhatcott/huguinho/contracts"
)

type PathLoader struct {
	files  contracts.Walk
	root   string
	output chan contracts.Page
	err    error
}

func NewPathLoader(
	files contracts.Walk,
	root string,
	output chan contracts.Page,
) *PathLoader {

	return &PathLoader{
		files:  files,
		root:   root,
		output: output,
	}
}

func (this *PathLoader) Listen() {
	this.err = this.files.Walk(this.root, this.walk)
	close(this.output)
}

func (this *PathLoader) walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	if !strings.HasSuffix(path, ".md") {
		return nil
	}
	this.output <- contracts.Page{SourcePath: path}
	return nil
}

func (this *PathLoader) Finalize() error {
	return this.err
}
