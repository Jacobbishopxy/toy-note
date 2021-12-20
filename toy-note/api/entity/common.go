package entity

import "time"

type UnitId struct {
	Id uint `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
}

type Dates struct {
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
}

type Pagination struct {
	Page int
	Size int
}

func NewPagination(page, size int) Pagination {
	return Pagination{
		Page: page,
		Size: size,
	}
}

type FileObject struct {
	Filename string
	Content  []byte
	Size     int64
}
