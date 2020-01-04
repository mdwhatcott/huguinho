package content

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	chromaHTML "github.com/alecthomas/chroma/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/renderer/html"

	"github.com/mdwhatcott/huguinho/contracts"
)

func ParseAll(files map[contracts.Path]contracts.File, drafts, future bool) contracts.Site {
	articles := parseAll(files)
	articles = filterDrafts(articles, drafts)
	articles = filterFutures(articles, future)
	return organizeContent(articles)
}

func filterDrafts(articles []contracts.Article, drafts bool) (filtered []contracts.Article) {
	for _, article := range articles {
		if drafts || !article.IsDraft {
			filtered = append(filtered, article)
		}
	}
	return filtered
}

func filterFutures(articles []contracts.Article, future bool) (filtered []contracts.Article) {
	now := time.Now()
	for _, article := range articles {
		if future || article.Date.Before(now) {
			filtered = append(filtered, article)
		}
	}
	return filtered
}

func parseAll(files map[contracts.Path]contracts.File) (articles []contracts.Article) {
	for path, file := range files {
		article := parse(file)
		article.Path = contracts.Path(strings.TrimSuffix(string(path), ".md")) + "/"
		articles = append(articles, article)
	}
	return articles
}

func parse(file contracts.File) (article contracts.Article) {
	frontMatter, content := splitFrontMatterFromContent(string(file))
	_, article.ParseError = toml.Decode(frontMatter, &article.FrontMatter)
	if article.ParseError == nil {
		article.OriginalContent = content

		markdown := goldmark.New(
			goldmark.WithRendererOptions(html.WithUnsafe()),
			goldmark.WithExtensions(
				highlighting.NewHighlighting(
					// All highlighting styles listed here: https://github.com/alecthomas/chroma/tree/master/styles
					// Playground for previewing each style: https://swapoff.org/chroma/playground/
					highlighting.WithStyle("friendly"),
					highlighting.WithFormatOptions(chromaHTML.WithLineNumbers(false)),
				),
			),
		)
		buffer := new(bytes.Buffer)
		err := markdown.Convert([]byte(content), buffer)
		if err != nil {
			log.Println("[INFO] Failed to convert markdown:", err)
		}
		article.Content = buffer.String()
	}
	return article
}

func splitFrontMatterFromContent(file string) (string, string) {
	scanner := bufio.NewScanner(strings.NewReader(file))
	front := scanFrontMatter(file, scanner)
	content := scanContent(scanner)
	return front, content
}

func scanContent(scanner *bufio.Scanner) string {
	content := new(strings.Builder)
	for scanner.Scan() {
		fmt.Fprintln(content, scanner.Text())
	}
	return strings.TrimSpace(content.String())
}

func scanFrontMatter(file string, scanner *bufio.Scanner) string {
	if !strings.HasPrefix(file, tomlSeparator) {
		return ""
	}
	scanFirstTOMLSeparator(scanner)
	return scanFrontMatterBody(scanner)
}

func scanFirstTOMLSeparator(scanner *bufio.Scanner) bool {
	return scanner.Scan()
}

func scanFrontMatterBody(scanner *bufio.Scanner) string {
	front := new(strings.Builder)
	for scanner.Scan() {
		line := scanner.Text()
		if line == tomlSeparator {
			break
		}
		fmt.Fprintln(front, line)
	}
	return strings.TrimSpace(front.String())
}

const tomlSeparator = "+++"
