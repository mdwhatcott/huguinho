package core

import (
	"testing"
	"text/template"

	"github.com/mdwhatcott/testing/should"

	"github.com/mdwhatcott/huguinho/contracts"
)

func TestTemplateRendererFixture(t *testing.T) {
	should.Run(&TemplateRendererFixture{T: should.New(t)}, should.Options.UnitTests())
}

type TemplateRendererFixture struct {
	*should.T

	templates *template.Template
	renderer  *TemplateRenderer
}

func (this *TemplateRendererFixture) Setup() {
	this.parseHomePageTemplate()
	this.parseArticleTemplate()
	this.parseTopicsTemplate()
	this.renderer = NewTemplateRenderer(this.templates)
	this.So(this.renderer.Validate(), should.BeNil)
}

func (this *TemplateRendererFixture) parseTopicsTemplate() {
	var err error
	if this.templates == nil {
		this.templates = template.New(contracts.TopicsTemplateName)
	} else {
		this.templates = this.templates.New(contracts.TopicsTemplateName)
	}
	this.templates, err = this.templates.Parse(contracts.TopicsTemplateName)
	this.So(err, should.BeNil)
}

func (this *TemplateRendererFixture) parseArticleTemplate() {
	var err error
	if this.templates == nil {
		this.templates = template.New(contracts.ArticleTemplateName)
	} else {
		this.templates = this.templates.New(contracts.ArticleTemplateName)
	}
	this.templates, err = this.templates.Parse(contracts.ArticleTemplateName)
	this.So(err, should.BeNil)
}

func (this *TemplateRendererFixture) parseHomePageTemplate() {
	var err error
	if this.templates == nil {
		this.templates = template.New(contracts.HomePageTemplateName)
	} else {
		this.templates = this.templates.New(contracts.HomePageTemplateName)
	}
	this.templates, err = this.templates.Parse(contracts.HomePageTemplateName)
	this.So(err, should.BeNil)
}

func (this *TemplateRendererFixture) TestMissingHomePageTemplate_ValidateErr() {
	this.templates = nil
	this.parseArticleTemplate()
	this.parseTopicsTemplate()
	this.renderer = NewTemplateRenderer(this.templates)
	this.So(this.renderer.Validate(), should.NOT.BeNil)
}

func (this *TemplateRendererFixture) TestMissingTopicsTemplate_ValidateErr() {
	this.templates = nil
	this.parseArticleTemplate()
	this.parseHomePageTemplate()
	this.renderer = NewTemplateRenderer(this.templates)
	this.So(this.renderer.Validate(), should.NOT.BeNil)
}

func (this *TemplateRendererFixture) TestMissingArticleTemplate_ValidateErr() {
	this.templates = nil
	this.parseHomePageTemplate()
	this.parseTopicsTemplate()
	this.renderer = NewTemplateRenderer(this.templates)
	this.So(this.renderer.Validate(), should.NOT.BeNil)
}

func (this *TemplateRendererFixture) TestCanRenderTypesCorrespondingToTemplates() {
	home, homeErr := this.renderer.Render(contracts.RenderedHomePage{})
	this.So(homeErr, should.BeNil)
	this.So(home, should.Equal, contracts.HomePageTemplateName)

	article, articleErr := this.renderer.Render(contracts.RenderedArticle{})
	this.So(articleErr, should.BeNil)
	this.So(article, should.Equal, contracts.ArticleTemplateName)

	topics, topicsErr := this.renderer.Render(contracts.RenderedTopicsListing{})
	this.So(topicsErr, should.BeNil)
	this.So(topics, should.Equal, contracts.TopicsTemplateName)
}

func (this *TemplateRendererFixture) TestCannotRenderUnknownTypes() {
	home, homeErr := this.renderer.Render(42)
	this.So(homeErr, should.WrapError, contracts.ErrUnsupportedRenderingType)
	this.So(home, should.BeEmpty)
}

func (this *TemplateRendererFixture) TestRenderError() {
	this.prepareRendererWithBadTemplate()

	rendered, err := this.renderer.Render(contracts.RenderedTopicsListing{})
	this.So(err, should.WrapError, contracts.ErrRenderingFailure)
	this.So(rendered, should.BeEmpty)
}

func (this *TemplateRendererFixture) prepareRendererWithBadTemplate() {
	var err error
	t := template.New(contracts.TopicsTemplateName)
	t, err = t.Parse("{{ .UnknownField }}")
	this.So(err, should.BeNil)

	this.renderer = NewTemplateRenderer(t)
}
