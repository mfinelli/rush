package db

import (
	"time"
)

type CACertificate struct {
	ID         int64 `gorm:"primaryKey"`
	Type       string
	PublicKey  string
	PrivateKey string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time `gorm:"index"`
}
