package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

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

func queryDB(db *sql.DB, query string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0)

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)
		if err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if val == nil {
				row[col] = nil
			} else {
				row[col] = val
			}
		}

		result = append(result, row)
	}

	return result, nil
}

func getUserID(db *sql.DB, username string) (int64, error) {
	var userID int64
	err := db.QueryRow("SELECT user_id FROM user WHERE username = ?", username).Scan(&userID)
	if err == sql.ErrNoRows {
		// If no rows are found, return nil
		return 0, nil
	} else if err != nil {
		// Return any other error
		return 0, err
	}

	return userID, nil
}

func formatDatetime(timestamp int64) string {
	t := time.Unix(timestamp, 0).UTC()
	return t.Format("2006-01-02 @ 15:04")
}

func gravatarURL(email string, size int) string {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)
	hash := md5.Sum([]byte(email))
	hashStr := hex.EncodeToString(hash[:])
	return fmt.Sprintf("http://www.gravatar.com/avatar/%s?d=identicon&s=%d", hashStr, size)
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


	query := "SELECT * FROM user WHERE user_id = 1"
	args := []interface{}{"value"}
	result, err := queryDB(db, query, args...)
	if err != nil {
		fmt.Println("Error querying the database:", err)
		return
	}

	fmt.Println(result)


	username := "LIKE %"
	userID, err := getUserID(db, username)
	if err != nil {
		fmt.Println("Error getting user ID:", err)
		return
	}

	if userID == 0 {
		fmt.Println("User not found")
	} else {
		fmt.Println("User ID:", userID)
	}


	timestamp := int64(1644942725) // Replace with your timestamp
	formatted := formatDatetime(timestamp)
	fmt.Println("Formatted datetime:", formatted)


	email := "example@example.com"
	size := 80
	gravatar := gravatarURL(email, size)
	fmt.Println("Gravatar URL:", gravatar)
}


