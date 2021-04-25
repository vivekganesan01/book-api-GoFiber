package main

import (
	"fmt"

	"book_api/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"gorm.io/gorm"
)

// APIVersion holds all app metadata
type APIVersion struct {
	Name       string
	APIVersion string
}

// BookParser to query records
type BookParser struct {
	Bookid        int     `json:"bookid" xml:"bookid" form:"bookid" query:"bookid"`
	Title         string  `json:"title" xml:"title" form:"title" query:"title"`
	Authors       string  `json:"author" xml:"author" form:"author" query:"author"`
	AverageRating float64 `json:"averagerating" xml:"averagerating" form:"averagerating" query:"averagerating"`
	ISBN          string  `json:"isbn" xml:"isbn" form:"isbn" query:"isbn"`
	ISBN13        string  `json:"isbn13" xml:"isbn13" form:"isbn13" query:"isbn13"`
	LanguageCode  string  `json:"languagecode" xml:"languagecode" form:"languagecode" query:"languagecode"`
	NumPages      int     `json:"numpage" xml:"numpage" form:"numpage" query:"numpage"`
	Ratings       int     `json:"ratings" xml:"ratings" form:"ratings" query:"ratings"`
	Reviews       int     `json:"reviews" xml:"reviews" form:"reviews" query:"reviews"`
}

var activeDB *gorm.DB

func main() {

	bookapi := fiber.New(fiber.Config{
		StrictRouting: true,
	})

	//Handles basic authorization to API
	// todo: need to hide/fetch the sha from boltDB
	bookapi.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			"YXBpYnl2aXZlawo=": "ac8deb2b53db3df8417de6ff9f2945ea601d3ab09c51f491323563f8a73e83d6",
		},
	}))

	// Initialize DB properties
	bookapi.Use(func(c *fiber.Ctx) error {
		var err error
		activeDB, err = db.ActiveDBInstance()
		if err != nil {
			c.Status(404)
			return c.JSON(fmt.Sprintf("error: db not working. %v", err))
		}
		return c.Next()
	})

	// Gets API Version metadata
	bookapi.Get("/", func(c *fiber.Ctx) error {
		data := &APIVersion{
			Name:       "Book-API",
			APIVersion: "v0.0.1",
		}
		return c.JSON(data)
	})

	// Gets firstbook
	bookapi.Get("/firstbook", func(c *fiber.Ctx) error {
		var firstBook db.BookArchieve
		activeDB.First(&firstBook)
		return c.JSON(firstBook)
	})
	// Gets last updated book
	bookapi.Get("/latest/updated", func(c *fiber.Ctx) error {
		var latestUpdated db.BookArchieve
		activeDB.Order("updated_at desc").Find(&latestUpdated)
		return c.JSON(latestUpdated)
	})
	// Gets last created book
	bookapi.Get("/latest/created", func(c *fiber.Ctx) error {
		var latestCreated db.BookArchieve
		activeDB.Order("created_at desc").Find(&latestCreated)
		return c.JSON(latestCreated)
	})

	// Gets book based on bookid
	bookapi.Get("/byid/:bookid", func(c *fiber.Ctx) error {
		bookid := c.Params("bookid")
		if bookid == "" {
			c.Status(404)
			return c.JSON("Please provide some book id (i.e., /book-api/v1/{BOOKID})")
		}
		var book db.BookArchieve
		result := activeDB.Where("bookid = ?", bookid).Find(&book)
		if result.RowsAffected == 0 {
			c.Status(404)
			return c.JSON(fmt.Sprintf("error: book id %v doesn't exists", bookid))
		}
		return c.JSON(book)
	})

	// Gets book based on given author name
	bookapi.Get("/byauthor", func(c *fiber.Ctx) error {
		authorname := c.Query("name")
		if authorname == "" {
			c.Status(404)
			return c.JSON(fmt.Sprintf("error: please pass valid author name via query /book-api/v1/byauthor?name=****"))
		}
		var books []db.BookArchieve
		activeDB.Where("Authors = ?", authorname).Find(&books)
		return c.JSON(books)
	})

	// Gets total record count
	bookapi.Get("/bookcount", func(c *fiber.Ctx) error {
		var books []db.BookArchieve
		result := activeDB.Find(&books)
		data := map[string]interface{}{
			"api":        "bookapi",
			"totalcount": result.RowsAffected,
		}
		return c.JSON(data)
	})

	// Gets book based on range from: to:
	bookapi.Get("/bookid/range", func(c *fiber.Ctx) error {
		from := c.Query("from")
		to := c.Query("to")
		var books []db.BookArchieve
		activeDB.Where("bookid BETWEEN ? AND ?", from, to).Find(&books)
		return c.JSON(books)
	})

	// Gets book based on given query data
	bookapi.Get("/book", func(c *fiber.Ctx) error {
		book := new(BookParser)
		if err := c.QueryParser(book); err != nil {
			return c.JSON(err)
		}
		var books []db.BookArchieve
		activeDB.Where(book).Find(&books)
		if len(books) == 0 {
			c.Status(404)
			return c.JSON(fmt.Sprintf("error: No record exists for given custom query"))
		}
		return c.JSON(books)
	})

	// Post an update/create to the record
	bookapi.Post("/book/record", func(c *fiber.Ctx) error {
		body := new(BookParser)
		if err := c.BodyParser(body); err != nil {
			return c.JSON(err)
		}
		if body.Bookid == 0 {
			c.Status(404)
			return c.JSON(fmt.Sprintf("error: bookid cannot be %d", body.Bookid))
		}
		var book db.BookArchieve
		result := activeDB.Where("bookid = ?", body.Bookid).Find(&book)
		if result.RowsAffected == 0 {
			result := activeDB.Create(&body)
			if result.Error != nil {
				c.Status(400)
				return c.JSON("error: Not able to update record")
			}
			c.Status(200)
			return c.JSON("record created successfully")
		}
		activeDB.Model(&book).Updates(body)
		return c.JSON("record updated successfully")
	})

	// Deletes a record based on book id
	bookapi.Delete("/book/delete/:bookid", func(c *fiber.Ctx) error {
		bookid := c.Params("bookid")
		if bookid == "" {
			c.Status(404)
			return c.JSON("Please provide some book id (i.e., /book-api/v1/{BOOKID})")
		}
		activeDB.Delete(&db.BookArchieve{}, bookid)
		return c.JSON("record deleted successfully")
	})

	app := fiber.New()
	app.Mount("/book-api/v1", bookapi)
	app.Listen(":9091")
}
