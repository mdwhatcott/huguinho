package core

import (
	"testing"

	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/huguinho/contracts"
)

func TestTemplateLoaderFixture(t *testing.T) {
	gunit.Run(new(TemplateLoaderFixture), t)
}

type TemplateLoaderFixture struct {
	*gunit.Fixture

	disk   *InMemoryFileSystem
	loader *TemplateLoader
}

func (this *TemplateLoaderFixture) Setup() {
	this.disk = NewInMemoryFileSystem()
	this.loader = NewTemplateLoader(this.disk, "templates")
	_ = this.disk.WriteFile("templates/"+contracts.HomePageTemplateName, []byte(ValidHomePageTemplate), 0644)
	_ = this.disk.WriteFile("templates/"+contracts.TopicsTemplateName, []byte(ValidTopicsPageTemplate), 0644)
	_ = this.disk.WriteFile("templates/"+contracts.ArticleTemplateName, []byte(ValidArticlePageTemplate), 0644)
}

func (this *TemplateLoaderFixture) Test() {
	// TODO
}

const ValidHomePageTemplate = ``
const ValidTopicsPageTemplate = ``
const ValidArticlePageTemplate = ``
