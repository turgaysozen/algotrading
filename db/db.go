package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/turgaysozen/algotrading/monitoring/metrics"
)

var Database *sql.DB

func InitializeDB() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		metrics.RecordError("db_load_env_error")
		return nil, err
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbSslMode := os.Getenv("DB_SSLMODE")

	connStr := "user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " host=" + dbHost + " port=" + dbPort + " sslmode=" + dbSslMode

	Database, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database connection: ", err)
		metrics.RecordError("db_open_connection_error")
		return nil, err
	}

	Database.SetMaxOpenConns(20)
	Database.SetMaxIdleConns(10)
	Database.SetConnMaxLifetime(30 * time.Minute)

	if err := Database.Ping(); err != nil {
		log.Fatal("Error connecting to the database: ", err)
		metrics.RecordError("db_ping_error")
		return nil, err
	}

	log.Println("Database connection established successfully")

	return Database, nil
}
