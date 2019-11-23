package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mdwhatcott/static/content"
	"github.com/mdwhatcott/static/fs"
)

const src = "/Users/mike/src/github.com/mdwhatcott/blog/content"
const dest = "./rendered"

func main() {
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
