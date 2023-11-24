package models

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	ID               uint `gorm:"primaryKey"`
	UserID           string
	OrderDescription string
	SourceServer     string
	SourceChannel    string
	Assignee         string
	Delivered        bool
	CreatedAt        time.Time `gorm:"<-:create type:datetime"`
	AcceptedAt       time.Time `gorm:"type:datetime"`
	DeliveredAt      time.Time `gorm:"type:datetime"`
	DeletedAt        gorm.DeletedAt
}
