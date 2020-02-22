package shell

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"text/template"

	"github.com/mdwhatcott/huguinho/contracts"
)

func ParseTemplates(glob string) *template.Template {
	templates, err := template.ParseGlob(glob)
	if err != nil {
		log.Fatalln("Could not parse templates:", err)
	}
	return templates
}

// TEST
type TemplateRenderer struct {
	templates *template.Template
}

func NewTemplateRenderer(templates *template.Template) *TemplateRenderer {
	return &TemplateRenderer{templates: templates}
}

func (this *TemplateRenderer) Render(v interface{}) (string, error) {
	switch v.(type) {

	case contracts.RenderedArticle:
		return this.render(contracts.ArticleTemplateName, v)

	case contracts.RenderedTopicsListing:
		return this.render(contracts.TopicsTemplateName, v)

	case contracts.RenderedHomePage:
		return this.render(contracts.HomePageTemplateName, v)

	default:
		return "", fmt.Errorf("unsupported type [%v]: %v", reflect.TypeOf(v).Name(), v)
	}
}

func (this *TemplateRenderer) render(name string, data interface{}) (string, error) {
	buffer := new(bytes.Buffer)
	err := this.templates.ExecuteTemplate(buffer, name, data)
	if err != nil {
		return "", fmt.Errorf(
			"failed to render template [%s] (err: %w) with data of type [%v]: %+v",
			name, err, reflect.TypeOf(data), data,
		)
	}
	return buffer.String(), nil
}
