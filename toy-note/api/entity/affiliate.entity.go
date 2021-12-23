package entity

/*
Affiliate

A pointer to a file stored in a remote storage.

- id
- object_id
- filename
- post_refer
- created_at
- updated_at
*/
type Affiliate struct {
	UintId
	ObjectId  string `json:"object_id,omitempty"`
	Filename  string `gorm:"not null" json:"filename"`
	PostRefer uint   `json:"post_refer,omitempty"`
	Dates
}
