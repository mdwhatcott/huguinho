package core

import (
	"os"
	"strings"

	"github.com/mdwhatcott/huguinho/contracts"
)

type ContentLoaderFileSystem interface {
	contracts.Walk
	contracts.ReadFile
}

type ContentLoader struct {
	files ContentLoaderFileSystem
	root  string
}

func NewContentLoader(files ContentLoaderFileSystem, root string) *ContentLoader {
	return &ContentLoader{files: files, root: root}
}

func (this *ContentLoader) LoadContent() (content []ContentFile, err error) {
	err = this.files.Walk(this.root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}
		contents, readErr := this.files.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		content = append(content, ContentFile{
			Path:    path,
			Content: string(contents),
		})
		return nil
	})
	return content, err
}

type ContentFile struct {
	Path    string
	Content string
}
