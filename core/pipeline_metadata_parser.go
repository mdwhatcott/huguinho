package core

import "time"

type JSONMetadata struct {
	Draft bool      `json:"draft"`
	Slug  string    `json:"slug"` // TODO: schema change (all articles)
	Title string    `json:"title"`
	Intro string    `json:"intro"` // TODO: schema change (all articles)
	Tags  []string  `json:"tags"`
	Date  time.Time `json:"date"`
}

type JSONMetadataParser struct {
}
