
package main

import (
	"database/sql"
	"fmt"
	"os"
	_ "github.com/mattn/go-sqlite3"
)

func initDB(db *sql.DB) error {
	// Connect to the database
	_, err := db.Exec("PRAGMA foreign_keys = ON") // Enable foreign key support
	if err != nil {
		return err
	}

	// Read schema.sql file
	content, err := os.ReadFile("schema.sql")
	if err != nil {
		return err
	}

	// Execute SQL commands from the schema file
	_, err = db.Exec(string(content))
	if err != nil {
		return err
	}

	return nil
}

func connectDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "/tmp/minitwit.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	// Open database connection
	db, err := connectDB()
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize the database
	err = initDB(db)
	if err != nil {
		fmt.Println("Error initializing the database:", err)
		os.Exit(1)
	}

	fmt.Println("Database initialized successfully.")
}


