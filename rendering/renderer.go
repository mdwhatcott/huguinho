package rendering

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"reflect"

	"github.com/mdwhatcott/huguinho/contracts"
)

type Renderer struct {
	templates *template.Template
}

func NewRenderer2(templates *template.Template) *Renderer {
	for _, t := range templates.Templates() {
		log.Println("Loaded template:", t.Name())
	}
	return &Renderer{templates: templates}
}

func (this *Renderer) Render(v interface{}) (string, error) {
	switch v.(type) {

	case contracts.RenderedArticle:
		return this.render2("article.tmpl", v)

	case contracts.RenderedAllTagsListing:
		return this.render2("all-tags.tmpl", v)

	case contracts.RenderedTagListing:
		return this.render2("tag.tmpl", v)

	case contracts.RenderedHomePage:
		return this.render2("home.tmpl", v)

	default:
		return "", fmt.Errorf(
			"unknown value for rendering of type [%v]: %v",
			reflect.TypeOf(v).Name(), v,
		)
	}
}

func (this *Renderer) render2(name string, data interface{}) (string, error) {
	buffer := new(bytes.Buffer)
	err := this.templates.ExecuteTemplate(buffer, name, data)
	if err != nil {
		log.Printf("Failed to render template [%s] (err: %v) with data: %+v", name, err, data)
		return "", fmt.Errorf("failed to render template [%s] (err: %w) with data: %+v", name, err, data)
	}
	return buffer.String(), nil
}

func NewRenderer(dir string) *Renderer {
	templates, err := template.ParseGlob(dir + "/*.html")
	if err != nil {
		log.Fatal(err)
	}
	return &Renderer{templates: templates}
}

func (this *Renderer) RenderHomePage(articles []contracts.Article__DEPRECATED) []byte {
	return this.render("index.html", Listing{Pages: articles})
}

func (this *Renderer) RenderListing(name string, articles []contracts.Article__DEPRECATED) []byte {
	return this.render("tag.html", Listing{
		Name:  name,
		Title: name,
		Path:  "/" + name,
		Pages: articles,
	})
}

func (this *Renderer) RenderPage(article contracts.Article__DEPRECATED) []byte {
	return this.render("page.html", article)
}

func (this *Renderer) render(name string, data interface{}) []byte {
	buffer := new(bytes.Buffer)
	err := this.templates.ExecuteTemplate(buffer, name, data)
	if err != nil {
		log.Printf("Failed to render template [%s] (err: %v) with data: %+v", name, err, data)
		return nil
	}
	return buffer.Bytes()
}

type Listing struct {
	Path        string
	Name        string
	Title       string
	Description string
	Pages       []contracts.Article__DEPRECATED
}
