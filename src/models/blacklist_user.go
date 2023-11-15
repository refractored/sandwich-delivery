package models

import (
	"time"
)

type BlacklistUser struct {
	ID        uint `gorm:"primaryKey"`
	UserID    string
	CreatedAt time.Time  `gorm:"type:datetime"`
	DeletedAt *time.Time `gorm:"index"`
}
