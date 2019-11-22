package content

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/mdwhatcott/static/contracts"
	"github.com/russross/blackfriday"
)

func ParseAll(files map[contracts.Path]contracts.File) contracts.ContentListing {
	return organizeContent(parseAll(files))
}

func parseAll(files map[contracts.Path]contracts.File) (articles []contracts.Article) {
	for path, file := range files {
		article := Parse(file)
		article.Path = path
		articles = append(articles, article)
	}
	return articles
}

func Parse(file contracts.File) (article contracts.Article) {
	frontMatter, content := splitFrontMatterFromContent(string(file))
	_, article.ParseError = toml.Decode(frontMatter, &article.FrontMatter)
	if article.ParseError == nil {
		article.OriginalContent = content
		article.HTMLContent = string(blackfriday.Run([]byte(content))) // TODO: footnotes option
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
