package models

import (
	"time"
)

type Order struct {
	ID               uint `gorm:"primaryKey"`
	UserID           string
	OrderDescription string
	Username         string
	Discriminator    string
	CreatedAt        time.Time  `gorm:"type:datetime"`
	DeletedAt        *time.Time `gorm:"index"`
}
