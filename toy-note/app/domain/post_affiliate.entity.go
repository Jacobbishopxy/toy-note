package domain

type PostAffiliate struct {
	Common
	PostID     uint   `gorm:"not null"`
	Filename   string `gorm:"not null"`
	Collection string `gorm:"not null"`
	ObjectId   string `gorm:"not null"`
}
