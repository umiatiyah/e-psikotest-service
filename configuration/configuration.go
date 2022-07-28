package configuration

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var err error

func OpenConnection() *sql.DB {

	err = godotenv.Load()
	var db *sql.DB
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	}

	if os.Getenv("ENV") == "dev" {
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USERNAME")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")

		psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)

		db, err = sql.Open("postgres", psqlInfo)

	} else {
		db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	}
	if err != nil {
		panic(err)
	}

	return db
}
