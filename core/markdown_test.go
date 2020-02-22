package core

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestGoldmarkMarkdownConverterFixture(t *testing.T) {
	gunit.Run(new(GoldmarkMarkdownConverterFixture), t)
}

type GoldmarkMarkdownConverterFixture struct {
	*gunit.Fixture

	converter *GoldmarkMarkdownConverter
}

func (this *GoldmarkMarkdownConverterFixture) Setup() {
	this.converter = NewGoldmarkMarkdownConverter()
}

func (this *GoldmarkMarkdownConverterFixture) Test() {
	output, err := this.converter.Convert(MARKDOWN_INPUT)
	this.So(err, should.BeNil)
	this.So(output, should.Equal, EXPECTED_HTML_OUTPUT)

	// Assert correct reuse of internal buffer:
	output2, err2 := this.converter.Convert(MARKDOWN_INPUT)
	this.So(err2, should.BeNil)
	this.So(output2, should.Equal, EXPECTED_HTML_OUTPUT)
}

const MARKDOWN_INPUT = `
# H1

## H2

### H3

#### H4

- a
- b
- c

1. 1
2. 2
3. 3

[link](/target)

` + "`code`" + `

` + "```go\nfenced code\n```" + `

> blockquote

_emphasized_

**bolded**

---

`

const EXPECTED_HTML_OUTPUT = `<h1>H1</h1>
<h2>H2</h2>
<h3>H3</h3>
<h4>H4</h4>
<ul>
<li>a</li>
<li>b</li>
<li>c</li>
</ul>
<ol>
<li>1</li>
<li>2</li>
<li>3</li>
</ol>
<p><a href="/target">link</a></p>
<p><code>code</code></p>
<pre><code class="language-go">fenced code
</code></pre>
<blockquote>
<p>blockquote</p>
</blockquote>
<p><em>emphasized</em></p>
<p><strong>bolded</strong></p>
<hr>
`
