package domain

import "time"

type Tag struct {
	Id          uint
	Name        string
	Description string
	Color       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
