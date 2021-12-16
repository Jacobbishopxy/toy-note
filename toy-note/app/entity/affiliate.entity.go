package entity

// Affiliate
//
// A pointer to a file stored in a remote storage.
type Affiliate struct {
	ObjectId  string `gorm:"not null"`
	Filename  string `gorm:"not null"`
	PostRefer uint   `gorm:"not null"`
	Dates
}
