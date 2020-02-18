package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mdwhatcott/huguinho/content"
	"github.com/mdwhatcott/huguinho/contracts"
	"github.com/mdwhatcott/huguinho/fs"
	"github.com/mdwhatcott/huguinho/rendering"
)

func main() {
	start := time.Now()
	config := ParseCLI()
	listing, err := ioutil.ReadDir(config.targetRoot)
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range listing {
		_ = os.RemoveAll(filepath.Join(config.targetRoot, item.Name()))
	}
	renderer := rendering.NewRenderer(config.templateDir)
	site := content.ParseAll(fs.LoadFiles(config.contentRoot), config.buildDrafts, config.buildFuture)
	renderArticles(config.targetRoot, renderer, site)
	renderListings(config.targetRoot, renderer, site)
	includeCSS(config.targetRoot, config.stylesDir)
	log.Println("Duration:", time.Since(start))
}

func renderArticles(root string, renderer *rendering.Renderer, site contracts.Site__DEPRECATED) {
	for _, article := range site[contracts.HomePageListingID__DEPRECATED] {
		fs.WriteFile(article.TargetPath(root), renderer.RenderPage(article))
	}
}

func renderListings(root string, renderer *rendering.Renderer, site contracts.Site__DEPRECATED) {
	for tag, articles := range site {
		if tag == contracts.HomePageListingID__DEPRECATED {
			fs.WriteFile(filepath.Join(root, "index.html"), renderer.RenderHomePage(articles))
		} else {
			fs.WriteFile(filepath.Join(root, "tags", tag, "index.html"), renderer.RenderListing(tag, articles))
		}
	}
}

func includeCSS(root, stylesDir string) {
	listing, err := ioutil.ReadDir(stylesDir)
	if err != nil {
		log.Println(err)
	}
	if len(listing) == 0 {
		return
	}
	for _, file := range listing {
		name := file.Name()
		if strings.HasSuffix(name, ".css") {
			path := filepath.Join(stylesDir, name)
			data := fs.ReadFile(path)
			fs.WriteFile(filepath.Join(root, "css", name), data)
		}
	}
}
