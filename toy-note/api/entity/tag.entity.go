package entity

/*
Tag

Various tags for notes, which later on can be used as search criteria.

- id
- name
- description
- color
- posts
- created_at
- updated_at
*/
type Tag struct {
	UnitId
	Name        string `gorm:"size:100;not null;unique" json:"name"`
	Description string `gorm:"size:100" json:"description,omitempty"`
	Color       string `gorm:"size:100" json:"color,omitempty"`
	Posts       []Post `gorm:"many2many:posts_tags;constraint:OnDelete:SET NULL;" json:"posts"`
	Dates
}
