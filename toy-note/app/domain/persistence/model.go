package persistence

import "time"

type UnitId struct {
	Id uint `gorm:"primaryKey;autoIncrement"`
}

type Common struct {
	UnitId
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Affiliate struct {
	Common
	PostID     uint   `gorm:"not null"`
	Filename   string `gorm:"not null"`
	Collection string `gorm:"not null"`
	ObjectId   string `gorm:"not null"`
}

type Tag struct {
	Common
	Name        string  `gorm:"size:100;not null;unique"`
	Description string  `gorm:"size:100"`
	Color       string  `gorm:"size:100"`
	Posts       []*Post `gorm:"many2many:post_tag;constraint:OnDelete:SET_NULL;"`
}

type Post struct {
	Common
	Title     string       `gorm:"size:100;not null"`
	Content   string       `gorm:"text;not null"`
	Affiliate []*Affiliate `gorm:"polymorphic:Owner"`
	Tags      []*Tag       `gorm:"many2many:article_tag;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
