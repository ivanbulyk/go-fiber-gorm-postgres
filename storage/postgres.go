package storage

import (
	"bytes"
	"fmt"
	"log"

	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
}

var Body bytes.Buffer

func NewConnection(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s password=%s user=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Password, config.User, config.DBName, config.SSLMode,
	)

	newLogger := logger.New(
		log.New(&Body, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
			LogLevel:                  logger.Warn,            // Log level
			IgnoreRecordNotFoundError: false,                  // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,                   // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		return db, err
	}

	return db, nil
}
