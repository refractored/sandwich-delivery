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
	Status           DeliveryStatus
	CreatedAt        time.Time `gorm:"<-:create type:datetime"`
	AcceptedAt       time.Time `gorm:"type:datetime"`
	InTransitAt      time.Time `gorm:"type:datetime"`
	DeliveredAt      time.Time `gorm:"type:datetime"`
	DeletedAt        gorm.DeletedAt
}
type DeliveryStatus uint8

const (
	StatusWaiting   DeliveryStatus = 0
	StatusAccepted  DeliveryStatus = 1
	StatusInTransit DeliveryStatus = 2
	StatusDelivered DeliveryStatus = 3
	StatusCancelled DeliveryStatus = 4
)
