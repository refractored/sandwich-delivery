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
	Tipped           bool
	Status           DeliveryStatus `gorm:"default:0"`
	CreatedAt        time.Time      `gorm:"<-:create type:datetime"`
	AcceptedAt       *time.Time     `gorm:"type:datetime"`
	InTransitAt      *time.Time     `gorm:"type:datetime"`
	DeliveredAt      *time.Time     `gorm:"type:datetime"`
	DeletedAt        *gorm.DeletedAt
}

type DeliveryStatus uint8

const (
	StatusWaiting   DeliveryStatus = 0 // Waiting to be accepted by a Sandwich Artist
	StatusAccepted  DeliveryStatus = 1 // Accepted by Sandwich Artist
	StatusInTransit DeliveryStatus = 2 // Invite sent to Sandwich Artist
	StatusDelivered DeliveryStatus = 3 // Sandwich Artist has marked as delivered
	StatusCancelled DeliveryStatus = 4 // Cancelled by user
	StatusModerated DeliveryStatus = 5 // Cancelled by staff
	StatusError     DeliveryStatus = 6 // Cancelled due to error (e.g. invalid location)
)

func (s DeliveryStatus) String() string {
	return [...]string{"Waiting", "Accepted", "In Transit", "Delivered", "Cancelled", "Cancelled (Staff)", "Cancelled (Error)"}[s]
}
