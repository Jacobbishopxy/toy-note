package entity

import (
	"time"
)

/*
Post

A data structure used for storing notes.

- id
- title
- subtitle
- content
- created_at
- updated_at
*/
type Post struct {
	UintId
	Title      string      `gorm:"size:100;not null" json:"title"`
	Subtitle   string      `gorm:"size:100" json:"subtitle,omitempty"`
	Content    string      `gorm:"text;not null" json:"content"`
	Date       time.Time   `gorm:"index;not null" json:"date"`
	Affiliates []Affiliate `gorm:"foreignKey:PostRefer;references:Id" json:"affiliates"`
	Tags       []Tag       `gorm:"many2many:posts_tags;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"tags"`
	Dates
}
