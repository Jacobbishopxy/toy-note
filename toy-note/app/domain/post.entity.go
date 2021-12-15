package domain

import "time"

// Post
//
// A data structure used for storing notes.
// - id
// - title
// - content
// - created_at
// - updated_at
type Post struct {
	Id        uint
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
