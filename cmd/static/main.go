package main

import (
	"fmt"
	"os"

	"github.com/mdwhatcott/static/fs"
	"github.com/mdwhatcott/static/pages"
)

func main() {
	root := "./rendered"
	_ = os.Mkdir(root, 0755)

	content := fs.LoadFiles("/Users/mike/src/github.com/mdwhatcott/blog/content")
	parsed := pages.ParseAll(content)
	listings := pages.OrganizeContent(parsed)

	fmt.Println("--", "ALL", "--")
	for _, post := range listings.All {
		// populate index file
		fmt.Println(post.Date.Format("2006-01-02"), post.Path, post.Title, post.Description)
	}

	for tag, posts := range listings.ByTag {
		fmt.Println("--", tag, "--")
		for _, page := range posts {
			// populate tag listing file
			fmt.Println(page.Date.Format("2006-01-02"), page.Path, page.Title, page.Description)
		}
	}
}
