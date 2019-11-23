package rendering

import (
	"bytes"
	"log"
	"text/template"

	"github.com/mdwhatcott/static/contracts"
)

type Renderer struct {
	templates *template.Template
}

func NewRenderer(glob string) *Renderer {
	templates, err := template.ParseGlob(glob)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(templates.DefinedTemplates())

	return &Renderer{
		templates: templates,
	}
}

func (this *Renderer) RenderHomePage(articles []contracts.Article) []byte {
	return this.execute("index.html", Listing{
		Pages:   articles,
	})
}

func (this *Renderer) RenderListing(name string, articles []contracts.Article) []byte {
	return this.execute("tag.html", Listing{
		Name:    name,
		Title:   name,
		Path:    "/" + name,
		Pages:   articles,
	})
}

func (this *Renderer) RenderPage(article contracts.Article) []byte {
	return this.execute("page.html", article)
}

func (this *Renderer) execute(name string, data interface{}) []byte {
	buffer := new(bytes.Buffer)
	err := this.templates.ExecuteTemplate(buffer, name, data)
	if err != nil {
		log.Printf("Failed to execute template [%s] (err: %v) with data: %+v", name, err, data)
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
