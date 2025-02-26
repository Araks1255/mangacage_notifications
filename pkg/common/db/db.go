package db

import (
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

func Init(dbUrl string) (db *gorm.DB, err error) {
	db, err = gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}