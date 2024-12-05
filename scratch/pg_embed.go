package scratch

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	pgembed "pg-vector-db/internal/pg_embed"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/pgvector/pgvector-go"
)

func PgVectorDbEmbed() {
	// Connect to the database
	db, err := sql.Open("postgres", pgembed.EmbedderUrl)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	// Sample text to embed
	text := "The quick brown fox jumps over the lazy dog. Paris is capital of France."

	// Get the embedded vector
	embeddedVector, err := embedText(text)
	if err != nil {
		log.Fatalf("Failed to embed text: %v", err)
	}

	// Insert vector into database
	if err := insertVector(db, text, embeddedVector); err != nil {
		log.Fatalf("Failed to insert vector: %v", err)
	}

	// query := "What is the capital of France?"

	// embeddedQuery, err := embedText(query)
	// if err != nil {
	// 	log.Fatalf("Failed to embed Query: %v", err)
	// }

	// // Query similar vectors
	// if err := QuerySimilarVectors(db, embeddedQuery, 1); err != nil {
	// 	log.Fatalf("Failed to query similar vectors: %v", err)
	// }
}

func embedText(text string) (pgvector.Vector, error) {
	// Call the embedding model API
	model := "nomic-embed-text"
	jsonData := fmt.Sprintf(`{"model": "%s", "prompt": "%s"}`, model, text)
	response, err := http.Post(pgembed.EmbedderUrl, "application/json", bytes.NewBuffer([]byte(jsonData)))
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
