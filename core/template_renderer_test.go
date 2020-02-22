package core

import (
	"errors"
	"testing"
	"text/template"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/huguinho/contracts"
)

func TestTemplateRendererFixture(t *testing.T) {
	gunit.Run(new(TemplateRendererFixture), t)
}

type TemplateRendererFixture struct {
	*gunit.Fixture

	renderer *TemplateRenderer
}

func (this *TemplateRendererFixture) Setup() {
	var err error

	t := template.New(contracts.HomePageTemplateName)
	t, err = t.Parse(contracts.HomePageTemplateName)
	this.So(err, should.Equal, nil)

	t = t.New(contracts.ArticleTemplateName)
	t, err = t.Parse(contracts.ArticleTemplateName)
	this.So(err, should.Equal, nil)

	t = t.New(contracts.TopicsTemplateName)
	t, err = t.Parse(contracts.TopicsTemplateName)
	this.So(err, should.Equal, nil)

	this.renderer = NewTemplateRenderer(t)
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
	this.So(errors.Is(homeErr, contracts.ErrUnsupportedRenderingType), should.BeTrue)
	this.So(home, should.BeBlank)
}

func (this *TemplateRendererFixture) TestRenderError() {
	this.prepareRendererWithBadTemplate()

	rendered, err := this.renderer.Render(contracts.RenderedTopicsListing{})
	this.So(errors.Is(err, contracts.ErrRenderingFailure), should.BeTrue)
	this.So(rendered, should.BeBlank)
}

func (this *TemplateRendererFixture) prepareRendererWithBadTemplate() {
	var err error
	t := template.New(contracts.TopicsTemplateName)
	t, err = t.Parse("{{ .UnknownField }}")
	this.So(err, should.BeNil)

	this.renderer = NewTemplateRenderer(t)
}
