package main

import (
	"flag"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mdwhatcott/static/content"
	"github.com/mdwhatcott/static/contracts"
	"github.com/mdwhatcott/static/fs"
	"github.com/mdwhatcott/static/rendering"
)

var (
	src    string
	dest   string
	drafts bool
	future bool
)

func main() {
	parseCLI()
	_ = os.RemoveAll(dest)
	renderer := buildRenderer()
	site := content.ParseAll(fs.LoadFiles(src), drafts, future)
	renderArticles(renderer, site)
	renderListings(renderer, site)
	includeCSS()
}

func parseCLI() {
	flag.StringVar(&src, "src", "", "The source directory (required)")
	flag.StringVar(&dest, "dest", "", "The destination directory (required)")
	flag.BoolVar(&drafts, "drafts", false, "When set, draft articles will be included in rendered output.")
	flag.BoolVar(&future, "future", false, "When set, future articles will be included in rendered output.")
	flag.Parse()
	if src == "" || dest == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func buildRenderer() *rendering.Renderer {
	_, thisFile, _, _ := runtime.Caller(0)
	templatesGlob := filepath.Join(filepath.Dir(thisFile), "..", "..", "templates") + "/*.html"
	return rendering.NewRenderer(templatesGlob)
}

func renderArticles(renderer *rendering.Renderer, site contracts.Site) {
	for _, article := range site[contracts.HomePageListingID] {
		fs.WriteFile(article.TargetPath(dest), renderer.RenderPage(article))
	}
}

func renderListings(renderer *rendering.Renderer, site contracts.Site) {
	for tag, articles := range site {
		if tag == contracts.HomePageListingID {
			fs.WriteFile(filepath.Join(dest, "index.html"), renderer.RenderHomePage(articles))
		} else {
			fs.WriteFile(filepath.Join(dest, "tags", tag, "index.html"), renderer.RenderListing(tag, articles))
		}
	}
}

func includeCSS() {
	_, thisFile, _, _ := runtime.Caller(0)
	cssFile := filepath.Join(filepath.Dir(thisFile), "..", "..", "css", "custom.css")
	data := fs.ReadFile(cssFile)
	fs.WriteFile(filepath.Join(dest, "css", "custom.css"), data)
}
