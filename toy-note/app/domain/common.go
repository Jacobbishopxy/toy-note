package domain

import "time"

type UnitId struct {
	Id uint `gorm:"primaryKey;autoIncrement"`
}

type Common struct {
	UnitId
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
