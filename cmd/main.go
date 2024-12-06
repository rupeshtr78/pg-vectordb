package main

import (
	"context"
	"fmt"
	"log"
	pgembed "pg-vector-db/internal/pg_embed"

	"github.com/pgvector/pgvector-go"
)

func main() {

	conn, err := pgembed.GormCreateConnection(pgembed.PgconnStr)
	if err != nil {
		fmt.Println(err)
	}

	table := "documents_gorm"
	err = pgembed.GormCreateVectorTable(context.Background(), conn, table)
	if err != nil {
		fmt.Println(err)
	}

	input := []string{
		"The dog is barking",
		"The cat is purring",
		"The bear is growling",
		"The lion is roaring",
		"The tiger is snarling",
		"The elephant is trumpeting",
		"Animals are amazing",
		"Sky is blue",
	}
	embedResults, err := pgembed.FetchEmbeddings(input, pgembed.EmbedderUrl, "nomic-embed-text")
	if err != nil {
		fmt.Println(err)
	}

	err = pgembed.GormLoadVectorData(context.Background(), input, embedResults, conn)
	if err != nil {
		fmt.Println(err)
	}

	query := []string{"which animal is roaring"}
	embeddedQuery, err := pgembed.FetchEmbeddings(query, pgembed.EmbedderUrl, pgembed.EmbedModel)
	if err != nil {
		log.Fatalf("Failed to embed Query: %v", err)
	}

	queryVector := pgvector.NewVector(embeddedQuery[0])

	err = pgembed.GormQuerySimilarVectors(context.Background(), conn, queryVector, 1, table)
	if err != nil {
		log.Fatalf("Failed to query similar vectors: %v", err)
	}

}

func RunEmbedding() {
	// pgembed.Connect()
	// pgembed.PgVectorDbEmbed()
	input := []string{
		"The dog is barking",
		"The cat is purring",
		"The bear is growling",
		"The lion is roaring",
		"The tiger is snarling",
		"The elephant is trumpeting",
		"Animals are amazing",
		"Sky is blue",
	}
	embedResults, err := pgembed.FetchEmbeddings(input, pgembed.EmbedderUrl, "nomic-embed-text")
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println(embedResults)
	// find dimension of embeddings
	fmt.Printf("Embedding dimension: %d\n", len(embedResults[0]))

	ctx := context.Background()

	conn, err := pgembed.CreateConnection(ctx)
	if err != nil {
		fmt.Println(err)
	}

	table := "documents"
	err = pgembed.CreateVectorTable(ctx, conn, table)
	if err != nil {
		fmt.Println(err)
	}

	err = pgembed.LoadVectorData(ctx, input, embedResults, conn)
	if err != nil {
		fmt.Println(err)
	}

	pgembed.RunQuery(ctx, conn, []string{"which animal is roaring"}, table)
}
