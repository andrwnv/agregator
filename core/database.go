package core

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const url = "host=localhost user=postgres password=852456 dbname=event_aggregator port=5432 sslmode=disable"

func InitDatabase() *gorm.DB {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		panic("Couldn't connect database")
	}

	return db
}
