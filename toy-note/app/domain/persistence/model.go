package persistence

import (
	"time"

	"gorm.io/gorm"
)

type UnitId struct {
	Id uint `gorm:"primaryKey;autoIncrement"`
}

type Dates struct {
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// Affiliate
type Affiliate struct {
	ObjectId  string `gorm:"not null"`
	Filename  string `gorm:"not null"`
	PostRefer uint   `gorm:"not null"`
	Dates
}

// Tag
type Tag struct {
	UnitId
	Name        string `gorm:"size:100;not null;unique"`
	Description string `gorm:"size:100"`
	Color       string `gorm:"size:100"`
	Posts       []Post `gorm:"many2many:post_tag;constraint:OnDelete:SET NULL;"`
	Dates
}

// Post
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
