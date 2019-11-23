package main

import (
	"flag"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mdwhatcott/huguinho/content"
	"github.com/mdwhatcott/huguinho/contracts"
	"github.com/mdwhatcott/huguinho/fs"
	"github.com/mdwhatcott/huguinho/rendering"
)

func main() {
	config := parseCLI()
	_ = os.RemoveAll(config.dest)
	renderer := buildRenderer()
	site := content.ParseAll(fs.LoadFiles(config.src), config.drafts, config.future)
	renderArticles(config.dest, renderer, site)
	renderListings(config.dest, renderer, site)
	includeCSS(config.dest)
}

type Config struct {
	src    string
	dest   string
	drafts bool
	future bool
}

func parseCLI() (config Config) {
	flag.StringVar(&config.src, "src", "", "The source directory (required)")
	flag.StringVar(&config.dest, "dest", "", "The destination directory (required)")
	flag.BoolVar(&config.drafts, "drafts", false, "When set, draft articles will be included in rendered output.")
	flag.BoolVar(&config.future, "future", false, "When set, future articles will be included in rendered output.")
	flag.Parse()
	if config.src == "" || config.dest == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	return config
}

func buildRenderer() *rendering.Renderer {
	_, thisFile, _, _ := runtime.Caller(0)
	templatesGlob := filepath.Join(filepath.Dir(thisFile), "..", "..", "templates") + "/*.html"
	return rendering.NewRenderer(templatesGlob)
}

func renderArticles(dest string, renderer *rendering.Renderer, site contracts.Site) {
	for _, article := range site[contracts.HomePageListingID] {
		fs.WriteFile(article.TargetPath(dest), renderer.RenderPage(article))
	}
}

func renderListings(dest string, renderer *rendering.Renderer, site contracts.Site) {
	for tag, articles := range site {
		if tag == contracts.HomePageListingID {
			fs.WriteFile(filepath.Join(dest, "index.html"), renderer.RenderHomePage(articles))
		} else {
			fs.WriteFile(filepath.Join(dest, "tags", tag, "index.html"), renderer.RenderListing(tag, articles))
		}
	}
}

func includeCSS(dest string) {
	_, thisFile, _, _ := runtime.Caller(0)
	cssFile := filepath.Join(filepath.Dir(thisFile), "..", "..", "css", "custom.css")
	data := fs.ReadFile(cssFile)
	fs.WriteFile(filepath.Join(dest, "css", "custom.css"), data)
}
