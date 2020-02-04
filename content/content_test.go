package content

import (
	"sort"
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/huguinho/contracts"
)

func TestArticleFixture(t *testing.T) {
	gunit.Run(new(ArticleFixture), t)
}

type ArticleFixture struct {
	*gunit.Fixture
}

func (this *ArticleFixture) TestParseEmptyFileToPage() {
	file := contracts.File("")
	page := parse(file)
	this.So(page, should.Resemble, contracts.Article__DEPRECATED{})
}

func (this *ArticleFixture) TestParseContentOnlyFileToPage() {
	file := contracts.File("I have some content")
	page := parse(file)
	this.So(page, should.Resemble, contracts.Article__DEPRECATED{
		OriginalContent: "I have some content",
		Content:         "<p>I have some content</p>\n",
	})
}

func (this *ArticleFixture) TestParseEmptyFrontMatterAndContentToPage() {
	file := contracts.File("+++\n\n+++\nI have some content")
	page := parse(file)
	this.So(page, should.Resemble, contracts.Article__DEPRECATED{
		OriginalContent: "I have some content",
		Content:         "<p>I have some content</p>\n",
	})
}

func (this *ArticleFixture) TestParseFrontMatterAndContentToPage() {
	file := contracts.File(`+++
title = "The Title"
description = "The Intro"
date = 2019-11-21
tags = ["a", "b", "c"]
draft = true
+++

The Content
`)
	page := parse(file)
	this.So(page, should.Resemble, contracts.Article__DEPRECATED{
		FrontMatter__DEPRECATED: contracts.FrontMatter__DEPRECATED{
			Title:       "The Title",
			Description: "The Intro",
			Date:        time.Date(2019, 11, 21, 0, 0, 0, 0, time.Local),
			Tags:        []string{"a", "b", "c"},
			IsDraft:     true,
		},
		OriginalContent: "The Content",
		Content:         "<p>The Content</p>\n",
	})
}

func (this *ArticleFixture) TestParseFrontMatterMalformed() {
	file := contracts.File(`+++
I am not front matter at all.
+++

The Content
`)
	page := parse(file)
	this.So(page.ParseError, should.NotBeNil)
	this.Println(page.ParseError)
}

func (this *ArticleFixture) TestDerivePath() {
	files := map[contracts.Path]contracts.File{
		"/a/b/c": "Hello",
		"/1/2/3": "World",
	}

	pages := parseAll(files)
	sort.Slice(pages, func(i, j int) bool {
		return pages[i].OriginalContent < pages[j].OriginalContent
	})

	this.So(pages, should.Resemble, []contracts.Article__DEPRECATED{
		{
			Path:            "/a/b/c/",
			OriginalContent: "Hello",
			Content:         "<p>Hello</p>\n",
		},
		{
			Path:            "/1/2/3/",
			OriginalContent: "World",
			Content:         "<p>World</p>\n",
		},
	})
}
