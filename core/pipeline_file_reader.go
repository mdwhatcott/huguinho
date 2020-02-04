package core

import "github.com/mdwhatcott/huguinho/contracts"

type FileReader struct {
	disk   contracts.ReadFile
	input  chan contracts.Article
	output chan contracts.Article
	err    error
}

func NewFileReader(
	disk contracts.ReadFile,
	input chan contracts.Article,
	output chan contracts.Article,
) *FileReader {
	return &FileReader{
		disk:   disk,
		input:  input,
		output: output,
	}
}

func (this *FileReader) Listen() {
	defer close(this.output)
	defer this.drain()

	for article := range this.input {
		raw, err := this.disk.ReadFile(article.Source.Path)
		if err != nil {
			this.err = NewStackTraceError(err)
			return
		}
		article.Source.Data = string(raw)
		this.output <- article
	}
}

func (this *FileReader) drain() {
	for range this.input {
	}
}

func (this *FileReader) Finalize() error {
	return this.err
}
