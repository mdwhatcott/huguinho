package main

import (
	"fmt"
	"os"

	"github.com/mdwhatcott/static/content"
	"github.com/mdwhatcott/static/fs"
)

const contentRoot = "/Users/mike/src/github.com/mdwhatcott/blog/content"

func main() {
	root := "./rendered"
	_ = os.Mkdir(root, 0755)

	files := fs.LoadFiles(contentRoot)
	listing := content.ParseAll(files)

	fmt.Println("--", "ALL", "--")
	for _, article := range listing.All {
		// populate index file
		fmt.Println(article.Date.Format("2006-01-02"), article.Path, article.Title, article.Description)
	}

	for tag, articles := range listing.ByTag {
		fmt.Println("--", tag, "--")
		for _, article := range articles {
			// populate tag listing file
			fmt.Println(article.Date.Format("2006-01-02"), article.Path, article.Title, article.Description)
		}
	}
}
