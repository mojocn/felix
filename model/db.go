package model

import (
	"gorm.io/driver/sqlite" // Sqlite driver based on CGO
	"gorm.io/gorm"
)

var db *gorm.DB

func initDb() {
	var err error
	db, err = gorm.Open(sqlite.Open("felix.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	migrate()
}

func DB() *gorm.DB {
	if db == nil {
		initDb()
	}
	return db
}
