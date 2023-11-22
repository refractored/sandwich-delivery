package models

import (
	"time"
)

type Order struct {
	ID               uint `gorm:"primaryKey"`
	UserID           string
	OrderDescription string
	SourceServer     string
	SourceChannel    string
	Delivered        bool
	CreatedAt        time.Time `gorm:"<-:create type:datetime"`
	DeliveredAt      time.Time `gorm:"type:datetime"`
	DeletedAt        time.Time `gorm:"index"`
}
