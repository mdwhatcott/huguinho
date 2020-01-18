package rendering

import (
	"bytes"
	"log"
	"text/template"

	"github.com/mdwhatcott/huguinho/contracts"
)

type Renderer struct {
	templates *template.Template
}

func NewRenderer(dir string) *Renderer {
	templates, err := template.ParseGlob(dir + "/*.html")
	if err != nil {
		log.Fatal(err)
	}
	return &Renderer{templates: templates}
}

func (this *Renderer) RenderHomePage(articles []contracts.Article) []byte {
	return this.render("index.html", Listing{Pages: articles})
}

func (this *Renderer) RenderListing(name string, articles []contracts.Article) []byte {
	return this.render("tag.html", Listing{
		Name:  name,
		Title: name,
		Path:  "/" + name,
		Pages: articles,
	})
}

func (this *Renderer) RenderPage(article contracts.Article) []byte {
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
	Pages       []contracts.Article
}
