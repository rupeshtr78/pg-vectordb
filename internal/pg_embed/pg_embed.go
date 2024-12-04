package pgembed

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/pgvector/pgvector-go"
)

const (
	connStr  = "host=10.0.0.213 port=5555 user=rupesh dbname=vectordb sslmode=disable"
	modelURL = "http://localhost:11434/api/embeddings" // replace with your url
)

func PgVectorDbEmbed() {
	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	// Sample text to embed
	text := "This is a sample text to embed."

	// Get the embedded vector
	embeddedVector, err := embedText(text)
	if err != nil {
		log.Fatalf("Failed to embed text: %v", err)
	}

	// Insert vector into database
	if err := insertVector(db, text, embeddedVector); err != nil {
		log.Fatalf("Failed to insert vector: %v", err)
	}

	// Query similar vectors
	if err := querySimilarVectors(db, embeddedVector, 5); err != nil {
		log.Fatalf("Failed to query similar vectors: %v", err)
	}
}

func embedText(text string) (pgvector.Vector, error) {
	// Call the embedding model API
	model := "nomic-embed-text"
	jsonData := fmt.Sprintf(`{"model": "%s", "prompt": "%s"}`, model, text)
	response, err := http.Post(modelURL, "application/json", bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		return pgvector.Vector{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return pgvector.Vector{}, fmt.Errorf("failed to get embedding: %s", response.Status)
	}

	var result struct {
		Vector []float32 `json:"embedding"`
	}

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return pgvector.Vector{}, err
	}

	return pgvector.NewVector(result.Vector), nil
}

func insertVector(db *sql.DB, text string, vector pgvector.Vector) error {
	_, err := db.Exec("INSERT INTO embeddings (text, embedding) VALUES ($1, $2)", text, vector)
	return err
}

func querySimilarVectors(db *sql.DB, vector pgvector.Vector, limit int) error {
	rows, err := db.Query("SELECT text FROM embeddings ORDER BY embedding <=> $1 LIMIT $2", vector, limit)
	if err != nil {
		return err
	}
	defer rows.Close()

	fmt.Println("Similar texts:")
	for rows.Next() {
		var text string
		if err := rows.Scan(&text); err != nil {
			return err
		}
		fmt.Println(text)
	}

	return rows.Err()
}
