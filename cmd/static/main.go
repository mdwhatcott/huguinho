package main

import (
	"fmt"
	"strings"

	"github.com/mdwhatcott/static/fs"
)

func main() {
	content := fs.LoadContent("/Users/mike/src/github.com/mdwhatcott/blog/content")
	for path, file := range content {
		fmt.Println(path, strings.ReplaceAll(string(file)[:20], "\n", "\\n"))
	}
}
