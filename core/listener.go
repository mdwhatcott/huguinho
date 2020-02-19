package core

import (
	"github.com/mdwhatcott/huguinho/contracts"
)

func Listen(in, out chan contracts.Article, handler contracts.Handler) {
	defer close(out)

	for article := range in {
		if article.Error == nil {
			handler.Handle(&article)
		}
		out <- article
	}
}
