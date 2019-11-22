package pages

import (
	"sort"

	"github.com/mdwhatcott/static/contracts"
)

func OrganizeContent(articles []contracts.Article) contracts.ContentListing {
	byTag := make(map[string][]contracts.Article)
	for _, tag := range allTags(articles) {
		byTag[tag] = filterByTag(articles, tag)
	}
	return contracts.ContentListing{
		All:   orderByMostRecentDate(articles),
		ByTag: byTag,
	}
}
func orderByMostRecentDate(articles []contracts.Article) (ordered []contracts.Article) {
	for _, article := range articles {
		ordered = append(ordered, article)
	}
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].Date.After(ordered[j].Date)
	})
	return ordered
}
func allTags(articles []contracts.Article) (tags []string) {
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
func filterByTag(articles []contracts.Article, tag string) (filtered []contracts.Article) {
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
