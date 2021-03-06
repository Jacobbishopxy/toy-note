package entity

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type UintId struct {
	Id uint `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
}

type Dates struct {
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
}

// make sure both created_at and updated_at cannot be manually set by user
func (d *Dates) BeforeCreate(tx *gorm.DB) (err error) {
	if !d.CreatedAt.IsZero() || !d.UpdatedAt.IsZero() {
		err = errors.New("can't manually set CreatedAt or UpdatedAt")
	}
	return
}

func (d *Dates) BeforeUpdate(tx *gorm.DB) (err error) {
	if !d.CreatedAt.IsZero() || !d.UpdatedAt.IsZero() {
		err = errors.New("can't manually set CreatedAt or UpdatedAt")
	}
	return
}

// request from frontend
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

// response to frontend
type FileObject struct {
	Filename string
	Content  []byte
	Size     int64
}

type TimeType uint

const (
	DATE TimeType = iota
	CREATE
	MODIFY
)

type TimeSearch struct {
	Start time.Time
	End   time.Time
	Type  TimeType
}
