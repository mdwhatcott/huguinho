package core

import (
	"os"
	"strings"

	"github.com/mdwhatcott/huguinho/contracts"
)

type PathLoader struct {
	files  contracts.Walk
	root   string
	output chan contracts.Article
}

func NewPathLoader(
	files contracts.Walk,
	root string,
	output chan contracts.Article,
) *PathLoader {
	return &PathLoader{
		files:  files,
		root:   root,
		output: output,
	}
}

func (this *PathLoader) Start() error {
	defer close(this.output)
	return this.files.Walk(this.root, this.walk)
}

func (this *PathLoader) walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		return NewStackTraceError(err)
	}
	if info.IsDir() {
		return nil
	}
	if !strings.HasSuffix(path, ".md") {
		return nil
	}
	this.output <- contracts.Article{
		Source: contracts.ArticleSource{Path: path},
	}
	return nil
}
