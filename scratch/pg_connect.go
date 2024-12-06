package scratch

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Connect connects to the Vector database
func Connect() {
	// Define connection parameters
	connStr := "host=10.0.0.213 port=5555 user=rupesh dbname=vectordb sslmode=disable"

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Check the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	fmt.Println("Successfully connected to the Vector database!")

	version, err := db.Query("SHOW server_version")
	if err != nil {
		log.Fatalf("Failed to query the database: %v", err)
	}

	for version.Next() {
		var serverVersion string
		if err := version.Scan(&serverVersion); err != nil {
			log.Fatalf("Failed to scan the database: %v", err)
		}
		fmt.Printf("Server version: %s\n", serverVersion)
	}

	// show all tables
	tables, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
	if err != nil {
		log.Fatalf("Failed to query the database: %v", err)
	}

	for tables.Next() {
		var tableName string
		if err := tables.Scan(&tableName); err != nil {
			log.Fatalf("Failed to scan the database: %v", err)
		}
		fmt.Printf("Table name: %s\n", tableName)
	}

	// describe a table
	tableName := "embeddings"
	columns, err := db.Query(fmt.Sprintf("SELECT column_name FROM information_schema.columns WHERE table_name = '%s'", tableName))
	if err != nil {
		log.Fatalf("Failed to query the database: %v", err)
	}

	for columns.Next() {
		var columnName string
		if err := columns.Scan(&columnName); err != nil {
			log.Fatalf("Failed to scan the database: %v", err)
		}
		fmt.Printf("Column name: %s\n", columnName)
	}

}
