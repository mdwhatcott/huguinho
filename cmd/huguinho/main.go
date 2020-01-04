package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mdwhatcott/huguinho/content"
	"github.com/mdwhatcott/huguinho/contracts"
	"github.com/mdwhatcott/huguinho/fs"
	"github.com/mdwhatcott/huguinho/rendering"
)

func main() {
	config := parseCLI()
	_ = os.RemoveAll(config.outputRoot)
	renderer := buildRenderer()
	site := content.ParseAll(fs.LoadFiles(config.contentRoot), config.buildDrafts, config.buildFuture)
	renderArticles(config.outputRoot, renderer, site)
	renderListings(config.outputRoot, renderer, site)
	includeCSS(config.outputRoot)
}

type Config struct {
	contentRoot string
	outputRoot  string
	buildDrafts bool
	buildFuture bool
}

func parseCLI() (config Config) {
	flag.StringVar(&config.contentRoot, "src", "", "The source directory (required)")
	flag.StringVar(&config.outputRoot, "dest", "", "The destination directory (required)")
	flag.BoolVar(&config.buildDrafts, "drafts", false, "When set, draft articles will be included in rendered output.")
	flag.BoolVar(&config.buildFuture, "future", false, "When set, future articles will be included in rendered output.")
	flag.Parse()
	if config.contentRoot == "" || config.outputRoot == "" {
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

func renderArticles(root string, renderer *rendering.Renderer, site contracts.Site) {
	for _, article := range site[contracts.HomePageListingID] {
		fs.WriteFile(article.TargetPath(root), renderer.RenderPage(article))
	}
}

func renderListings(root string, renderer *rendering.Renderer, site contracts.Site) {
	for tag, articles := range site {
		if tag == contracts.HomePageListingID {
			fs.WriteFile(filepath.Join(root, "index.html"), renderer.RenderHomePage(articles))
		} else {
			fs.WriteFile(filepath.Join(root, "tags", tag, "index.html"), renderer.RenderListing(tag, articles))
		}
	}
}

func includeCSS(root string) {
	_, thisFile, _, _ := runtime.Caller(0)
	cssFolder := filepath.Join(filepath.Dir(thisFile), "..", "..", "css")
	listing, err := ioutil.ReadDir(cssFolder)
	if err != nil {
		log.Println(err)
	}
	if len(listing) == 0 {
		return
	}
	for _, file := range listing {
		name := file.Name()
		if strings.HasSuffix(name, ".css") {
			path := filepath.Join(cssFolder, name)
			data := fs.ReadFile(path)
			fs.WriteFile(filepath.Join(root, "css", name), data)
		}
	}
}
