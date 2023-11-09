package models

import (
	"time"
)

type Order struct {
	ID               uint `gorm:"primaryKey"`
	UserID           string
	OrderDescription string
	DisplayName      string
	CreatedAt        time.Time  `gorm:"type:datetime"`
	DeletedAt        *time.Time `gorm:"index"`
}
