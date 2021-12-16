package domain

import "time"

// Affiliate
//
// A pointer to a file stored in a remote storage.
type Affiliate struct {
	Id        uint
	Filename  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
