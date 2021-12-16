package entity

// Tag
//
// Various tags for notes, which later on can be used as search criteria.
// Tag
type Tag struct {
	UnitId
	Name        string `gorm:"size:100;not null;unique"`
	Description string `gorm:"size:100"`
	Color       string `gorm:"size:100"`
	Posts       []Post `gorm:"many2many:post_tag;constraint:OnDelete:SET NULL;"`
	Dates
}
