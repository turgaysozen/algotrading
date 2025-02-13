package db

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/turgaysozen/algotrading/monitoring/metrics"
)

var Database *sql.DB

const maxRetries = 3
const retryDelay = 5 * time.Second

func InitializeDB() (*sql.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbSslMode := os.Getenv("DB_SSLMODE")

	connStr := "user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " host=" + dbHost + " port=" + dbPort + " sslmode=" + dbSslMode

	for i := 0; i < maxRetries; i++ {
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Error opening database connection (attempt %d/%d): %v", i+1, maxRetries, err)
			metrics.RecordError("db_open_connection_error")
			time.Sleep(retryDelay)
			continue
		}

		db.SetMaxOpenConns(50)
		db.SetMaxIdleConns(30)
		db.SetConnMaxLifetime(30 * time.Minute)

		if err := db.Ping(); err != nil {
			log.Printf("Error pinging the database (attempt %d/%d): %v", i+1, maxRetries, err)
			metrics.RecordError("db_ping_error")
			time.Sleep(retryDelay)
			continue
		}

		log.Println("Database connection established successfully")
		Database = db
		return db, nil
	}

	log.Println("Failed to establish database connection after retries")
	metrics.RecordError("db_connection_retry_failure")
	return nil, errors.New("unable to connect to the database after multiple attempts")
}
