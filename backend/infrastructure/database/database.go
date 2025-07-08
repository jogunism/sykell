package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

// InitDB initializes and returns a new database connection
func InitDB(dbConnStr string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dbConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Ping the database to verify connection
	err = db.Ping()
	if err != nil {
		db.Close() // Close the connection if ping fails
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Successfully connected to MySQL database!")
	return db, nil
}
