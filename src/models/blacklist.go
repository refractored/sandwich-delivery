package models

import (
	"time"
)

// TODO: USE uint64 FOR USERID IN THE FUTURE
type BlacklistUser struct {
	ID        uint `gorm:"primaryKey"`
	UserID    string
	CreatedAt time.Time  `gorm:"type:datetime"`
	DeletedAt *time.Time `gorm:"index"`
}
