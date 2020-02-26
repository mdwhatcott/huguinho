package core

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/clock"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/logging"
)

func TestPipelineRunnerFixture(t *testing.T) {
	gunit.Run(new(PipelineRunnerFixture), t)
}

type PipelineRunnerFixture struct {
	*gunit.Fixture

	started  time.Time
	finished time.Time
	args     []string
	disk     *InMemoryFileSystem
	runner   *PipelineRunner
}

func (this *PipelineRunnerFixture) Setup() {
	this.disk = NewInMemoryFileSystem()
	this.runner = NewPipelineRunner(this.args, this.disk)
	this.runner.log = logging.Capture()
	this.runner.log.SetFlags(0)

	this.started = time.Now()
	this.finished = this.started.Add(time.Millisecond)
	this.runner.clock = clock.Freeze(this.started, this.finished)

	this.file("content/a.md", ContentA)
	this.file("content/b.md", ContentB)
	this.file("content/c.md", ContentC)

	this.file("templates/home.tmpl", TemplateHome)
	this.file("templates/topics.tmpl", TemplateTopics)
	this.file("templates/article.tmpl", TemplateArticle)
}

func (this *PipelineRunnerFixture) arg(values ...string) {
	this.args = append(this.args, values...)
}
func (this *PipelineRunnerFixture) file(path, content string) {
	_ = this.disk.WriteFile(path, []byte(content), 0644)
}
func (this *PipelineRunnerFixture) ls(root string) {
	err := this.disk.Walk(root, func(path string, info os.FileInfo, err error) error {
		this.Println(path)
		return nil
	})

	this.So(err, should.BeNil)
}
func (this *PipelineRunnerFixture) assertFolder(path string) {
	dir := this.disk.Files[path]
	if this.So(dir, should.NotBeNil) {
		this.So(dir.IsDir(), should.BeTrue)
	}
}
func (this *PipelineRunnerFixture) assertFile(path, expectedContent string) {
	file := this.disk.Files[path]
	if this.So(file, should.NotBeNil) {
		actual := strings.TrimSpace(file.Content())
		expected := strings.TrimSpace(expectedContent)
		this.So(actual, should.EqualTrimSpace, expected)
	}
}

func (this *PipelineRunnerFixture) TODOTestBadConfigPreventsProcessing_Error() {
	// TODO
}

func (this *PipelineRunnerFixture) TODOTestTemplateLoadFailurePreventsProcessing_Error() {
	// TODO
}

func (this *PipelineRunnerFixture) TODOTestTemplateValidationFailurePreventsProcessing_Error() {
	// TODO
}

func (this *PipelineRunnerFixture) TestValidConfigAndTemplates_PipelineRuns() {
	this.So(len(this.disk.Files), should.Equal, 6) // 3 articles, 3 templates

	errors := this.runner.Run()

	this.So(errors, should.Equal, 0)
	this.So(len(this.disk.Files), should.Equal, 12) // 3 folders, 3 html files

	this.assertFolder("rendered")
	this.assertFolder("rendered/topics")
	this.assertFolder("rendered/article-a")

	this.assertFile("rendered/index.html", RenderedHome)
	this.assertFile("rendered/topics/index.html", RenderedTopics)
	this.assertFile("rendered/article-a/index.html", RenderedArticleA)
}

const (
	ContentA = `
slug:   /article-a/
title:  Article A
intro:  The introduction for Article A.
topics: important misc
date:   2020-02-08

+++

This is the first article.
`

	ContentB = `
slug:   /article-b/
title:  Article B
intro:  The introduction for Article B.
topics:  important
date:   2021-02-09

+++

This is the second article.
`

	ContentC = `
slug:   /article-c/
title:  Article C
intro:  The introduction for Article C.
topics: important
date:   2020-02-10
draft:  true

+++

This is the third article.
`

	TemplateHome = `
{{ range .Pages }}
Slug:  {{ .Slug }}
Title: {{ .Title }}
Date:  {{ .Date.Format "2006-01-02" }}
Intro: {{ .Intro }}
------------------------------------------------------------------
{{ end }}
`

	TemplateArticle = `
Title:  {{ .Title }}
Intro:  {{ .Intro }}
Date:   {{ .Date.Format "2006-01-02" }}
Topics: {{ range .Topics }}{{ . }}{{ end }}
Content:

{{ .Content }}
`

	TemplateTopics = `
{{ range .Topics }}{{ .Topic }}
	{{ range .Articles }}
	Slug:  {{ .Slug }}
	Title: {{ .Title }}
	Intro: {{ .Intro }}
	Date:  {{ .Date.Format "2006-01-02" }}
	{{ end }}
{{ end }}
`

	RenderedArticleA = `
Title:  Article A
Intro:  The introduction for Article A.
Date:   2020-02-08
Topics: importantmisc
Content:

<p>This is the first article.</p>
`

	RenderedTopics = `
important
	
	Slug:  /article-a/
	Title: Article A
	Intro: The introduction for Article A.
	Date:  2020-02-08
	
misc
	
	Slug:  /article-a/
	Title: Article A
	Intro: The introduction for Article A.
	Date:  2020-02-08

`

	RenderedHome = `
Slug:  /article-a/
Title: Article A
Date:  2020-02-08
Intro: The introduction for Article A.
------------------------------------------------------------------
`
)
