package db

import (
	"github.com/acme-sky/airline-api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB = nil

func InitDb(dsn string) (*gorm.DB, error) {
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err == nil {
		db.AutoMigrate(&models.User{}, &models.Airport{}, &models.Flight{})
	}

	return db, err
}

func GetDb() *gorm.DB {
	return db
}
