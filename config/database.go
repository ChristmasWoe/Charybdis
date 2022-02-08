package config

import (
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"gorm.io/driver/postgres"
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

	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}
	// migrate our model for artist
	//  db.AutoMigrate(&api.Artist{})
	return db
}