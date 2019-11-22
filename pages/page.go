package pages

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/mdwhatcott/static/contracts"
	"github.com/russross/blackfriday"
)

func ParseAll(files map[contracts.Path]contracts.File) (pages []contracts.Page) {
	for path, file := range files {
		page := Parse(file)
		page.Path = path
		pages = append(pages, page)
	}
	return pages
}

func Parse(file contracts.File) (page contracts.Page) {
	frontMatter, content := splitFrontMatterFromContent(string(file))
	_, page.ParseError = toml.Decode(frontMatter, &page.FrontMatter)
	if page.ParseError == nil {
		page.OriginalContent = content
		page.HTMLContent = string(blackfriday.Run([]byte(content))) // TODO: footnotes option
	}
	return page
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
