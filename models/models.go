package models

import (
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

type Books struct {
	ID        uint    `gorm:"primary key;autoIncrement" json:"id"`
	Author    *string `json:"author"`
	Title     *string `json:"title"`
	Publisher *string `json:"publisher"`
}

type Queries struct {
	ID        uint     `gorm:"primary key;autoIncrement" json:"id"`
	Date      *string  `json:"date"`
	Time      *string  `json:"time"`
	TimeSpent *float64 `json:"time_spent"`
	SQL       *string  `json:"sql"`
}

func MigrateBooks(db *gorm.DB) error {
	err := db.AutoMigrate(&Books{})
	return err
}

func MigrateQueries(db *gorm.DB) error {
	err := db.AutoMigrate(&Queries{})
	return err
}
