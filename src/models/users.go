package models

import "time"

type User struct {
	ID              uint `gorm:"primaryKey"`
	UserID          string
	Credits         uint64
	OrdersCreated   uint64
	OrdersAccepted  uint64
	PermissionLevel uint64
	IsBlacklisted   bool
	CreatedAt       time.Time `gorm:"type:datetime"`
}
