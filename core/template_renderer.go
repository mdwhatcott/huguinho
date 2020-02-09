package core

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"reflect"

	"github.com/mdwhatcott/huguinho/contracts"
)

type TemplateRenderer struct {
	templates *template.Template
}

func NewTemplateRenderer(templates *template.Template) *TemplateRenderer {
	return &TemplateRenderer{templates: templates}
}

func (this *TemplateRenderer) Render(v interface{}) (string, error) {
	switch v.(type) {

	case contracts.RenderedArticle:
		return this.render("article.tmpl", v)

	case contracts.RenderedAllTagsListing:
		return this.render("all-tags.tmpl", v)

	case contracts.RenderedTagListing:
		return this.render("tag.tmpl", v)

	case contracts.RenderedHomePage:
		return this.render("home.tmpl", v)

	default:
		return "", fmt.Errorf(
			"unknown value for rendering of type [%v]: %v",
			reflect.TypeOf(v).Name(), v,
		)
	}
}

func (this *TemplateRenderer) render(name string, data interface{}) (string, error) {
	buffer := new(bytes.Buffer)
	err := this.templates.ExecuteTemplate(buffer, name, data)
	if err != nil {
		log.Printf("Failed to render template [%s] (err: %v) with data: %+v", name, err, data)
		return "", fmt.Errorf("failed to render template [%s] (err: %w) with data: %+v", name, err, data)
	}
	return buffer.String(), nil
}
