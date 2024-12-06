package pgembed

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
)

// RunQuery runs a query against the vector database
func RunQuery(ctx context.Context, conn *pgx.Conn, query []string, table string) {

	embeddedQuery, err := FetchEmbeddings(query, EmbedderUrl, EmbedModel)
	if err != nil {
		log.Fatalf("Failed to embed Query: %v", err)
	}

	// Query similar vectors
	queryVector := pgvector.NewVector(embeddedQuery[0])
	if err := QuerySimilarVectors(ctx, conn, queryVector, QueryLimit, table); err != nil {
		log.Fatalf("Failed to query similar vectors: %v", err)
	}
}

// Query similar vectors to a given vector
// <-> - L2 distance
// <#> - (negative) inner product
// <=> - cosine distance
// <+> - L1 distance (added in 0.7.0)
// <~> - Hamming distance (binary vectors, added in 0.7.0)
// <%> - Jaccard distance (binary vectors, added in 0.7.0)
func QuerySimilarVectors(ctx context.Context, conn *pgx.Conn, vector pgvector.Vector, limit int, table string) error {

	// Get the nearest neighbors to a vector
	query := fmt.Sprintf("SELECT content FROM %s ORDER BY embedding <=> $1 LIMIT $2", table)

	// Get the nearest neighbors to a row
	// nearestneighbors := fmt.Sprintf("SELECT content FROM %s WHERE id != 1 ORDER BY embedding <-> (SELECT embedding FROM %s WHERE id = 1) LIMIT $1;", table, table)

	rows, err := conn.Query(ctx, query, vector, limit)
	// rows, err := conn.Query(ctx, nearestneighbors, limit)
	if err != nil {
		return err
	}
	defer rows.Close()

	fmt.Println("Querying VectorDb:")
	for rows.Next() {
		var text string
		if err := rows.Scan(&text); err != nil {
			return err
		}
		fmt.Println(text)
	}

	return rows.Err()
}
