package entities

import "time"

type User struct {
	ID                     string `gorm:"primaryKey"`
	Email                  string
	RefreshToken           string
	RefreshTokenExpiryTime time.Time
}
