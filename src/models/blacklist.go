package models

import (
	"time"
)

// TODO: USE int FOR USERID IN THE FUTURE
type BlacklistUser struct {
	ID        uint `gorm:"primaryKey"`
	UserID    string
	CreatedAt time.Time  `gorm:"type:datetime"`
	DeletedAt *time.Time `gorm:"index"`
}
