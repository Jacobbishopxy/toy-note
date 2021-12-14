package domain

type Tag struct {
	Common
	Name        string  `gorm:"size:100;not null;unique"`
	Description string  `gorm:"size:100"`
	Color       string  `gorm:"size:100"`
	Posts       []*Post `gorm:"many2many:post_tag;constraint:OnDelete:SET_NULL;"`
}
