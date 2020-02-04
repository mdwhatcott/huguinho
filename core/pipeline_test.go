package core

import "github.com/mdwhatcott/huguinho/contracts"

func gather(output chan contracts.Article) (pages []contracts.Article) {
	for page := range output {
		pages = append(pages, page)
	}
	return pages
}
