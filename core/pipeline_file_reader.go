package core

import "github.com/mdwhatcott/huguinho/contracts"

type FileReader struct {
	disk contracts.ReadFile
}

func NewFileReader(disk contracts.ReadFile) *FileReader {
	return &FileReader{disk: disk}
}

func (this *FileReader) Handle(article *contracts.Article) error {
	raw, err := this.disk.ReadFile(article.Source.Path)
	if err != nil {
		return NewStackTraceError(err)
	}
	article.Source.Data = string(raw)
	return nil
}
