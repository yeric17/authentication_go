package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	PORT,
	HOST,
	MODE,
	DB_USER,
	DB_PASS,
	DB_HOST,
	DB_PORT,
	DB_NAME,
	DB_DRIVER,
	JWT_SECRET_KEY,
	GOOGLE_CLIENT_ID,
	GOOGLE_CLIENT_SECRET,
	REDIRECT_GOOGLE_URL,
	CONNECTION_STRING string
)

func init() {
	err := godotenv.Load("app.env")

	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}
	PORT = os.Getenv("PORT")
	MODE = os.Getenv("MODE")
	HOST = os.Getenv("HOST")
	DB_USER = os.Getenv("DB_USER")
	DB_PASS = os.Getenv("DB_PASS")
	DB_HOST = os.Getenv("DB_HOST")
	DB_PORT = os.Getenv("DB_PORT")
	DB_NAME = os.Getenv("DB_NAME")
	DB_DRIVER = os.Getenv("DB_DRIVER")
	GOOGLE_CLIENT_ID = os.Getenv("GOOGLE_CLIENT_ID")
	GOOGLE_CLIENT_SECRET = os.Getenv("GOOGLE_CLIENT_SECRET")
	REDIRECT_GOOGLE_URL = os.Getenv("REDIRECT_GOOGLE_URL")
	JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")

	CONNECTION_STRING = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME)

}
