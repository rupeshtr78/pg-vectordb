package pgembed

import (
	"context"
	"fmt"

	"github.com/pgvector/pgvector-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GormCreateConnection(connStr string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Check the connection
	if db.Error != nil {
		return nil, db.Error
	}

	// If you want to use context with GORM, you can use it in your queries.
	// GORM does not directly handle context in the Open() method, but you
	// can use it in transaction, find, etc.

	return db, nil
}

func GormCreateVectorTable(ctx context.Context, conn *gorm.DB, tableName string) error {

	fmt.Printf("Creating Vector Table %s\n", tableName)

	createExt := conn.Exec("CREATE EXTENSION IF NOT EXISTS vector")
	if createExt.Error != nil {
		return createExt.Error
	}

	dropTable := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
	drop := conn.Exec(dropTable)
	if drop.Error != nil {
		return drop.Error
	}

	createTable := fmt.Sprintf("CREATE TABLE %s (id bigserial PRIMARY KEY, content text, embedding vector(%d))", tableName, 768)
	createVecTable := conn.Exec(createTable)
	if createVecTable.Error != nil {
		return createVecTable.Error
	}

	createIndex := fmt.Sprintf("CREATE INDEX ON %s USING hnsw(embedding vector_cosine_ops)", tableName)
	index := conn.Exec(createIndex)
	if index.Error != nil {
		return index.Error
	}

	fmt.Printf("Vector Table %s created\n", tableName)
	return nil
}

func GormLoadVectorData(ctx context.Context, input []string, embeddings [][]float32, conn *gorm.DB) error {

	if len(input) != len(embeddings) {
		return fmt.Errorf("input and embeddings size mismatch: %d != %d", len(input), len(embeddings))
	}

	for i, content := range input {
		insert := conn.Exec("INSERT INTO documents (content, embedding) VALUES ($1, $2)", content, pgvector.NewVector(embeddings[i]))
		if insert.Error != nil {
			return fmt.Errorf("failed to insert vector: %v", insert.Error)
		}
	}

	fmt.Println("Data loaded successfully")

	return nil

}
