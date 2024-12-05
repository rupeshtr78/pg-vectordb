package pgembed

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
	pgxvector "github.com/pgvector/pgvector-go/pgx"
)

func CreateConnection(ctx context.Context) (*pgx.Conn, error) {

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func CreateVectorTable(ctx context.Context, conn *pgx.Conn, tableName string) error {

	_, err := conn.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS vector")
	if err != nil {
		return err
	}

	err = pgxvector.RegisterTypes(ctx, conn)
	if err != nil {
		return err
	}

	dropTable := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
	_, err = conn.Exec(ctx, dropTable)
	if err != nil {
		return err
	}

	createTable := fmt.Sprintf("CREATE TABLE %s (id bigserial PRIMARY KEY, content text, embedding vector(%d))", tableName, modelDimension)
	_, err = conn.Exec(ctx, createTable)
	if err != nil {
		return err
	}

	createIndex := fmt.Sprintf("CREATE INDEX ON %s USING hnsw(embedding vector_cosine_ops)", tableName)
	_, err = conn.Exec(ctx, createIndex)
	if err != nil {
		return err
	}

	fmt.Printf("Vector Table %s created\n", tableName)
	return nil
}

func LoadVectorData(ctx context.Context, input []string, embeddings [][]float32, conn *pgx.Conn) error {

	if len(input) != len(embeddings) {
		return fmt.Errorf("input and embeddings size mismatch: %d != %d", len(input), len(embeddings))
	}

	for i, content := range input {
		_, err := conn.Exec(ctx, "INSERT INTO documents (content, embedding) VALUES ($1, $2)", content, pgvector.NewVector(embeddings[i]))
		if err != nil {
			return fmt.Errorf("failed to insert vector: %v", err)
		}
	}

	fmt.Println("Data loaded successfully")

	return nil

}
