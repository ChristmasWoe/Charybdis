package config

import (
	"charybdis/api"
	"charybdis/logging"
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"gorm.io/driver/postgres"

	// "gorm.io/logger"
	"gorm.io/gorm"
)

const (
	host     = "172.17.0.2"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "postgres"
)

func OpenConnection() *gorm.DB {
	// Create db instance with gorm
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("Failed to connect to database!")
	}

	db.AutoMigrate(&api.User{})
	db.AutoMigrate(&api.Executor{})
	db.AutoMigrate(&logging.Log{})
	// migrate our model for artist
	//  db.AutoMigrate(&api.Artist{})
	return db
}
