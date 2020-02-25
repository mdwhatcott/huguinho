package core

import (
	"errors"
	"testing"

	"github.com/mdwhatcott/huguinho/fs"
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/huguinho/contracts"
)

func TestTemplateLoaderFixture(t *testing.T) {
	gunit.Run(new(TemplateLoaderFixture), t)
}

type TemplateLoaderFixture struct {
	*gunit.Fixture

	disk   *fs.InMemoryFileSystem
	loader *TemplateLoader
}

func (this *TemplateLoaderFixture) Setup() {
	this.disk = fs.NewInMemoryFileSystem()
	this.loader = NewTemplateLoader(this.disk, "templates")
	_ = this.disk.MkdirAll("templates", 0755)
	_ = this.disk.WriteFile("templates/supplemental-template.tmpl", []byte(""), 0644)
	_ = this.disk.WriteFile("templates/"+contracts.HomePageTemplateName, []byte(ValidHomePageTemplate), 0644)
	_ = this.disk.WriteFile("templates/"+contracts.TopicsTemplateName, []byte(ValidTopicsPageTemplate), 0644)
	_ = this.disk.WriteFile("templates/"+contracts.ArticleTemplateName, []byte(ValidArticlePageTemplate), 0644)
}

func (this *TemplateLoaderFixture) TestLoadsEachTemplate() {
	templates, err := this.loader.Load()
	this.So(err, should.BeNil)
	this.So(templates.Lookup("supplemental-template.tmpl"), should.NotBeNil)
	this.So(templates.Lookup(contracts.HomePageTemplateName), should.NotBeNil)
	this.So(templates.Lookup(contracts.TopicsTemplateName), should.NotBeNil)
	this.So(templates.Lookup(contracts.ArticleTemplateName), should.NotBeNil)
}

func (this *TemplateLoaderFixture) TestNonTemplateFilesIgnored() {
	_ = this.disk.WriteFile("templates/not-a-template", []byte(ValidHomePageTemplate), 0644)
	templates, err := this.loader.Load()
	this.So(err, should.BeNil)
	this.So(templates.Lookup("not-a-template"), should.BeNil)
	this.So(templates.Lookup(contracts.HomePageTemplateName), should.NotBeNil)
	this.So(templates.Lookup(contracts.TopicsTemplateName), should.NotBeNil)
	this.So(templates.Lookup(contracts.ArticleTemplateName), should.NotBeNil)
}

func (this *TemplateLoaderFixture) TestInvalidTemplateFiles_Error() {
	_ = this.disk.WriteFile("templates/invalid-template.tmpl", []byte("{{ .invalid {{{}{{{})"), 0644)
	templates, err := this.loader.Load()
	this.So(err, should.NotBeNil)
	this.So(templates, should.BeNil)
}

func (this *TemplateLoaderFixture) TestReadFileErr() {
	gophers := errors.New("GOPHERS")
	this.disk.ErrReadFile["templates/"+contracts.HomePageTemplateName] = gophers
	templates, err := this.loader.Load()
	this.So(errors.Is(err, gophers), should.BeTrue)
	this.So(templates, should.BeNil)
}

func (this *TemplateLoaderFixture) TestNestedDirectoriesSkipped() {
	_ = this.disk.MkdirAll("templates/nested", 0755)
	_ = this.disk.WriteFile("templates/nested/template.tmpl", []byte(""), 0644)
	templates, err := this.loader.Load()
	this.So(err, should.BeNil)
	this.So(templates.Lookup("nested/template.tmpl"), should.BeNil)
	this.So(templates.Lookup("template.tmpl"), should.BeNil)
}

const ValidHomePageTemplate = ``
const ValidTopicsPageTemplate = ``
const ValidArticlePageTemplate = ``
