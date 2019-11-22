package main

import (
	"fmt"

	"github.com/mdwhatcott/static/fs"
	"github.com/mdwhatcott/static/site"
)

func main() {
	content := fs.LoadContent("/Users/mike/src/github.com/mdwhatcott/blog/content")
	pages := site.ParsePages(content)
	for _, page := range pages {
		fmt.Println(page.Date.Format("2006-01-02"), page.Path, page.Title)
	}
}
