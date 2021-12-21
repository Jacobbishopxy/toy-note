package entity

/*
Affiliate

A pointer to a file stored in a remote storage.
*/
type Affiliate struct {
	UnitId
	ObjectId  string `json:"object_id,omitempty"`
	Filename  string `gorm:"not null" json:"filename"`
	PostRefer uint   `json:"post_refer,omitempty"`
	Dates
}
