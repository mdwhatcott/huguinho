package site

import (
	"sort"

	"github.com/mdwhatcott/static/contracts"
)

func OrderByDate(pages []contracts.Page) (ordered []contracts.Page) {
	for _, page := range pages {
		ordered = append(ordered, page)
	}
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].Date.Before(ordered[j].Date)
	})
	return ordered
}
