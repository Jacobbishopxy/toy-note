package domain

import "time"

// Post
//
// A data structure used for storing notes.
// - id
// - title
// - subtitle
// - content
// - created_at
// - updated_at
type Post struct {
	Id        uint
	Title     string
	Subtitle  string
	Content   string
	Date      time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
