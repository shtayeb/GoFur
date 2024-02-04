package internal

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/shtayeb/rssfeed/http/session"
	"github.com/shtayeb/rssfeed/http/types"
	"github.com/shtayeb/rssfeed/internal/database"
)

var DB *database.Queries
var Config *types.Config

func InitApp() {
	godotenv.Load(".env")

	appKey := os.Getenv("APP_KEY")
	if appKey == "" {
		log.Fatal("App key is not set. Please create an application key !")
	}

	DB_URL := os.Getenv("DATABASE_URL")
	if DB_URL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
		port = "8000"
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "local"
	}

	db, err := sql.Open("postgres", DB_URL)
	if err != nil {
		log.Fatal(err)
	}

	// database.Queries
	DB = database.New(db)

	Config = &types.Config{
		APP_KEY: appKey,
		APP_ENV: appEnv,
		APP_URL: os.Getenv("APP_URL"),
		PORT:    os.Getenv("PORT"),
		MailConfig: types.MailConfig{
			MAIL_HOST:         os.Getenv("MAIL_HOST"),
			MAIL_PORT:         os.Getenv("MAIL_PORT"),
			MAIL_USERNAME:     os.Getenv("MAIL_USERNAME"),
			MAIL_PASSWORD:     os.Getenv("MAIL_PASSWORD"),
			MAIL_FROM_ADDRESS: os.Getenv("MAIL_FROM_ADDRESS"),
			MAIL_FROM_NAME:    os.Getenv("MAIL_FROM_NAME"),
			MAIL_ENCRYPTION:   os.Getenv("MAIL_ENCRYPTION"),
		},
	}

	// Sesison Manager
	session.InitSessionManager(db)
}
