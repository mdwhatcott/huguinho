package pages

import (
	"sort"
	"testing"
	"time"

	"github.com/mdwhatcott/static/contracts"
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestPageFixture(t *testing.T) {
	gunit.Run(new(PageFixture), t)
}

type PageFixture struct {
	*gunit.Fixture
}

func (this *PageFixture) TestParseEmptyFileToPage() {
	file := contracts.File("")
	page := Parse(file)
	this.So(page, should.Resemble, contracts.Page{})
}

func (this *PageFixture) TestParseContentOnlyFileToPage() {
	file := contracts.File("I have some content")
	page := Parse(file)
	this.So(page, should.Resemble, contracts.Page{
		OriginalContent: "I have some content",
		HTMLContent:     "<p>I have some content</p>\n",
	})
}

func (this *PageFixture) TestParseEmptyFrontMatterAndContentToPage() {
	file := contracts.File("+++\n\n+++\nI have some content")
	page := Parse(file)
	this.So(page, should.Resemble, contracts.Page{
		OriginalContent: "I have some content",
		HTMLContent:     "<p>I have some content</p>\n",
	})
}

func (this *PageFixture) TestParseFrontMatterAndContentToPage() {
	file := contracts.File(`+++
title = "The Title"
description = "The Description"
date = 2019-11-21
tags = ["a", "b", "c"]
draft = true
+++

The Content
`)
	page := Parse(file)
	this.So(page, should.Resemble, contracts.Page{
		FrontMatter: contracts.FrontMatter{
			Title:       "The Title",
			Description: "The Description",
			Date:        time.Date(2019, 11, 21, 0, 0, 0, 0, time.Local),
			Tags:        []string{"a", "b", "c"},
			IsDraft:     true,
		},
		OriginalContent: "The Content",
		HTMLContent:     "<p>The Content</p>\n",
	})
}

func (this *PageFixture) TestParseFrontMatterMalformed() {
	file := contracts.File(`+++
I am not front matter at all.
+++

The Content
`)
	page := Parse(file)
	this.So(page.ParseError, should.NotBeNil)
	this.Println(page.ParseError)
}

func (this *PageFixture) TestDeriveSlug() {
	files := map[contracts.Path]contracts.File{
		"/a/b/c": "Hello",
		"/1/2/3": "World",
	}

	pages := ParseAll(files)
	sort.Slice(pages, func(i, j int) bool {
		return pages[i].OriginalContent < pages[j].OriginalContent
	})

	this.So(pages, should.Resemble, []contracts.Page{
		{
			Path:            "/a/b/c",
			OriginalContent: "Hello",
			HTMLContent:     "<p>Hello</p>\n",
		},
		{
			Path:            "/1/2/3",
			OriginalContent: "World",
			HTMLContent:     "<p>World</p>\n",
		},
	})
}
