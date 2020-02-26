package core

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mdwhatcott/huguinho/contracts"
)

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

func (this *TemplateLoader) Load() (templates *template.Template, err error) {
	err = this.disk.Walk(this.folder, func(path string, info os.FileInfo, err error) error {
		if path != this.folder && info.IsDir() {
			return filepath.SkipDir
		}
		if templates == nil {
			templates = template.New(info.Name())
		} else {
			templates = templates.New(info.Name())
		}
		if !strings.HasSuffix(info.Name(), ".tmpl") {
			return nil
		}
		all, err := this.disk.ReadFile(path)
		if err != nil {
			return StackTraceError(err)
		}
		templates, err = templates.Parse(string(all))
		if err != nil {
			return StackTraceError(err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return templates, nil
}
