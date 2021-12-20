package entity

import (
	"time"
)

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
	UnitId
	Title      string      `gorm:"size:100;not null"`
	Subtitle   string      `gorm:"size:100"`
	Content    string      `gorm:"text;not null"`
	Date       time.Time   `gorm:"index;not null"`
	Affiliates []Affiliate `gorm:"foreignKey:PostRefer;references:Id"`
	Tags       []Tag       `gorm:"many2many:posts_tags;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Dates
}
