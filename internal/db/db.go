package db

import (
	"go-recipes-server/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	dsn := "host=localhost port=5432 user=gorecipe password=example dbname=gorecipe sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&model.Recipe{}, &model.User{})
	return db
}
