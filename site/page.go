package site

import (
	"strings"

	"github.com/mdwhatcott/static/contracts"
	"github.com/russross/blackfriday"
)

// TODO: translate map[fs.Path]fs.File to map[fs.Path]Page (parse front-matter, render md)

func ConvertToPage(file contracts.File) contracts.Page {
	content := string(file)
	var endFront int
	if strings.HasPrefix(content, TOML) {
		endFront = strings.Index(content[len(TOML):], TOML) + len(TOML) + len(TOML)
	}
	content = content[endFront:]
	return contracts.Page{
		Content: content,
		HTML:    string(blackfriday.Run([]byte(content))), // TODO: footnotes option
	}
}

const TOML = "+++\n"
