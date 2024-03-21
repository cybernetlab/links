package models

import (
	"log"
	"os"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var KV *redis.Client

func Setup() {
	setupDB()
	setupKV()
}

func setupDB() {
	dsn, exists := os.LookupEnv("DB_DSN")
	if !exists {
		dsn = "host=localhost user=go password=go dbname=links port=5432 TimeZone=Europe/Moscow"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}
	DB = db
}

func setupKV() {
	dsn, exists := os.LookupEnv("KV_DSN")
	if !exists {
		dsn = "redis://localhost:6379/0?protocol=3"
	}
	opts, err := redis.ParseURL(dsn)
	if err != nil {
		log.Fatal(err)
	}
	KV = redis.NewClient(opts)
}
