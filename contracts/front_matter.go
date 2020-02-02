package contracts

import "time"

type JSONFrontMatter struct {
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Tags        []string  `json:"tags"`
	IsDraft     bool      `json:"draft"`
}

const FRONT_MATTER_DIVIDER = "+++"
