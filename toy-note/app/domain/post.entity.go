package domain

type Post struct {
	Common
	Title     string           `gorm:"size:100;not null"`
	Content   string           `gorm:"text;not null"`
	Affiliate []*PostAffiliate `gorm:"polymorphic:Owner"`
	Tags      []*Tag           `gorm:"many2many:article_tag;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
