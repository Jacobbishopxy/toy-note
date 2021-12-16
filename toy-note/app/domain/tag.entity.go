package domain

import "time"

// Tag
//
// Various tags for notes, which later on can be used as search criteria.
type Tag struct {
	Id          uint
	Name        string
	Description string
	Color       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
