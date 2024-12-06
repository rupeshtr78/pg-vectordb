package pgembed

import (
	"context"
	"fmt"

	"github.com/pgvector/pgvector-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GormCreateConnection creates a connection to the database using GORM
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

type Documents struct {
	ID        uint            `gorm:"primaryKey"`
	Content   string          `gorm:"type:text"`
	Embedding pgvector.Vector `gorm:"type:vector(768)"`
}

func (Documents) TableName() string {
	return "documents_gorm"
}

// GormCreateVectorTable creates a table with vector extension and hnsw index using GORM
func GormCreateVectorTable(ctx context.Context, conn *gorm.DB, tableName string) error {

	fmt.Printf("Creating Vector Table %s\n", tableName)

	createExt := conn.Exec("CREATE EXTENSION IF NOT EXISTS vector")
	if createExt.Error != nil {
		return createExt.Error
	}

	err := conn.Migrator().DropTable(&Documents{})
	if err != nil {
		return err
	}
	fmt.Printf("Table %s dropped\n", tableName)

	err = conn.Migrator().CreateTable(&Documents{})
	if err != nil {
		return err
	}
	fmt.Printf("Vector Table %s created\n", tableName)

	createIndex := fmt.Sprintf("CREATE INDEX ON %s USING hnsw(embedding vector_cosine_ops)", tableName)
	index := conn.Exec(createIndex)
	if index.Error != nil {
		return index.Error
	}

	fmt.Printf("Index created on %s\n", tableName)

	return nil
}

// GormLoadVectorData loads the vector data into the table using GORM
func GormLoadVectorData(ctx context.Context, input []string, embeddings [][]float32, conn *gorm.DB) error {

	if len(input) != len(embeddings) {
		return fmt.Errorf("input and embeddings size mismatch: %d != %d", len(input), len(embeddings))
	}

	for i, content := range input {
		doc := &Documents{
			Content:   content,
			Embedding: pgvector.NewVector(embeddings[i]),
		}
		// insert := conn.Exec("INSERT INTO documents (content, embedding) VALUES ($1, $2)", content, pgvector.NewVector(embeddings[i]))
		insert := conn.Create(doc)
		if insert.Error != nil {
			return fmt.Errorf("failed to insert vector: %v", insert.Error)
		}
	}

	fmt.Println("Data loaded successfully")

	return nil

}

type VectorDBQuery struct {
	Content   string          `gorm:"column:content"`
	Embedding pgvector.Vector `gorm:"column:embedding"`
}

// GormQuerySimilarVectors queries similar vectors using GORM
func GormQuerySimilarVectors(ctx context.Context, db *gorm.DB, vector pgvector.Vector, limit int, table string) error {
	// Query similar vectors using GORM
	results := []VectorDBQuery{}

	// Build the query
	vectorStr := "('" + vector.String() + "')"
	order := fmt.Sprintf("embedding <=> %s", vectorStr)
	query := db.Table(table).Select("content").Order(order).Limit(limit)
	if query.Error != nil {
		return query.Error
	}

	// Execute the query
	execQuery := query.Find(&results)
	if execQuery.Error != nil {
		return execQuery.Error
	}

	fmt.Println("Getting Results from VectorDb:")
	for _, result := range results {
		fmt.Println(result.Content)
	}

	return nil
}
