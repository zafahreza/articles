package app

import (
	"articles/command/domain"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"
)

func InitDB() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	host := os.Getenv("POSTGRES_HOST")

	dsn := fmt.Sprintf("host=%s user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai", host)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	pool, err := db.DB()
	if err != nil {
		panic(err)
	}

	pool.SetMaxIdleConns(5)
	pool.SetMaxOpenConns(100)
	pool.SetConnMaxIdleTime(10 * time.Minute)
	pool.SetConnMaxLifetime(60 * time.Minute)

	err = db.AutoMigrate(&domain.Article{})
	if err != nil {
		panic(err)
	}

	return db
}
