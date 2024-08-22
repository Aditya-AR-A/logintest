package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() error {
	// Database connection parameters
	dbUser := "aditya"
	dbPass := "9211" // Replace with actual password
	dbName := "login_credentials"
	dbHost := "localhost"
	dbPort := "3306"

	// Create the database connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)

	// Open database connection
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	// Test the connection
	err = DB.Ping()
	if err != nil {
		return err
	}

	fmt.Println("Successfully connected to the database")
	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
