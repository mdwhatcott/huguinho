package main

import (
	"fmt"

	"github.com/mdwhatcott/static/fs"
	"github.com/mdwhatcott/static/pages"
)

func main() {
	content := fs.LoadContent("/Users/mike/src/github.com/mdwhatcott/blog/content")
	parsed := pages.ParseAll(content)
	listings := pages.OrganizePages(parsed)

	fmt.Println("--", "ALL", "--")
	for _, post := range listings.All {
		fmt.Println(post.Date.Format("2006-01-02"), post.Path, post.Title, post.Description)
	}

	for tag, posts := range listings.ByTag {
		fmt.Println("--", tag, "--")
		for _, page := range posts {
			fmt.Println(page.Date.Format("2006-01-02"), page.Path, page.Title, page.Description)
		}
	}
}
