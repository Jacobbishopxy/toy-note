package entity

import (
	"time"

	"gorm.io/gorm"
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
	gorm.Model
	UnitId
	Title     string      `gorm:"size:100;not null"`
	Subtitle  string      `gorm:"size:100"`
	Content   string      `gorm:"text;not null"`
	Date      time.Time   `gorm:"index;not null"`
	Affiliate []Affiliate `gorm:"foreignKey:PostRefer"`
	Tags      []Tag       `gorm:"many2many:article_tag;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Dates
}
