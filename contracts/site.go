package contracts

type SiteListings struct {
	All   []Page
	ByTag map[string][]Page
}
