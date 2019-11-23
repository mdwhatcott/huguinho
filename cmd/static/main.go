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
	"github.com/mdwhatcott/static/fs"
)

const (
	root      = "/Users/mike/src/github.com/mdwhatcott/blog"
	src       = root + "/content"
	dest      = "./rendered"
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

	listing := content.ParseAll(fs.LoadFiles(src))

	writeFile(filepath.Join(dest, "index.html"), listing.All)

	for _, article := range listing.All {
		writeFile(article.TargetPath(dest), article)
	}

	for tag, articles := range listing.ByTag {
		writeFile(filepath.Join(dest, tag, "index.html"), articles)
	}
}

func writeFile(path string, data interface{}) {
	stuff, _ := json.MarshalIndent(data, "", "  ")
	_ = os.MkdirAll(filepath.Dir(path), 0755)
	_ = ioutil.WriteFile(path, stuff, 0644)
}
