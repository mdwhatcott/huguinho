package core

import "github.com/mdwhatcott/huguinho/contracts"

type FileReadingHandler struct {
	disk contracts.ReadFile
}

func NewFileReadingHandler(disk contracts.ReadFile) *FileReadingHandler {
	return &FileReadingHandler{disk: disk}
}

func (this *FileReadingHandler) Handle(article *contracts.Article) error {
	raw, err := this.disk.ReadFile(article.Source.Path)
	if err != nil {
		return contracts.NewStackTraceError(err)
	}
	article.Source.Data = string(raw)
	return nil
}
