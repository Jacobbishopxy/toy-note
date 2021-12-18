package entity

import "time"

type UnitId struct {
	Id uint `gorm:"primaryKey;autoIncrement"`
}

type Dates struct {
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
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
