package content

import (
	"sort"

	"github.com/mdwhatcott/huguinho/contracts"
)

func organizeContent(articles []contracts.Article__DEPRECATED) contracts.Site__DEPRECATED {
	site := make(contracts.Site__DEPRECATED)
	site[contracts.HomePageListingID__DEPRECATED] = orderByMostRecentDate(articles)

	for _, tag := range allTags(articles) {
		site[tag] = filterByTag(articles, tag)
	}
	return site
}
func orderByMostRecentDate(articles []contracts.Article__DEPRECATED) (ordered []contracts.Article__DEPRECATED) {
	for _, article := range articles {
		ordered = append(ordered, article)
	}
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].Date.After(ordered[j].Date)
	})
	return ordered
}
func allTags(articles []contracts.Article__DEPRECATED) (tags []string) {
	all := make(map[string]struct{})
	for _, article := range articles {
		for _, tag := range article.Tags {
			all[tag] = struct{}{}
		}
	}
	for tag := range all {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}
func filterByTag(articles []contracts.Article__DEPRECATED, tag string) (filtered []contracts.Article__DEPRECATED) {
	for _, article := range articles {
		if contains(article.Tags, tag) {
			filtered = append(filtered, article)
		}
	}
	return orderByMostRecentDate(filtered)
}
func contains(haystack []string, needle string) bool {
	for _, straw := range haystack {
		if straw == needle {
			return true
		}
	}
	return false
}
