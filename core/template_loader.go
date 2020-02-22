package core

import "github.com/mdwhatcott/huguinho/contracts"

type TemplateLoaderFileSystem interface {
	contracts.ReadFile
	contracts.Walk
}

type TemplateLoader struct {
	disk   TemplateLoaderFileSystem
	folder string
}

func NewTemplateLoader(disk TemplateLoaderFileSystem, folder string) *TemplateLoader {
	return &TemplateLoader{disk: disk, folder: folder}
}
