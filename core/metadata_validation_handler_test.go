package core

import (
	"errors"
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/huguinho/contracts"
)

func TestMetadataValidationHandlerFixture(t *testing.T) {
	gunit.Run(new(MetadataValidationHandlerFixture), t)
}

type MetadataValidationHandlerFixture struct {
	*gunit.Fixture

	handler *MetadataValidationHandler
	article *contracts.Article
}

func (this *MetadataValidationHandlerFixture) Setup() {
	this.handler = NewMetadataValidationHandler()
	this.article = &contracts.Article{
		Metadata: contracts.ArticleMetadata{
			Draft: false,
			Slug:  "/slug1",
			Title: "Title",
			Intro: "Introduction",
			Tags:  []string{"a", "b", "c"},
			Date:  Date(2020, 2, 2),
		},
	}
}

func (this *MetadataValidationHandlerFixture) TestAllPresentAndAccountedFor() {
	err := this.handler.Handle(this.article)
	this.So(err, should.BeNil)
}
func (this *MetadataValidationHandlerFixture) TestMissingTitle_Err() {
	this.article.Metadata.Title = ""
	err := this.handler.Handle(this.article)
	this.So(errors.Is(err, errBlankMetadataTitle), should.BeTrue)
}
func (this *MetadataValidationHandlerFixture) TestMissingSlug_Err() {
	this.article.Metadata.Slug = ""
	err := this.handler.Handle(this.article)
	this.So(errors.Is(err, errBlankMetadataSlug), should.BeTrue)
}
func (this *MetadataValidationHandlerFixture) TestMissingDate_Err() {
	this.article.Metadata.Date = time.Time{}
	err := this.handler.Handle(this.article)
	this.So(errors.Is(err, errBlankMetadataDate), should.BeTrue)
}
func (this *MetadataValidationHandlerFixture) TestUniqueSlugs_OK() {
	this.assertHandleWithSlugOK("a")
	this.assertHandleWithSlugOK("b")
	this.assertHandleWithSlugOK("c")
}
func (this *MetadataValidationHandlerFixture) TestRepeatedSlugs_Err() {
	this.assertHandleWithSlugOK("A")
	this.assertHandleWithSlugOK("b")
	this.assertHandleWithSlugOK("c")
	this.assertHandleWithSlugFAIL("A")
}
func (this *MetadataValidationHandlerFixture) assertHandleWithSlugOK(slug string) {
	this.article.Metadata.Slug = slug
	err := this.handler.Handle(this.article)
	this.So(err, should.BeNil)
}
func (this *MetadataValidationHandlerFixture) assertHandleWithSlugFAIL(slug string) {
	this.article.Metadata.Slug = slug
	err := this.handler.Handle(this.article)
	this.So(errors.Is(err, errRepeatedMetadataSlug), should.BeTrue)
}
