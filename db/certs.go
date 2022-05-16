package db

import (
	"time"
)

type HostCertificate struct {
	ID              int64  `gorm:"primaryKey"`
	KeyId           string `gorm:"uniqueIndex"`
	Hostname        string
	Type            string
	NotBefore       time.Time
	NotAfter        time.Time
	PublicKey       string
	SignedPublicKey string
	PrivateKey      string
	CACertificateID int64
	CACertificate   CACertificate
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       time.Time `gorm:"index"`
}

type UserCertificate struct {
	ID              int64  `gorm:"primaryKey"`
	KeyId           string `gorm:"uniqueIndex"`
	Username        string
	Type            string
	NotBefore       time.Time
	NotAfter        time.Time
	PublicKey       string
	SignedPublicKey string
	PrivateKey      string
	CACertificateID int64
	CACertificate   CACertificate
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       time.Time `gorm:"index"`
}
