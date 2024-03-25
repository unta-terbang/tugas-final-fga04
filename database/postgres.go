package database

import (
	"main.go/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	POSTGRES_URL="postgres://default:CHNhY9Uw6fjQ@ep-late-queen-a1jhdpqs-pooler.ap-southeast-1.aws.neon.tech:5432/verceldb?sslmode=require"
	POSTGRES_PRISMA_URL="postgres://default:CHNhY9Uw6fjQ@ep-late-queen-a1jhdpqs-pooler.ap-southeast-1.aws.neon.tech:5432/verceldb?sslmode=require&pgbouncer=true&connect_timeout=15"
	POSTGRES_URL_NO_SSL="postgres://default:CHNhY9Uw6fjQ@ep-late-queen-a1jhdpqs-pooler.ap-southeast-1.aws.neon.tech:5432/verceldb"
	POSTGRES_URL_NON_POOLING="postgres://default:CHNhY9Uw6fjQ@ep-late-queen-a1jhdpqs.ap-southeast-1.aws.neon.tech:5432/verceldb?sslmode=require"
	POSTGRES_USER="default"
	POSTGRES_HOST="ep-late-queen-a1jhdpqs-pooler.ap-southeast-1.aws.neon.tech"
	POSTGRES_PASSWORD="CHNhY9Uw6fjQ"
	POSTGRES_DATABASE="verceldb"
)

var (
	db  *gorm.DB
	err error
)

func StartDB() {
	config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)

	db, err = gorm.Open(postgres.Open(config), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.Debug().AutoMigrate(models.User{}, models.Photo{}, models.Comment{}, models.SocialMedia{})
}

func GetDB() *gorm.DB {
	if DEBUG_MODE {
		return db.Debug()
	}

	return db
}
