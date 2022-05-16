package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("./test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&CACertificate{})
	db.AutoMigrate(&HostCertificate{})
	db.AutoMigrate(&UserCertificate{})

	return db, nil
}
