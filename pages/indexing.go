package pages

import (
	"sort"

	"github.com/mdwhatcott/static/contracts"
)

func OrganizePages(pages []contracts.Page) contracts.SiteListings {
	byTag := make(map[string][]contracts.Page)
	for _, tag := range allTags(pages) {
		byTag[tag] = filterByTag(pages, tag)
	}
	return contracts.SiteListings{
		All:   orderByMostRecentDate(pages),
		ByTag: byTag,
	}
}
func orderByMostRecentDate(pages []contracts.Page) (ordered []contracts.Page) {
	for _, page := range pages {
		ordered = append(ordered, page)
	}
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].Date.After(ordered[j].Date)
	})
	return ordered
}
func allTags(pages []contracts.Page) (tags []string) {
	all := make(map[string]struct{})
	for _, page := range pages {
		for _, tag := range page.Tags {
			all[tag] = struct{}{}
		}
	}
	for tag := range all {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}
func filterByTag(pages []contracts.Page, tag string) (filtered []contracts.Page) {
	for _, page := range pages {
		if contains(page.Tags, tag) {
			filtered = append(filtered, page)
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
