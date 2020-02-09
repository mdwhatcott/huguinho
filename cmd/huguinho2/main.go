package main

import (
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/mdwhatcott/huguinho/contracts"
	"github.com/mdwhatcott/huguinho/core"
	"github.com/mdwhatcott/huguinho/rendering"
	"github.com/mdwhatcott/huguinho/shell"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config := ParseCLI()
	renderer := rendering.NewRenderer2(parseTemplates(filepath.Join(config.templateDir, "*.tmpl")))
	pipeline := NewPipeline(config, shell.NewDisk(), renderer)
	os.Exit(pipeline.Run())
}
func parseTemplates(glob string) *template.Template {
	templates, err := template.ParseGlob(glob)
	if err != nil {
		log.Fatal("Could not parse templates:", err)
	}
	return templates
}

type Pipeline struct {
	config   *Config
	disk     contracts.FileSystem
	renderer contracts.Renderer
	errs     chan error
}

func NewPipeline(config *Config, disk contracts.FileSystem, renderer contracts.Renderer) *Pipeline {
	return &Pipeline{
		config:   config,
		disk:     disk,
		renderer: renderer,
		errs:     make(chan error),
	}
}
func (this *Pipeline) Run() int {
	// TODO: Move CSS into place
	out := this.startAll()
	go this.terminate(out)
	return this.errCount()
}
func (this *Pipeline) startAll() (out chan contracts.Article) {
	out = this.goLoad()
	out = this.goListen(out, core.NewFileReader(this.disk))
	out = this.goListen(out, core.NewJSONMetadataParser())
	// out = this.goListen(out, core.NewDraftFiltering(...)) // TODO
	// out = this.goListen(out, core.NewFutureFiltering(...)) // TODO
	out = this.goListen(out, core.NewContentParser(shell.NewGoldmarkMarkdownConverter()))
	// out = this.goListen(out, core.NewHomePageRenderer(...)) // TODO
	out = this.goListen(out, core.NewArticleRenderingHandler(this.disk, this.renderer, this.config.targetRoot))
	// out = this.goListen(out, core.NewTagsRenderer(...)) // TODO
	// out = this.goListen(out, core.NewAllTagsRenderer(...)) // TODO
	return out
}
func (this *Pipeline) goLoad() (out chan contracts.Article) {
	out = make(chan contracts.Article)
	loader := core.NewPathLoader(this.disk, this.config.contentRoot, out)
	go func() { this.errs <- loader.Start() }()
	return out
}
func (this *Pipeline) goListen(in chan contracts.Article, handler contracts.Handler) (out chan contracts.Article) {
	out = make(chan contracts.Article)
	listener := core.NewListener(in, out, handler)
	go func() { this.errs <- listener.Listen() }()
	return out
}
func (this *Pipeline) terminate(out chan contracts.Article) {
	for item := range out {
		log.Println("Published article:", item.Metadata.Slug)
	}
	close(this.errs)
}
func (this *Pipeline) errCount() (errCount int) {
	for err := range this.errs {
		if err != nil {
			errCount++
			log.Println("[WARN]", err)
		}
	}
	return errCount
}
