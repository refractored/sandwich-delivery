package models

import "time"

type User struct {
	ID              uint `gorm:"primaryKey"`
	UserID          string
	Credits         uint32
	OrdersCreated   uint32
	OrdersAccepted  uint32
	PermissionLevel UserPermissionLevel
	IsBlacklisted   bool
	CreatedAt       time.Time `gorm:"type:datetime"`
}

type UserPermissionLevel uint8

const (
	PermissionLevelUser   UserPermissionLevel = 0
	PermissionLevelMod    UserPermissionLevel = 1 // Can delete orders
	PermissionLevelArtist UserPermissionLevel = 2 // Can accept orders
	PermissionLevelAdmin  UserPermissionLevel = 3 // Can (un)blacklist users and purge orders
	PermissionLevelOwner  UserPermissionLevel = 4 // Can do everything else (e.g. shutdown)
)
