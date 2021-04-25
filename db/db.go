package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// BookArchieve defines book metadata schema
type BookArchieve struct {
	gorm.Model
	Bookid        int `gorm:"primary_key"`
	Title         string
	Authors       string
	AverageRating float64
	ISBN          string
	ISBN13        string
	LanguageCode  string
	NumPages      int
	Ratings       int
	Reviews       int
}

// Postgres SQL connection string
const dsn string = "host=db user=root password=root dbname=books port=5432 sslmode=disable"

// schemeLoader helps to connect to postgres SQL
func schemeLoader() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("ERROR: Unable to open connection: %s", err)
		return db, err
	}
	db.AutoMigrate(&BookArchieve{})
	return db, err
}

// Inserts record to DB
func dataLoader(activeDB *gorm.DB, data BookArchieve) {
	activeDB.Create(&data)
}

// ActiveDBInstance returns active DB connection
func ActiveDBInstance() (*gorm.DB, error) {
	activeDB, err := schemeLoader()
	return activeDB, err
}
