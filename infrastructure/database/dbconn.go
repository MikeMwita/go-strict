package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DBConn *sql.DB

func InitDB() error {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		return errors.New("missing database environment variables")
	}

	dbConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	DBConn, err = sql.Open("postgres", dbConnStr)
	if err != nil {
		return err
	}

	err = DBConn.Ping()
	if err != nil {
		return err
	}

	return nil
}
