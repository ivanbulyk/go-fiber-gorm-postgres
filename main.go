package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ivanbulyk/go-fiber-gorm-postgres/models"
	"github.com/ivanbulyk/go-fiber-gorm-postgres/storage"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}
type Query struct {
	Date      string  `json:"date"`
	Time      string  `json:"time"`
	TimeSpent float64 `json:"time_spent"`
	SQL       string  `json:"sql"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{}

	err := context.BodyParser(&book)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Debug().Create(&book).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create book"})
		return err
	}

	query := ReadLineNext()

	err = r.DB.Create(query).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create query"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "book has been added"})
	return nil
}

func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	bookModel := models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.Debug().Delete(bookModel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete book",
		})
		return err.Error
	}

	query := ReadLineNext()

	err = r.DB.Create(query)
	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create query"})
		return err.Error
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book deleted successfully",
	})

	return nil
}

func (r *Repository) GetBooks(context *fiber.Ctx) error {
	bookModels := &[]models.Books{}

	err := r.DB.Debug().Find(bookModels).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get books"})
		return err
	}

	query := ReadLineNext()

	err = r.DB.Create(query).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create query"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "books fetched successfully",
		"data": bookModels})
	return nil
}

func (r *Repository) GetBookByID(context *fiber.Ctx) error {
	id := context.Params("id")
	bookModel := &models.Books{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	fmt.Println("The ID is: ", id)

	err := r.DB.Debug().Where("id = ?", id).First(bookModel).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not get the book",
		})
		return err
	}

	query := ReadLineNext()

	err = r.DB.Create(query).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create query"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book by ID fetched successfully",
		"data":    bookModel,
	})

	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("/delete_books/:id", r.DeleteBook)
	api.Get("/get_books/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)
	api.Get("/queries", r.GetQueriesDesc)
	api.Get("/queries_asc", r.GetQueriesAsc)
	api.Get("/sql/:sql", r.GetBySQLQuery)
	api.Get("/pages", r.GetQueriesDescPaginated)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := storage.NewConnection(&config)
	if err != nil {
		log.Fatal("could not load the database")
	}

	err = models.MigrateQueries(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	r := Repository{
		DB: db,
	}
	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}

func ReadLineNext() (res Query) {

	scanner := bufio.NewScanner(&storage.Body)

	// Default scanner is bufio.ScanLines.
	// Could also use a custom function of SplitFunc type
	scanner.Split(bufio.ScanWords)

	scannedSlice := make([]string, 0)
	backedSQL := ""

	// Scan for next token
	for scanner.Scan() {

		if scanner.Text() != "" {

			scannedSlice = append(scannedSlice, scanner.Text())
			if InSetSQL(scanner.Text()) {
				backedSQL = backedSQL + scanner.Text()
			}

		}
		continue
	}

	res.Date = scannedSlice[0]
	res.Time = scannedSlice[1]
	if strings.Contains(scannedSlice[3], "ms") {
		str := strings.Split(scannedSlice[3], "")

		timeSpent := ""
		ret := make([]string, 0)

		for ind, val := range str {
			if val == "[" {
				ret = str[ind+1:]
			}
		}

		for ind, val := range ret {
			if val == "m" {

				ret = ret[:ind]
				continue

			}
		}

		for _, c := range ret {

			timeSpent += c

		}

		f, err := strconv.ParseFloat(timeSpent, 64)
		if err != nil {
			log.Println("could not parse str to float64")
		}
		res.TimeSpent = f
	} else {
		res.TimeSpent = 0.0
	}

	if InSetSQL(scannedSlice[5]) {
		res.SQL = scannedSlice[5]
	} else {
		res.SQL = backedSQL
	}

	fmt.Println(res)

	return res

}

func InSeNottString(c string) bool {
	switch c {
	case "[", "]", "m", "s", "ms":
		return false
	}
	return true
}

func InSetSQL(sql string) bool {
	switch sql {
	case "SELECT", "INSERT", "UPDATE", "DELETE":
		return true
	}
	return false
}

func (r *Repository) GetQueriesDesc(context *fiber.Ctx) error {
	queries := &[]models.Queries{}

	err := r.DB.Order("time_spent desc").Find(&queries).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get queries"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "queries fetched successfully",
		"data": queries})
	return nil
}

func (r *Repository) GetQueriesAsc(context *fiber.Ctx) error {
	queries := &[]models.Queries{}

	err := r.DB.Order("time_spent asc").Find(&queries).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get queries"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "queries fetched successfully",
		"data": queries})
	return nil
}

func (r *Repository) GetBySQLQuery(context *fiber.Ctx) error {
	sql := context.Params("sql")
	queries := &[]models.Queries{}
	if sql == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "sql cannot be empty",
		})
		return nil
	}

	fmt.Println("The SQL is: ", sql)

	f := 0.0
	err := r.DB.Order("time_spent desc").Where("sql = ? AND time_spent <> ?", sql, f).Find(&queries).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not get the queries",
		})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "queries by SQL fetched successfully",
		"data":    queries,
	})

	return nil
}

func (r *Repository) GetQueriesDescPaginated(context *fiber.Ctx) error {
	queries := &[]models.Queries{}

	err := r.DB.Scopes(Paginate(context)).Find(&queries).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get queries"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "queries fetched successfully",
		"data": queries})
	return nil
}

func Paginate(context *fiber.Ctx) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		q := context.Query("page")
		page, _ := strconv.Atoi(q)
		if page == 0 {
			page = 1
		}

		qs := context.Query("page_size")
		pageSize, _ := strconv.Atoi(qs)
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
