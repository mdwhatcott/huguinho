package rendering

import (
	"bytes"
	"html/template"
	"log"

	"github.com/mdwhatcott/static/contracts"
)

type Renderer struct {
	baseURL   string
	templates *template.Template
}

func NewRenderer(baseURL, glob string) *Renderer {
	templates, err := template.ParseGlob(glob)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(templates.DefinedTemplates())

	return &Renderer{
		baseURL:   baseURL,
		templates: templates,
	}
}

func (this *Renderer) RenderHomePage(articles []contracts.Article) []byte {
	return this.execute("index.html", Listing{
		BaseURL: this.baseURL,
		Pages:   articles,
	})
}

func (this *Renderer) RenderListing(name string, articles []contracts.Article) []byte {
	return this.execute("tag.html", Listing{
		Name:    name,
		Path:    "/"+name,
		BaseURL: this.baseURL,
		Pages:   articles,
	})
}

func (this *Renderer) RenderPage(article contracts.Article) []byte {
	return this.execute("page.html", Page{
		BaseURL: this.baseURL,
		Article: article,
	})
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

type Page struct {
	BaseURL string
	contracts.Article
}

type Listing struct {
	Path        string
	Name        string
	Title       string
	Description string
	BaseURL     string
	Pages       []contracts.Article
}
