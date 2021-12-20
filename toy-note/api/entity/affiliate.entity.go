package entity

// Affiliate
//
// A pointer to a file stored in a remote storage.
type Affiliate struct {
	UnitId
	ObjectId  string
	Filename  string `gorm:"not null"`
	PostRefer uint
	Dates
}
