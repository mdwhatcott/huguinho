package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mdwhatcott/static/content"
	"github.com/mdwhatcott/static/contracts"
	"github.com/mdwhatcott/static/fs"
)

const (
	root = "/Users/mike/src/github.com/mdwhatcott/blog"
	src  = root + "/content"
	dest = "./rendered"
	base = "https://michaelwhatcott.com"
)

func main() {
	_, thisFile, _, _ := runtime.Caller(0)
	templatesGlob := filepath.Join(filepath.Dir(thisFile), "..", "..", "templates") + "/*.html"
	templates, err := template.ParseGlob(templatesGlob)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(templates.DefinedTemplates())

	_ = os.RemoveAll(dest)

	site := content.ParseAll(fs.LoadFiles(src))

	for _, article := range site[contracts.HomePageListingID] {
		writeFile(article.TargetPath(dest), Page{
			BaseURL: base,
			Article: article,
		})
	}

	for tag, articles := range site {
		if tag == contracts.HomePageListingID {
			writeFile(filepath.Join(dest, "index.html"), Listing{
				BaseURL: base,
				Pages:   articles,
			})
		} else {
			writeFile(filepath.Join(dest, tag, "index.html"), Listing{
				Name:    tag,
				BaseURL: base,
				Pages:   articles,
			})
		}
	}
}

func writeFile(path string, data interface{}) {
	stuff, _ := json.MarshalIndent(data, "", "  ")
	_ = os.MkdirAll(filepath.Dir(path), 0755)
	_ = ioutil.WriteFile(path, stuff, 0644)
}

type Page struct {
	BaseURL string
	contracts.Article
}
type Listing struct {
	Name    string
	BaseURL string
	Pages   []contracts.Article
}
