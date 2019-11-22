package site

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/mdwhatcott/static/contracts"
	"github.com/russross/blackfriday"
)

// TODO: translate map[fs.Path]fs.File to map[fs.Path]Page (parse front-matter, render md)

func ConvertToPage(file contracts.File) contracts.Page {
	_, content := splitFrontMatterFromContent(string(file))
	var page contracts.Page
	// TODO: decode TOML in front matter
	page.OriginalContent = content
	page.HTMLContent = string(blackfriday.Run([]byte(content))) // TODO: footnotes option
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
