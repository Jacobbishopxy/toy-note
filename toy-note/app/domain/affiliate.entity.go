package domain

import "time"

type Affiliate struct {
	Id        uint
	Filename  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
