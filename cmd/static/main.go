package main

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/mdwhatcott/static/content"
	"github.com/mdwhatcott/static/contracts"
	"github.com/mdwhatcott/static/fs"
	"github.com/mdwhatcott/static/rendering"
)

const (
	root = "/Users/mike/src/github.com/mdwhatcott/blog"
	src  = root + "/content"
	dest = "./rendered"
	base = "https://michaelwhatcott.com"
)

func main() {
	_ = os.RemoveAll(dest)
	renderer := buildRenderer()
	site := content.ParseAll(fs.LoadFiles(src))
	renderArticles(site, renderer)
	renderListings(site, renderer)
}

func buildRenderer() *rendering.Renderer {
	_, thisFile, _, _ := runtime.Caller(0)
	templatesGlob := filepath.Join(filepath.Dir(thisFile), "..", "..", "templates") + "/*.html"
	return rendering.NewRenderer(base, templatesGlob)
}

func renderArticles(site contracts.Site, renderer *rendering.Renderer) {
	for _, article := range site[contracts.HomePageListingID] {
		fs.WriteFile(article.TargetPath(dest), renderer.RenderPage(article))
	}
}

func renderListings(site contracts.Site, renderer *rendering.Renderer) {
	for tag, articles := range site {
		if tag == contracts.HomePageListingID {
			fs.WriteFile(filepath.Join(dest, "index.html"), renderer.RenderHomePage(articles))
		} else {
			fs.WriteFile(filepath.Join(dest, "tags", tag, "index.html"), renderer.RenderListing(tag, articles))
		}
	}
}
