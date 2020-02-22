package core

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/logging"

	"github.com/mdwhatcott/huguinho/contracts"
)

func TestReporterFixture(t *testing.T) {
	gunit.Run(new(ReporterFixture), t)
}

type ReporterFixture struct {
	*gunit.Fixture
}

func (this *ReporterFixture) Test() {
	started := time.Now()
	stopped := started.Add(time.Millisecond * 42).Add(time.Microsecond * 1)

	stream := make(chan contracts.Article)
	go this.load(stream)

	reporter := NewReporter(started)
	reporter.log = logging.Capture()
	reporter.log.SetFlags(0)

	reporter.ProcessStream(stream)
	reporter.RenderFinalReport(stopped)

	this.So(reporter.Errors(), should.Equal, 3)
	this.So(reporter.log.Log.String(), should.Equal, strings.Join([]string{
		"[INFO] published article: /a",
		"[INFO] dropped article",
		"[INFO] published article: /c",
		"[WARN] error: GOPHERS",
		"[INFO] published article: /e",
		"[INFO] errors encountered:  1",
		"[INFO] dropped articles:    1",
		"[INFO] published articles:  3",
		"[INFO] processing duration: 42ms",
		"",
	}, "\n"))
}

func (this *ReporterFixture) load(stream chan contracts.Article) {
	defer close(stream)
	stream <- contracts.Article{Metadata: contracts.ArticleMetadata{Slug: "/a"}}
	stream <- contracts.Article{Metadata: contracts.ArticleMetadata{Slug: "/b"}, Error: contracts.ErrDroppedArticle}
	stream <- contracts.Article{Metadata: contracts.ArticleMetadata{Slug: "/c"}}
	stream <- contracts.Article{Metadata: contracts.ArticleMetadata{Slug: "/d"}, Error: errors.New("GOPHERS")}
	stream <- contracts.Article{Metadata: contracts.ArticleMetadata{Slug: "/e"}}
}
