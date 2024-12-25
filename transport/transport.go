package transport

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var Db *sqlx.DB
var err error

func InitDatabase() {
	// Load environment variables from .env file
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve database configuration from environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Initialize sqlx DB instead of sql.DB
	Db, err = sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	// Ping to check if the connection is valid
	err = Db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}
