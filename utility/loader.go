package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

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

var root, _ = os.Getwd()

var path = root + "/books.csv"

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

// Loader helps to load the csv data into postgres
func loader() {
	fmt.Println("** Welcome to data loader")
	activeDB, _ := schemeLoader()
	f, _ := os.Open(path)
	readers := csv.NewReader(f)
	for {
		row, err := readers.Read()
		if err == io.EOF {
			fmt.Println("WARN: EOD reached.")
			break
		} else if err != nil {
			fmt.Printf("WARN: Unable to read row record, %s", err)
		} else if len(row) != 10 {
			fmt.Printf("WARN: row  record count is not 10, %s", err)
		} else {
			fmt.Println("Entering record ....", row[0])
			avg, _ := strconv.ParseFloat(row[3], 64)
			bookid, _ := strconv.Atoi(row[0])
			numPage, _ := strconv.Atoi(row[7])
			ratings, _ := strconv.Atoi(row[8])
			reviews, _ := strconv.Atoi(row[9])
			book := BookArchieve{
				Bookid:        bookid,
				Title:         row[1],
				Authors:       row[2],
				AverageRating: avg,
				ISBN:          row[4],
				ISBN13:        row[5],
				LanguageCode:  row[6],
				NumPages:      numPage,
				Ratings:       ratings,
				Reviews:       reviews,
			}
			dataLoader(activeDB, book)
		}
	}
	fmt.Println("-- END -- *")
}

func main() {
	loader()
}
